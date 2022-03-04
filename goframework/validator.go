package goframework

import (
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Validation struct {
	Data   url.Values // map[string][]string
	Errors map[string]string
}

func (g *Goframework) Validator(data url.Values) *Validation {
	return &Validation{
		Errors: make(map[string]string),
		Data:   data,
	}
}

func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validation) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// check if any value is in the form post
func (v *Validation) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// specific for form post
func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "this field can not be blank")
		}

	}
}

// useful to check any kind of validation err anywhere in the app
func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validation) IsMail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "Invalid email address")
	}
}

func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "This field must be an integer")
	}
}

func (v *Validation) IsFloat(field, value string) {
	_, err := strconv.ParseFloat(value, 64) // parse float64 (but we do not care)
	if err != nil {
		v.AddError(field, "This field must be an floating point number")
	}
}

func (v *Validation) IsDateISO(field, value string) {
	_, err := time.Parse("2006-01-02", value) // "2006-01-02" unusual way to refer to ISO
	if err != nil {
		v.AddError(field, "This field must be a date in the form of YYYY-MM-DD")
	}
}

func (v *Validation) NoSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "Spaces are not permitted")
	}
}
