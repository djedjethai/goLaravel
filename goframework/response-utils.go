package goframework

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
)

// make sure the incoming json payload is usable
func (g *Goframework) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{}) // to make sure there is only one value in it(that why{}{})
	if err != io.EOF {
		return errors.New("body must only have a single json value")
	}

	return nil
}

func (g *Goframework) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	// json.MarshalIndent is only for dev
	// second arg is a prefix, i don't know for what, but it's optional
	// we indent things using \t
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// can one or more headers
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (g *Goframework) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (g *Goframework) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile string, fileName string) error {
	fp := path.Join(pathToFile, fileName)
	// clean the path for client who could trick the path to get another file
	fileToServe := filepath.Clean(fp)

	// set the header to tell him to download the file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; file=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)

	return nil
}

func (g *Goframework) Error404(w http.ResponseWriter) {
	g.ErrorStatus(w, http.StatusNotFound)
}

func (g *Goframework) Error500(w http.ResponseWriter) {
	g.ErrorStatus(w, http.StatusInternalServerError)
}

func (g *Goframework) ErrorUnauthorized(w http.ResponseWriter) {
	g.ErrorStatus(w, http.StatusUnauthorized)
}

func (g *Goframework) ErrorForbidden(w http.ResponseWriter) {
	g.ErrorStatus(w, http.StatusForbidden)
}

func (g *Goframework) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
