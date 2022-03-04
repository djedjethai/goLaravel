package render

// not that we set up our views test into Goframwork
// because we do not want myapp to have to stick with Goframwork to run test

// this file will allow us to run some functions before test are executed

import (
	"github.com/CloudyKit/jet/v6"
	"os"
	"testing"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./testdata/views"),
	jet.InDevelopmentMode(),
)

var testRenderer = Render{
	Renderer: "",
	RootPath: "",
	JetViews: views,
}

// any time i run tests in this dir, if go sees the file setup_test.go
// with a func named TestMain(*testing.M), it will run that func
// and then that func will run any test found in that directory
func TestMain(m *testing.M) {

	// this m.Run() is the one runnning the tests
	os.Exit(m.Run())
}
