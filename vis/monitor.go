package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"jarvis"
)

const (
	SQL_NEW_METRIC_RECORD = `INSERT INTO metric_records VALUES (node, metric, value, ctime)
		SELECT ?, metrics.id, ?, datetime('now','localtime') FROM metrics WHERE name = '?';`
)

func (v *Vis) runMonitor() {

	mux := http.NewServeMux()

	mux.HandleFunc(jarvis.URL_INDEX, safeHandler(v.handleIndex))
	mux.HandleFunc(jarvis.URL_LOGIN, safeHandler(v.handleLogin))
	mux.HandleFunc(jarvis.URL_PING, safeHandler(v.handlePing))
	mux.HandleFunc(jarvis.URL_REPORT, safeHandler(v.handleReport))

	mux.Handle(jarvis.URL_PUBLIC, http.FileServer(http.Dir(jarvis.DIR_PUBLIC)))

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

func (v *Vis) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (v *Vis) handlePing(w http.ResponseWriter, r *http.Request) {

}

func (v *Vis) handleReport(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	check(err)

	var report jarvis.MetricReport

	err = json.Unmarshal(body, &report)

	check(err)

	for name, value := range report.Metrics {

		_, err := v.db.Query(SQL_NEW_METRIC_RECORD, report.ID, value, name)

		if err != nil {
			log.Println("[WARN]", err, string(body[:]))
		}
	}

	if _, err = w.Write([]byte(jarvis.COMMON_RSP_OK)); err != nil {
		log.Println("[WARN]", "Write response failed.")
	}
}
