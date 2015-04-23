package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory:  "app/views/simple",
		Extensions: []string{".tmpl", ".html"},
	}))

	m.Post("/report", func(req *http.Request) {
		var postData []byte
		n, err := req.Body.Read(postData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(postData[:n]))
	})

	m.RunOnAddr(":8080")
}
