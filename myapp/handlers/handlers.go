package handlers

import (
	// "fmt"
	"fmt"
	"myapp/data"
	"net/http"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/djedjethai/goframework"
)

type Handlers struct {
	App    *goframework.Goframework
	Models data.Models
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	defer h.App.LoadTime(time.Now())
	err := h.render(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}

func (h *Handlers) GoPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.GoPage(w, r, "home", nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}

func (h *Handlers) JetPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.JetPage(w, r, "jet-template", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}

func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	data := "bar"

	// add some data to the session
	h.App.Session.Put(r.Context(), "foo", data)

	// get the data out of the session(for the exercice)
	value := h.App.Session.GetString(r.Context(), "foo")

	// define a structure of data to pass to jet
	vars := make(jet.VarMap)
	vars.Set("foo", value)

	err := h.App.Render.JetPage(w, r, "sessions", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}

// handler to try various response format, here JSON
func (h *Handlers) JSON(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID      int64    `json:"id"`
		Name    string   `json:"name"`
		Hobbies []string `json:"hobbies"`
	}

	payload.ID = 10
	payload.Name = "jack json"
	payload.Hobbies = []string{"tenis", "program", "swimming"}

	err := h.App.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}

// handler to try various response format, here XML
func (h *Handlers) XML(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		ID      int64    `xml:"id"`
		Name    string   `xml:"name"`
		Hobbies []string `xml:"hobbies>hobby"`
	}

	var payload Payload

	payload.ID = 10
	payload.Name = "jack xml"
	payload.Hobbies = []string{"tenis", "program", "swimming"}

	err := h.App.WriteXML(w, http.StatusOK, payload)
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}

// handler to try various response format, here DownloadFile
func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	err := h.App.DownloadFile(w, r, "./public/images", "celeritas.jpg")
	if err != nil {
		h.App.ErrorLog.Println("error rendering", err)
	}
}

func (h *Handlers) TestCrypto(w http.ResponseWriter, r *http.Request) {
	plainText := "Hello world"
	fmt.Fprint(w, "Unencrypted: "+plainText+"\n")

	encrypted, err := h.encrypt(plainText)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.Error500(w)
		return
	}
	fmt.Fprint(w, "Encrypted: "+encrypted+"\n")

	decrypted, err := h.decrypt(encrypted)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.Error500(w)
		return
	}
	fmt.Fprint(w, "Decrypted: "+decrypted+"\n")

}
