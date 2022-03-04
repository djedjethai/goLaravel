//go:build integration

// run  test with this command(--count=1, to do not cache any test)
// sudo go test . --tags integration --count=1
// sudo go test -v . --tags integration --count=1 // to get some comments

// see the integration coverage
// sudo go test -cover . --tags integration

// see our coverage in the browser, run that first
// sudo go test -coverprofile=coverage.out . --tags integration
// then
// go tool cover -html=coverage.out

package data

// package to run docker from here :)
// go get github.com/ory/dockertest/v3
// go get github.com/ory/dockertest/v3/docker

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "secret"
	dbName   = "goframework_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var dummyUser = User{
	FirstName: "Some",
	LastName:  "Guy",
	Email:     "me@here.com",
	Active:    "1",
	Password:  "password",
}

var models Models
var testDB *sql.DB
var resource *dockertest.Resource
var pool *dockertest.Pool

// can not have 2 TestMain Func in a single package
// but we already have one in setup_test, that s we are using build tags
func TestMain(m *testing.M) {
	os.Setenv("DATABASE_TYPE", "postgres")
	// set var env for upperDB to set its log to ERROR (no WARNING which is default)
	os.Setenv("UPPER_DB_LOG", "ERROR")

	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.4",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// to spin up the container and the db will take time
	// so we have to wait before running our tests, lets use a docker feature
	if err := pool.Retry(func() error {
		var err error

		// open a connection with the db
		// "pgx" bc we use jackc/pgx package
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		// cocker will still work until the ping still ping
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to docker: %s", err)
	}

	// when we passed all of that we got docker and postgres running
	// and a database named goframework_test, but still no tables
	err = createTables(testDB)
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	models = New(testDB)

	code := m.Run()

	// get rid out of the built image
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables(db *sql.DB) error {
	stmt := `
	CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

drop table if exists users cascade;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    user_active integer NOT NULL DEFAULT 0,
    email character varying(255) NOT NULL UNIQUE,
    password character varying(60) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists remember_tokens;

CREATE TABLE remember_tokens (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    remember_token character varying(100) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON remember_tokens
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists tokens;

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    first_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    token character varying(255) NOT NULL,
    token_hash bytea NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    expiry timestamp without time zone NOT NULL
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON tokens
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
	`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

// test User
func TestUser_Table(t *testing.T) {
	s := models.Users.Table()
	if s != "users" {
		t.Error("wrong table name returned: ", s)
	}
}

func TestUser_Insert(t *testing.T) {
	// Act
	id, err := models.Users.Insert(dummyUser)

	if err != nil {
		t.Error("While testing Insert a User, the err should be nil", err)
	}
	if id == 0 {
		t.Error("While testing Insert a User, the userID should be returned")
	}
}

func TestUser_Get(t *testing.T) {
	u, err := models.Users.Get(1)

	if err != nil {
		t.Error("While testing Get() a User, the err should be nil", err)
	}
	if u.FirstName != "Some" {
		t.Error("While testing Get() a User, the user should be returned")
	}
}

func TestUser_GetByEmail(t *testing.T) {
	u, err := models.Users.GetByEmail("me@here.com")
	if err != nil {
		t.Error("While testing GetByEmail() a User, the err should be nil", err)
	}

	if u.FirstName != "Some" {
		t.Error("While testing GetByEmail a User, the user should be returned")
	}
}

func TestUser_GetAll(t *testing.T) {
	var dummyUser1 = User{
		FirstName: "Some1",
		LastName:  "Guy1",
		Email:     "me@here.com1",
		Active:    "1",
		Password:  "password1",
	}

	_, _ = models.Users.Insert(dummyUser1)
	allUser, err := models.Users.GetAll()
	if err != nil {
		t.Error("While testing GetAll() User, the err should be nil", err)
	}

	if len(allUser) != 2 {
		t.Error("While testing Insert a User, the userID should be returned")
	}
}

func TestUser_Update(t *testing.T) {
	user, _ := models.Users.Get(1)
	user.FirstName = "updated"

	err := models.Users.Update(*user)
	updatedUser, _ := models.Users.Get(1)

	if err != nil {
		t.Error("While testing Update a User, the err should be nil", err)
	}
	if updatedUser.FirstName != "updated" {
		t.Error("while testing Update a User, the user FirstName should be updated")
	}
}

func TestUser_PasswordMatches(t *testing.T) {
	user, _ := models.Users.Get(1)

	isMatch, err := user.PasswordMatches("password")
	if err != nil {
		t.Error("While testing a User PasswordMatches, the err should be nil", err)
	}
	if !isMatch {
		t.Error("While testing a User correct password, should return true")
	}

	isMatch, err = user.PasswordMatches("wrongPassword")
	if err != nil {
		t.Error("While testing a User PasswordMatches, the err should be nil", err)
	}
	if isMatch {
		t.Error("While testing a User incorrect password, should return false")
	}
}

func TestUser_ResetPassword(t *testing.T) {
	err := models.Users.ResetPassword(1, "newpassword")
	if err != nil {
		t.Error("while tetsting a User ResetPassword, the err should be nil", err)
	}

	err = models.Users.ResetPassword(3, "newpassword")
	if err == nil {
		t.Error("while testing a User ResetPassword, the password the err should not be nil")
	}
}

func TestUser_Delete(t *testing.T) {
	err := models.Users.Delete(1)
	if err != nil {
		t.Error("While testing a User Delete should not return an err", err)
	}
	err = models.Users.Delete(1)
	if err != nil {
		t.Error("While testing a User Delete should return an err")
	}
	_, err = models.Users.Get(1)
	if err == nil {
		t.Error("While testing a User Delete, User deleted should be gone")
	}
}

// test Token
func TestToken_Table(t *testing.T) {
	tb := models.Tokens.Table()

	if tb != "tokens" {
		t.Error("While testing Token, Table() should return 'tokens'")
	}
}

func TestToken_GenerateToken(t *testing.T) {
	// ttl for 1 year
	ttl := time.Hour * 24 * 365
	_, err := models.Tokens.GenerateToken(2, ttl)

	if err != nil {
		t.Error("While testing Token, GenerateToken() should not return an err", err)
	}
}

func TestToken_Insert(t *testing.T) {
	tkn, _ := models.Tokens.GenerateToken(2, time.Hour*24*365)
	u, _ := models.Users.Get(2)

	err := models.Tokens.Insert(*tkn, *u)
	if err != nil {
		t.Error("While testing Token, Insert() should not return an err", err)
	}

}

func TestToken_GetUserForToken(t *testing.T) {
	tkn, _ := models.Tokens.GenerateToken(2, time.Hour*24*365)
	u, _ := models.Users.Get(2)
	_ = models.Tokens.Insert(*tkn, *u)

	user, err := models.Tokens.GetUserForToken(tkn.PlainText)
	if err != nil {
		t.Error("While testing Token, GetUserForToken() should not return an err", err)
	}
	if user.FirstName != u.FirstName {
		t.Error("While testing Token, GetUserForToken() return wrong user")
	}

	_, err = models.Tokens.GetUserForToken("unknowtoken")
	if err == nil {
		t.Error("While testing Token, GetUserForToken() should return an err when unknowToken", err)
	}
}

func TestToken_GetTokensForUser(t *testing.T) {
	tkns, err := models.Tokens.GetTokensForUser(2)
	if err != nil {
		t.Error("While testing Token, GetTokensForUser() should not return an err", err)
	}

	if len(tkns) < 1 {
		t.Error("While testing Token, GetTokensForUser() tokens should not be nil")
	}

	_, err = models.Tokens.GetTokensForUser(10)
	if err == nil {
		t.Error("While testing Token, GetTokensForUser() should return an err")
	}

}

func TestToken_Get(t *testing.T) {
	u, _ := models.Users.Get(2)
	tkn, err := models.Tokens.Get(u.Token.ID)
	if err != nil {
		t.Error("While testing Token, Get() should not return an err", err)
	}

	if tkn == nil {
		t.Error("While testing Token, Get(), the token should not be nil")
	}
}

func TestToken_GetByToken(t *testing.T) {
	u, _ := models.Users.Get(2)
	tkn, err := models.Tokens.GetByToken(u.Token.PlainText)
	if err != nil {
		t.Error("While testing Token, GetByToken() should not return an err", err)
	}
	if tkn == nil {
		t.Error("While testing Token, GetByToken(), the token should not be nil")
	}

	tkn, err = models.Tokens.GetByToken("invalidToken")
	if err == nil {
		t.Error("While testing Token, GetByToken() should return an err")
	}
	if tkn != nil {
		t.Error("While testing Token, GetByToken(), the token should be nil")
	}
}

// test authentification
var authData = []struct {
	name        string
	token       string
	email       string
	errExpected bool
	message     string
}{
	{"invalid", "abcdefghejklmnopqrstuvwxyz", "a@here.com", true, "invalid token accepted as valid"},
	{"invalid_length", "abcdefghejklmnopqrstuvwxy", "a@here.com", true, "invalid token of wrong length accepted as valid"},
	{"no_user", "abcdefghejklmnopqrstuvwxyz", "a@here.com", true, "no user but token accepted as valid"},
	{"valid", "", "me@here.com1", false, "valid token reported as invalid"},
}

var dummyUser1 = User{
	FirstName: "Some1",
	LastName:  "Guy1",
	Email:     "me@here.com1",
	Active:    "1",
	Password:  "password",
}

func TestToken_AuthenticateToken(t *testing.T) {
	// get user id = 2, as dummyUser(id = 1) has been deleted
	u, _ := models.Users.Get(2)

	for _, tt := range authData {
		token := ""
		if tt.email == u.Email {
			user, err := models.Users.GetByEmail(tt.email)

			if err != nil {
				t.Error("While testing auth, fail to get user", err)
			}
			token = user.Token.PlainText
		} else {
			// set token setted up in test slice
			token = tt.token
		}

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		_, err := models.Tokens.AuthenticateToken(req)
		if tt.errExpected && err == nil {
			t.Errorf("%s: %s", tt.name, tt.message)
		} else if !tt.errExpected && err != nil {
			t.Errorf("%s: %s - %s", tt.name, tt.message, err)
		} else {
			t.Logf("passed %s", tt.name)
		}
	}
}

func TestToken_Delete(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUser1.Email)
	if err != nil {
		t.Error(err)
	}

	err = models.Tokens.DeleteByToken(u.Token.PlainText)
	if err != nil {
		t.Error("While testToken, DeleteByToken should not return an err", err)
	}
}

func TestToken_ExpiredToken(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUser1.Email)
	if err != nil {
		t.Error(err)
	}

	// generate an expired token(can do that!)
	tkn, err := models.Tokens.GenerateToken(u.ID, -time.Hour)
	if err != nil {
		t.Error(err)
	}

	err = models.Tokens.Insert(*tkn, *u)
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer "+tkn.PlainText)

	_, err = models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("While testToken, ExpiredToken should return an err")
	}
}

