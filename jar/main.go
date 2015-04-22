package main

import (
	"flag"
	"fmt"
	"log"
	"log/syslog"
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
	// Configure logger to write to the syslog.
	logWriter, err := syslog.New(syslog.LOG_NOTICE, appName)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer logWriter.Close()

	log.SetOutput(logWriter)
	log.SetFlags(0)

	var jar Jar
	err = jar.Init(*root, appName)

	if err != nil {
		log.Println("[ERRO]", err)
		return
	}

	jar.Run()
}
