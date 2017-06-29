package main

import (
	"mtest/controller"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
)

func main() {
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("secret123"))
	store.Options(sessions.Options{
		MaxAge: 0,
	})
	m.Use(sessions.Sessions("my_session", store))
	m.Use(sessionauth.SessionUser(controller.GenerateAnonymousUser))
	sessionauth.RedirectUrl = "/new-login"
	sessionauth.RedirectParam = "new-next"
	m.Use(render.Renderer(render.Options{
		Directory: "view", // Specify what path to load the templates from.
		//Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".html"}, // Specify extensions to load for templates.
		//Funcs:           []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		//Delims:          render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
		Charset:    "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,    // Output human readable JSON
		IndentXML:  true,    // Output human readable XML
		//HTMLContentType: "application/xhtml+xml",     // Output XHTML content type instead of default "text/html"
	}))
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/index", func(r render.Render) {
		r.HTML(200, "index", "jeremy")
	})
	//m.Post("/handler", func(r render.Render) {
	//	r.JSON(200, map[string]interface{}{"field": "value"})
	//})
	m.Post("/handler", binding.Bind(controller.PostRequest{}), controller.Handler)
	m.Post("/auth", binding.Bind(controller.UserAuth{}))
	m.RunOnAddr(":8088")
}

//TODO: implement auth
//TODO: implement sessions and sessionauth
//TODO: create makefile
