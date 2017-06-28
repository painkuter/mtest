package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/binding"
	"mtest/controller"
)

func main() {
	m := martini.Classic()
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
	m.RunOnAddr(":8088")
}

//TODO: implement auth
//TODO: create makefile
