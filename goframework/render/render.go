package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
	Session    *scs.SessionManager
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
	Error           string
	Flash           string
}

func (c *Render) defaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Secure = c.Secure
	td.ServerName = c.ServerName
	td.CSRFToken = nosurf.Token(r)
	td.Port = c.Port

	// determine if there is a value call userID in the session
	if c.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = true
	}
	// PopString() return the session value from a given key and the delete it
	td.Error = c.Session.PopString(r.Context(), "error")
	td.Flash = c.Session.PopString(r.Context(), "flash")
	return td
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
		// the word "variables" is require by jet
		return c.JetPage(w, r, view, variables, data)
	default:

	}

	return errors.New("no rendering engine specified")
}

// goPage renders a standard go template
func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		// cast the data interface{} to be *TemplateData
		td = data.(*TemplateData)
	}

	// at this point we have data(if any)
	err = tmpl.Execute(w, &td)
	if err != nil {
		return err
	}

	return nil
}

// jetpage render a template using the het templating engine
func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	// jet.varMap is a data structure that jet use to pass datas to the template
	var vars jet.VarMap

	if variables == nil {
		// initialize  vars
		vars = make(jet.VarMap)
	} else {
		// in case variables have datas (we cast it as well)
		vars = variables.(jet.VarMap)
	}

	// deal with the case we already have datas in our templateData
	td := &TemplateData{}
	if data != nil {
		// will add the passed data with the already existing one(from TemplateData{})
		td = data.(*TemplateData)
	}

	// add the some property(IsAuthenticated, Secure, etc) into any rendered templateData
	td = c.defaultData(td, r)

	// load the template to render it
	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
