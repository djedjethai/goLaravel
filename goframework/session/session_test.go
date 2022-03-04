package session

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"reflect"
	"testing"
)

func TestSession_InitSession(t *testing.T) {
	// Arrange
	session := &Session{
		CookieLifetime: "60",
		CookiePersist:  "true",
		CookieName:     "mycookie",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}

	// Act
	var sm *scs.SessionManager
	sess := session.InitSession()

	var sessKind reflect.Kind
	var sessType reflect.Type

	rv := reflect.ValueOf(sess)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fmt.Println("for loop: ", rv.Kind(), rv.Type(), rv)
		sessKind = rv.Kind()
		sessType = rv.Type()

		rv = rv.Elem()
	}

	// Assert, assert the type of sess which should be os sm's type
	// Method .IsValid() built into rv
	if !rv.IsValid() {
		t.Error("Invalid type or kind, kind: ", rv.Kind(), " type: ", rv.Type())
	}

	// get the kind of the var we declared before
	if sessKind != reflect.ValueOf(sm).Kind() {
		t.Error("Wrong kind returned testing cookie session, expected: ", reflect.ValueOf(sm).Kind(), " and got: ", sessKind)
	}

	// get the type of the var we declared before
	if sessType != reflect.ValueOf(sm).Type() {
		t.Error("Wrong Type returned testing cookie session, expected: ", reflect.ValueOf(sm).Type(), " and got: ", sessType)
	}
}
