package main

import (
	//"html/template"
	"mtest/controller"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"mtest/common/errors"
)

func main() {
	db := controller.TestGorp()
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("secret123"))
	store.Options(sessions.Options{
		MaxAge: 0,
	})
	m.Use(sessions.Sessions("my_session", store))
	m.Use(sessionauth.SessionUser(controller.GenerateAnonymousUser))
	sessionauth.RedirectUrl = "/auth"
	sessionauth.RedirectParam = "new-next"
	m.Use(render.Renderer(render.Options{
		Directory: "view", // Specify what path to load the templates from.
		//Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		//Extensions: []string{".html"}, // Specify extensions to load for templates.
		//Funcs:           []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		//Delims:          render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
		Charset:    "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,    // Output human readable JSON
		IndentXML:  true,    // Output human readable XML
		//HTMLContentType: "application/xhtml+xml",     // Output XHTML content type instead of default "text/html"
	}))

	//ROUTES
	m.Get("/", func(r render.Render) {
		r.Redirect("/index")
	})
	// A martini.Classic() instance automatically serves static files from the "public" directory in the root of your server.
	m.Get("/index", func(session sessions.Session, user sessionauth.User, r render.Render) {
		response := map[string]interface{}{}
		if user.IsAuthenticated() {
			response["isAuthenticated"] = true
		}
		r.HTML(200, "index", response)
	})
	//m.Get("/login", func(r render.Render) { r.HTML(200, "login", nil) })
	//m.Post("/handler", binding.Bind(controller.PostRequest{}), controller.Handler)
	m.Post("/login", binding.Bind(controller.UserAuth{}),
		func(session sessions.Session, postedUser controller.UserAuth, r render.Render, req *http.Request) {
			// You should verify credentials against a database or some other mechanism at this point.
			// Then you can authenticate this session.
			user := controller.UserAuth{}
			_, err := postedUser.CheckAuth()
			user = postedUser
			if err != nil {
				r.JSON(200, map[string]interface{}{"response": "wrong login/password"})
				return
			} else {
				err := sessionauth.AuthenticateSession(session, &user)
				if err != nil {
					r.JSON(500, err)
				}
				r.JSON(200, map[string]interface{}{"response": "ok", "name": user.Name, "last_access": user.LastAccess})
				return
			}
		})
	m.Post("/signup", binding.Bind(controller.UserSignUp{}),
		func(newUser controller.UserSignUp, r render.Render, req *http.Request) {
			// Check and save new user to database

			if newUser.Login == "" || newUser.PassHash == "" {
				r.JSON(400, errors.New("Wrong login/password"))
				return
			}
			db.Insert()
			r.JSON(200, newUser)
			//r.JSON(200, map[string]interface{}{"response": "ok"})
		})
	m.Get("/private", sessionauth.LoginRequired, func(r render.Render, user sessionauth.User) {
		r.HTML(200, "private", user.(*controller.UserAuth))
	})
	m.Get("/user", sessionauth.LoginRequired, func(r render.Render, user sessionauth.User) {
		r.HTML(200, "user", user.(*controller.UserAuth))
	})
	m.Get("/logout", sessionauth.LoginRequired, func(session sessions.Session, user sessionauth.User, r render.Render) {
		sessionauth.Logout(session, user)
		r.Redirect("/")
	})
	m.NotFound(func(r render.Render) {
		// handle 404
		r.HTML(404, "errors/404", nil)
	})
	m.RunOnAddr(":8088")
}
