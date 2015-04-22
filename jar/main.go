package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	// path and name
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	root := flag.String("r", wd, "Path to Jar's Home")
	flag.Parse()

	appName := path.Base(os.Args[0])

	// log
	logFile := path.Clean(*root + "/" + appName + ".log")
	logWriter, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer logWriter.Close()

	log.SetOutput(logWriter)
	log.SetFlags(log.LstdFlags)

	var jar Jar
	err = jar.Init(*root, appName)

	if err != nil {
		log.Println("[ERRO]", err)
		return
	}

	jar.Run()
}
