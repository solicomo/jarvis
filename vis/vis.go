package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"path"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	MonitorType string
	MonitorAddr string
	PortalType  string
	PortalAddr  string
}

type Vis struct {
	root    string
	appName string
	config  Config
	db      *sql.DB
}

func (v *Vis) Init(root, appName string) (err error) {

	v.root = root
	v.appName = appName

	err = v.initConfig()

	if err != nil {
		return err
	}

	err = v.initDB()

	return
}

func (v *Vis) initConfig() (err error) {

	configFile := path.Join(v.root, v.appName+".json")

	configData, err := ioutil.ReadFile(configFile)

	if err != nil {
		return
	}

	err = json.Unmarshal(configData, &v.config)

	return
}

func (v *Vis) initDB() (err error) {

	v.db, err = sql.Open("sqlite3", path.Join(v.root, "app/data", v.appName+".db"))

	return
}

func (v *Vis) clearHistory() {

	for _ = range time.Tick(1 * time.Hour) {
		v.db.Exec(SQL_CLEAR_HISTORY)
	}
}

func (v *Vis) Run() {

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		v.runMonitor()
	}()

	go func() {
		defer wg.Done()
		v.runPortal()
	}()

	go func() {
		defer wg.Done()
		v.clearHistory()
	}()

	wg.Wait()
}

func check(err error) {
	if err != nil {
		log.Println("[ERRO]", err)
		log.Println("[DEBU]", string(debug.Stack()[:]))
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

		log.Println("[INFO]", "request from", r.RemoteAddr, r.URL)
		fn(w, r)
	}
}

func martiniSafeHandler(layout string, h martini.Handler) martini.Handler {

	return func(req *http.Request, params martini.Params, r render.Render) {

		data := make(map[string]interface{})

		defer func() {
			if err, ok := recover().(error); ok {
				data["Status"] = "500"
				data["ErrMsg"] = err.Error()

				r.HTML(500, layout, data)
			}
		}()

		log.Println("[INFO]", "request from", req.RemoteAddr, req.URL)

		h(req, params, data)

		r.HTML(200, layout, data)
	}
}
