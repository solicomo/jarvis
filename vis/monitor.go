package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"jarvis"
)

func (v *Vis) runMonitor() {

	mux := http.NewServeMux()

	mux.HandleFunc(jarvis.INDEX_URL, safeHandler(v.handleIndex))
	mux.HandleFunc(jarvis.REPORT_URL, safeHandler(v.handleReport))

	mux.Handle(jarvis.PUBLIC_URL, http.FileServer(http.Dir(jarvis.PUBLIC_DIR)))

	server := &http.Server{Addr: v.config.MonitorAddr, Handler: mux}
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

func check(err error) {
	if err != nil {
		log.Println("[ERRO]", err)
		panic(err)
	}
}

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err, ok := recover().(error); ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		fn(w, r)
	}
}

//=-=-=-=-=-=-=-=-=-=-=-=

func (v *Vis) handleIndex(w http.ResponseWriter, r *http.Request) {

}

func (v *Vis) handleReport(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	check(err)

	fmt.Println(string(body[:]))
}
