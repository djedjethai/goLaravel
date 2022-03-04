package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"myapp/data"
	"net/http"
)

// handler to display the form
func (h *Handlers) Form(w http.ResponseWriter, r *http.Request) {
	// to pass var to jet templates, we do like follow
	vars := make(jet.VarMap)
	validator := h.App.Validator(nil)

	// pass the validator to the template using jetTemplates vars
	vars.Set("validator", validator)
	vars.Set("user", data.User{})

	err := h.App.Render.Page(w, r, "form", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// handler to submit the form and validate the datas
func (h *Handlers) PostForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	validator := h.App.Validator(nil)

	validator.Required(r, "first_name", "last_name", "email")
	validator.IsMail("email", r.Form.Get("email"))

	validator.Check(len(r.Form.Get("first_name")) > 1, "first_name", "Must be at least 2 characters")
	validator.Check(len(r.Form.Get("last_name")) > 1, "last_name", "Must be at least 2 characters")

	if !validator.Valid() {
		// we gonna take back the values to the form
		// and display errors if have some
		vars := make(jet.VarMap)
		vars.Set("validator", validator)
		var user data.User
		user.FirstName = r.Form.Get("first_name")
		user.LastName = r.Form.Get("last_name")
		user.Email = r.Form.Get("email")
		vars.Set("user", user)

		if err := h.App.Render.Page(w, r, "form", vars, nil); err != nil {
			h.App.ErrorLog.Println(err)
			return
		}
		return
	}

	fmt.Fprint(w, "valid data")
}
