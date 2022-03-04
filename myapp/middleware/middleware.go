package middleware

import (
	"github.com/djedjethai/goframework"
	"myapp/data"
)

type Middleware struct {
	App    *goframework.Goframework
	Models data.Models
}
