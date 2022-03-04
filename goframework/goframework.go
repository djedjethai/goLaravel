package goframework

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger/v3"
	"github.com/djedjethai/goframework/cache"
	"github.com/djedjethai/goframework/mailer"
	"github.com/djedjethai/goframework/render"
	"github.com/djedjethai/goframework/session"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

var myRedisCache *cache.RedisCache
var myBadgerCache *cache.BadgerCache
var redisPool *redis.Pool
var badgerConn *badger.DB

type Goframework struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	RootPath      string
	Routes        *chi.Mux
	Render        *render.Render
	Session       *scs.SessionManager
	DB            Database
	JetViews      *jet.Set
	config        config
	EncryptionKey string
	Cache         cache.Cache
	Scheduler     *cron.Cron
	Mail          mailer.Mail
	Server        Server
}

type Server struct {
	ServerName string
	Port       string
	Secure     bool
	URL        string
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
	redis       RedisConfig
}

func (g *Goframework) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "mail", "data", "public", "tmp", "logs", "middleware"},
	}

	err := g.Init(pathConfig)
	if err != nil {
		return err
	}

	err = g.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create loggers && set various properties in Goframework
	infoLog, errorLog := g.startLoggers()

	// connect to database
	if os.Getenv("DATABASE_TYPE") != "" {
		// db is the connection pool
		db, err := g.OpenDB(os.Getenv("DATABASE_TYPE"), g.BuildDSN())
		if err != nil {
			errorLog.Println(err)
			// stop as we can not connect to the db
			os.Exit(1)
		}

		g.DB = Database{
			DataType: os.Getenv("DATABASE_TYPE"),
			Pool:     db,
		}
	}

	if os.Getenv("CACHE") == "redis" || os.Getenv("SESSION_TYPE") == "redis" {
		myRedisCache = g.CreateClientRedisCache()
		g.Cache = myRedisCache
		redisPool = myRedisCache.Conn
	}

	// set badger Scheduler before to use it
	scheduler := cron.New()
	g.Scheduler = scheduler

	if os.Getenv("CACHE") == "badger" {
		myBadgerCache = g.CreateClientBadgerCache()
		g.Cache = myBadgerCache
		badgerConn = myBadgerCache.Conn

		// with Badger it need periodic house keeping
		// need garbage collection on the Badger db
		// need to run that once a day
		_, err = g.Scheduler.AddFunc("@daily", func() {
			// that the garbage collection
			_ = myBadgerCache.Conn.RunValueLogGC(0.7)
		})
		if err != nil {
			return err
		}
	}

	g.InfoLog = infoLog
	g.ErrorLog = errorLog
	g.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	g.Version = version
	g.RootPath = rootPath
	// cannot use g.routes() (value of type http.Handler), so we cast the type
	g.Routes = g.routes().(*chi.Mux)
	g.Mail = g.createMailer()

	// set goframework configuration
	g.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSISTS"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			// call BuildDSN() twice,
			// but that will be trigger just at start, so ok.
			dsn:      g.BuildDSN(),
			database: os.Getenv("DATABASE_TYPE"),
		},
		redis: RedisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
	}

	secure := true
	if strings.ToLower(os.Getenv("SECURE")) == "false" {
		secure = false
	}

	g.Server = Server{
		ServerName: os.Getenv("SERVER_NAME"),
		Port:       os.Getenv("PORT"),
		Secure:     secure,
		URL:        os.Getenv("APP_URL"),
	}

	// create session
	sess := session.Session{
		CookieLifetime: g.config.cookie.lifetime,
		CookiePersist:  g.config.cookie.persist,
		CookieName:     g.config.cookie.name,
		SessionType:    g.config.sessionType,
		CookieDomain:   g.config.cookie.domain,
	}

	// set the  session storage to redis(if used) or to the db
	switch g.config.sessionType {
	case "redis":
		sess.RedisPool = myRedisCache.Conn
	case "mysql", "postgres", "mariadb", "postgresql":
		sess.DBPool = g.DB.Pool
	}

	g.Session = sess.InitSession()

	g.EncryptionKey = os.Getenv("KEY")

	// jet.InDevelopmentMode() is good for dev
	// as it allows to modify templates without re-compile
	// but for prod, it's not as it use some memory
	// and go Gargabe collector have clean it up after
	if g.Debug {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
			jet.InDevelopmentMode(),
		)
		g.JetViews = views
	} else {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		)
		g.JetViews = views
	}

	g.createRenderer()
	// create the mail channel to listen on the mailer
	// so this process will run on the background
	go g.Mail.ListenForMail()

	return nil
}

func (g *Goframework) CreateClientRedisCache() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Conn:   g.createRedisPool(),
		Prefix: g.config.redis.prefix,
	}

	return &cacheClient
}

func (g *Goframework) CreateClientBadgerCache() *cache.BadgerCache {
	cacheClient := cache.BadgerCache{
		Conn: g.createBadgerConn(),
	}
	return &cacheClient
}

func (g *Goframework) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				g.config.redis.host,
				redis.DialPassword(g.config.redis.password))
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}

func (g *Goframework) createBadgerConn() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(g.RootPath + "/tmp/badger"))
	if err != nil {
		return nil
	}
	return db
}

func (g *Goframework) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it does not exist
		err := g.createDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// listenAndServe starts the web server
func (g *Goframework) ListenAndServe() {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     g.ErrorLog,
		Handler:      g.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	// as soon as we stop our app this func will stop
	// then it close the database pool
	if g.DB.Pool != nil {
		defer g.DB.Pool.Close()
	}

	if redisPool != nil {
		defer redisPool.Close()
	}

	if badgerConn != nil {
		defer badgerConn.Close()
	}

	g.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := server.ListenAndServe()
	g.ErrorLog.Fatal(err)
}

func (g *Goframework) checkDotEnv(rootPath string) error {
	err := g.CreateFileIfNotExists(fmt.Sprintf("%s/.env", rootPath))
	if err != nil {
		return err
	}
	return nil
}

func (g *Goframework) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

// we don t use g to make sure we have the necessary values
func (g *Goframework) createRenderer() {
	g.Render = &render.Render{
		Renderer: g.config.renderer,
		RootPath: g.RootPath,
		Port:     g.config.port,
		JetViews: g.JetViews,
		Session:  g.Session,
	}
}

func (g *Goframework) createMailer() mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	m := mailer.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Templates:   g.RootPath + "/mail",
		Host:        os.Getenv("SMTP_HOST"),
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Encryption:  os.Getenv("SMTP_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		Jobs:        make(chan mailer.Message, 20), // max size, 20 mails
		Results:     make(chan mailer.Result, 20),
		API:         os.Getenv("MAILER_API"),
		APIKeys:     os.Getenv("MAILER_KEY"),
		APIUrl:      os.Getenv("MAILER_URL"),
	}

	return m
}

func (g *Goframework) BuildDSN() string {
	// that where we store our final ready to go connection string
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		// build a connection without password
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))

		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
	default:

	}

	return dsn
}
