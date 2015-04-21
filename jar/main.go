package main

import (
	"flag"
	"fmt"
	"jarvis/jar/jar"
	"os"
	"path"
)

func main() {
	root := flag.String("r", os.Getwd(), "Path to Jar's Home")
	flag.Parse()

	appName := path.Base(os.Args[0])

	var jar Jar
	err := jar.Init(*root, appName)

	if err != nil {
		fmt.Println(err)
		return
	}

	jar.Run()
}
