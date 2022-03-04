package main

import (
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"

	"github.com/djedjethai/goframework/mailer"
	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must comes before any routes
	a.use(a.Middleware.CheckRemember)

	// add routes here, all routes are call like so
	a.App.Routes.Get("/", a.Handlers.Home)
	// using the helper func to shortdown the name
	a.get("/go-page", a.Handlers.GoPage)
	a.get("/jet-page", a.Handlers.JetPage)
	a.get("/sessions", a.Handlers.SessionTest)

	a.get("/users/login", a.Handlers.UserLogin)
	a.App.Routes.Post("/users/login", a.Handlers.PostUserLogin)
	a.get("/users/logout", a.Handlers.UserLogout)
	a.get("/users/forgot-password", a.Handlers.Forgot)
	a.post("/users/forgot-password", a.Handlers.PostForgot)
	a.get("/users/reset-password", a.Handlers.ResetPasswordForm)
	a.post("/users/reset-password", a.Handlers.PostResetPassword)

	a.get("/form", a.Handlers.Form)
	a.post("/form", a.Handlers.PostForm)

	// routes to test response format handlers
	a.get("/json", a.Handlers.JSON)
	a.get("/xml", a.Handlers.XML)
	a.get("/download-file", a.Handlers.DownloadFile)

	// route to test encryption/decryption
	a.get("/crypto", a.Handlers.TestCrypto)

	// cache routes
	a.get("/cache-test", a.Handlers.ShowCachePage)
	a.post("/api/save-in-cache", a.Handlers.SaveInCache)
	a.post("/api/get-from-cache", a.Handlers.GetFromCache)
	a.post("/api/delete-from-cache", a.Handlers.DeleteFromCache)
	a.post("/api/empty-cache", a.Handlers.EmptyCache)

	// route test the mailer
	a.get("/test-mail", func(w http.ResponseWriter, r *http.Request) {
		msg := mailer.Message{
			From:        "test@example.com",
			To:          "you@there.com",
			Subject:     "Test subject - sent using func",
			Template:    "test",
			Attachments: nil,
			Data:        nil,
		}

		// that works
		// a.App.Mail.Jobs <- msg
		// res := <-a.App.Mail.Results
		// if res.Error != nil {
		// 	a.App.ErrorLog.Println(res.Error)
		// }

		// can do as well
		err := a.App.Mail.SendSMTPMessage(msg)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprint(w, "mail sent")

	})

	// test route, insert user
	a.App.Routes.Get("/create-user", func(w http.ResponseWriter, r *http.Request) {
		u := data.User{}
		u.FirstName = "jerome"
		u.LastName = "Mylastname"
		u.Email = "me@here.com"
		u.Password = "password"
		u.Active = "1"

		id, err := a.Models.Users.Insert(u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "%d %s", id, u.FirstName)
	})

	// test
	a.App.Routes.Get("/get-all-users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		for _, x := range users {
			fmt.Fprint(w, x.LastName)
		}
	})

	// test
	a.App.Routes.Get("/get-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		u, _ := a.Models.Users.Get(id)
		fmt.Fprintf(w, "%s %s %s", u.FirstName, u.LastName, u.Email)
	})

	// test
	a.App.Routes.Get("/update-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		u, _ := a.Models.Users.Get(id)

		// update name for testing purpose
		u.LastName = a.App.RandomString(10)

		validator := a.App.Validator(nil)
		// validator.Check(len(u.LastName) > 20, "last_name", "last name must have more than 10 chars")
		u.LastName = ""

		u.Validate(validator)

		if !validator.Valid() {
			fmt.Fprint(w, "failed validation")
			return
		}

		err := u.Update(*u)
		if err != nil {
			a.App.ErrorLog.Println(err)
		}

		fmt.Fprintf(w, "updated last name to %s", u.LastName)
	})

	// static routes
	// allow us to serve all the contents in the public folder
	// by going to the route what ever our server name is
	// what ever the path, the server will serve /public/myFileName.ext
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	// get the routes propertie from goFamework struct
	return a.App.Routes
}

// test database connection
// a.App.Routes.Get("/test-database", func(w http.ResponseWriter, r *http.Request) {
// 	query := "select id, first_name from users where id=1"
// 	// the row is from database/sql stdlib
// 	row := a.App.DB.Pool.QueryRowContext(r.Context(), query)

// 	var id int
// 	var name string
// 	err := row.Scan(&id, &name)
// 	if err != nil {
// 		a.App.ErrorLog.Println(err)
// 		return
// 	}

// 	fmt.Fprintf(w, "%d %s", id, name)
// })
