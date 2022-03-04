package cache

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/dgraph-io/badger/v3"
	"github.com/gomodule/redigo/redis"
)

var testRedisCache RedisCache
var testBadgerCache BadgerCache

func TestMain(m *testing.M) {
	// set redis cache test
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	pool := redis.Pool{
		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", s.Addr()) // s.Addr give back in-memory of redis
		},
	}

	testRedisCache.Conn = &pool
	testRedisCache.Prefix = "test-goframework"

	defer testRedisCache.Conn.Close()

	// set badger cache test
	// delete the test badger db
	_ = os.RemoveAll("./testdata/tmp/badger")

	// create a badger database
	if _, err := os.Stat("./testdata/tmp"); os.IsNotExist(err) {
		// if this folder do not exist, create it
		err := os.Mkdir("./testdata/tmp", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	// now we make the db
	err = os.Mkdir("./testdata/tmp/badger", 0755)
	if err != nil {
		log.Fatal(err)
	}

	db, _ := badger.Open(badger.DefaultOptions("./testdata/tmp/badger"))

	testBadgerCache.Conn = db

	os.Exit(m.Run())

}
