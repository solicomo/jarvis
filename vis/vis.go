package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type Config struct {
	ListenType string
	ListenAddr string
}

type Vis struct {
	root    string
	appName string
	config  Config
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory:  "app/views/simple",
		Extensions: []string{".tmpl", ".html"},
	}))

	m.Post("/report", func(req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body[:]))
	})

	m.RunOnAddr(":8081")
}
