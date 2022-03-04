package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// to avoid long Testing function, we can use a
// common convention in the Go world known as Table Test
var pageData = []struct {
	name          string
	renderer      string
	template      string
	errorExpected bool
	errorMessage  string
}{
	{"go_page", "go", "home", false, "error rendering go template"},
	{"go_page_no_template", "go", "no-file", true, "no error rendering no existent go template when one is expecting"},
	{"jet_page", "jet", "home", false, "error rendering jet template"},
	{"jet_page_no_template", "jet", "no-file", true, "no error rendering no existent jet template when one is expecting"},
	{"invalid renderer engine", "foo", "home", true, "no error rendering with no existent template engine"},
}

// test the Page()
func TestRender_Page(t *testing.T) {

	w := httptest.NewRecorder()
	testRenderer.RootPath = "./testdata"

	for _, dt := range pageData {
		r, err := http.NewRequest(http.MethodGet, "home", nil)
		if err != nil {
			t.Error(err)
		}

		testRenderer.Renderer = dt.renderer
		err = testRenderer.Page(w, r, dt.template, nil, nil)
		if dt.errorExpected {
			if err == nil {
				t.Errorf("%s: %s", dt.name, dt.errorMessage)
			}
		} else {
			if err != nil {
				t.Errorf("%s: %s: %s", dt.name, dt.errorMessage, err.Error())
			}
		}
	}
}

// useless test but the teacher made it
// func Test_GoPage(t *testing.T) {
// 	// Arrange
// 	r, err := http.NewRequest(http.MethodGet, "some-url", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	w := httptest.NewRecorder()
//
// 	testRenderer.Renderer = "go"
// 	testRenderer.RootPath = "./testdata"
//
// 	// Act and Assert
// 	err = testRenderer.Page(w, r, "home", nil, nil)
// 	if err != nil {
// 		t.Error("Error rendering go page using GoPage() should not return an err", err)
// 	}
// }
//
// // useless test but the teacher made it
// func Test_JetPage(t *testing.T) {
// 	// Arrange
// 	r, err := http.NewRequest(http.MethodGet, "some-url", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	w := httptest.NewRecorder()
//
// 	testRenderer.Renderer = "jet"
// 	testRenderer.RootPath = "./testdata"
//
// 	// Act and Assert
// 	err = testRenderer.Page(w, r, "home", nil, nil)
// 	if err != nil {
// 		t.Error("Error rendering go page using JetPage() should not return an err", err)
// 	}
// }
