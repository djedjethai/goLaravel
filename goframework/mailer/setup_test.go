package mailer

// !!!! to run the tests
// --count=1 to avoid the cache to run
// sudo as docker need to be run
// sudo go test . --count=1

// get the coverage in browser
// sudo go test --coverprofile=coverage.out && go tool cover -html=coverage.out

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var pool *dockertest.Pool
var resource *dockertest.Resource

// mailhog port is 1025,
// but we may already use this one
// so lets use 1026
var mailer = Mail{
	Domain:      "localhost",
	Templates:   "./testdata/mail",
	Host:        "localhost",
	Port:        1026,
	Encryption:  "none",
	FromAddress: "me@whatever.com",
	FromName:    "John",
	Jobs:        make(chan Message, 1),
	Results:     make(chan Result, 1),
}

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("could not connect to docker", err)
	}

	pool = p

	// 1025 on the docker img, 8025 webinterface
	opts := dockertest.RunOptions{
		Repository:   "mailhog/mailhog",
		Tag:          "latest",
		Env:          []string{},
		ExposedPorts: []string{"1025", "8025"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"1025": {
				{HostIP: "0.0.0.0", HostPort: "1026"},
			},
			"8025": {
				{HostIP: "0.0.0.0", HostPort: "8026"},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Println(err)
		_ = pool.Purge(resource)
		log.Fatal("Could not start resource")
	}

	time.Sleep(2 * time.Second)

	// that will run on the background for the  duration of our tests
	go mailer.ListenForMail()

	code := m.Run()

	// if i comment these lines the docker containers will not stop
	// a good way to check into mailhog
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	// after we ran that check, we exit
	// at exit time all process will be kill, even the channel
	os.Exit(code)
}
