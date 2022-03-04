//go:build unit

// go tag to indicate this test will only run for unit test
package data

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