func TestToken_BadHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	_, err := models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("While testToken missing auth header, BadHeader should return an err")
	}

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "abc")
	_, err = models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("While testToken fail to catch bad auth header, BadHeader should return an err")
	}

	newUser := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "you@there.com",
		Active:    "1",
		Password:  "abc",
	}

	id, err := models.Users.Insert(newUser)
	if err != nil {
		t.Error(err)
	}

	token, err := models.Tokens.GenerateToken(id, 1*time.Hour)
	if err != nil {
		t.Error(err)
	}

	err = models.Tokens.Insert(*token, newUser)
	if err != nil {
		t.Error(err)
	}

	// here we test the deletion of the token
	// as in the db there is a CASCADE DELETION for the token
	// at the time we delet the user
	err = models.Users.Delete(id)
	if err != nil {
		t.Error(err)
	}

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer "+token.PlainText)
	_, err = models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("While testToken when user has been deleted, should return an err as token has been deleted (cascade delete from user)")
	}
}

func TestToken_DeleteNonExistingToken(t *testing.T) {
	err := models.Tokens.DeleteByToken("abc")
	if err != nil {
		t.Error("While testToken delete inexisting token, should return an err")
	}
}

func TestToken_ValidToken(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUser1.Email)
	if err != nil {
		t.Error(t)
	}

	newToken, err := models.Tokens.GenerateToken(u.ID, 24*time.Hour)
	if err != nil {
		t.Error(t)
	}

	err = models.Tokens.Insert(*newToken, *u)
	if err != nil {
		t.Error(t)
	}

	// test valid token
	okay, err := models.Tokens.ValidToken(newToken.PlainText)
	if err != nil {
		t.Error("While testToken valid token, should not return an err", err)
	}

	if !okay {
		t.Error("While testToken valid token should be reported as valid")
	}

	// test invalid token
	okay, _ = models.Tokens.ValidToken("abc")
	if okay {
		t.Error("While testToken valid token should be reported as invalid")
	}

	// lets refresh our user to work with fresh datas from db
	u, err = models.Users.GetByEmail(dummyUser1.Email)
	if err != nil {
		t.Error(t)
	}

	_ = models.Tokens.DeleteById(u.Token.ID)

	// check for a gone token
	okay, err = models.Tokens.ValidToken(u.Token.PlainText)
	if err == nil {
		t.Error(t)
	}
	if okay {
		t.Error("While testToken unexistent token should be reported as invalid")
	}
}
