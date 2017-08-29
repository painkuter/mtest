package main

import (
	"html/template"
	"mtest/controller"
	"net/http"

	"fmt"

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
	sessionauth.RedirectUrl = "/auth"
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
	m.Get("/index", func(res http.ResponseWriter) {
		//if user.IsAuthenticated(){
		//
		//}
		t, _ := template.New("test").ParseFiles("view/test.html")
		t.Execute(res, nil)
		//r.HTML(200, "index", map[string]interface{}{"name": "wor1ld"})
	})
	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", "jeremy")
	})
	//m.Post("/handler", func(r render.Render) {
	//	r.JSON(200, map[string]interface{}{"field": "value"})
	//})
	m.Post("/handler", binding.Bind(controller.PostRequest{}), controller.Handler)
	m.Post("/auth", binding.Bind(controller.UserAuth{}),
		func(session sessions.Session, postedUser controller.UserAuth, r render.Render, req *http.Request) {
			fmt.Println("AUTH!!!")
			// You should verify credentials against a database or some other mechanism at this point.
			// Then you can authenticate this session.
			//r.JSON(200, map[string]interface{}{"field": "auth"})
			user := controller.UserAuth{}
			name, err := postedUser.CheckAuth()
			user.Name = name
			if err != nil {
				//wrong login/pass
				fmt.Println("WRONG LOGIN")
				//r.Redirect("/user")

				r.HTML(200, "logged-in", map[string]interface{}{"name": user.Name, "time": "10-00"})
				return
			} else {
				fmt.Println(user)
				err := sessionauth.AuthenticateSession(session, &user)
				fmt.Println(user)
				if err != nil {
					r.JSON(500, err)
				}

				//params := req.URL.Query()
				//redirect := params.Get(sessionauth.RedirectParam)
				//r.Redirect(redirect)
				r.Redirect("/user")
				fmt.Println("RIGHT LOGIN")

				return
			}
		})
	m.Post("/private", sessionauth.LoginRequired, func(r render.Render, user sessionauth.User) {
		r.HTML(200, "private", user.(*controller.UserAuth))
	})
	m.Get("/user", sessionauth.LoginRequired, func(r render.Render, user sessionauth.User) {
		//fmt.Println(user.GetById())
		r.HTML(200, "user", user.(*controller.UserAuth))
	})
	m.Get("/logout", sessionauth.LoginRequired, func(session sessions.Session, user sessionauth.User, r render.Render) {
		sessionauth.Logout(session, user)
		r.Redirect("/")
	})
	m.RunOnAddr(":8088")
}

//TODO: implement auth
//TODO: implement sessions and sessionauth
//TODO: +create login form
//TODO: create makefile
