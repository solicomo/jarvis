package main

import (
	"path"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func (v *Vis) runPortal() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  path.Join(v.root, "app/views/simple"),
		Extensions: []string{".tmpl", ".html"},
	}))

	m.Post("/report", v.handleReport)

	m.RunOnAddr(":8081")
}
