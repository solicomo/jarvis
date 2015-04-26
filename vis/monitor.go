package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"jarvis"
)

const (
	SQL_NEW_METRIC_RECORD = `INSERT INTO metric_records (node, metric, value, ctime)
		SELECT ?, metrics.id, '?', datetime('now','localtime') FROM metrics WHERE name = '?';`

	SQL_NEW_NODE = `INSERT INTO nodes (addr, type) VALUES ('?', '?');`

	SQL_NEW_DEFAULT_METRICS = `INSERT INTO metric_bindings (node, metric, interval, atime, ctime) 
		SELECT ?, id, interval, datetime('now','localtime'), datetime('now','localtime') 
		FROM default_metrics;`

	SQL_UPDATE_NODE = `UPDATE nodes SET type = '?', os = '?', cpu = '?', core = '?', mem = '?',
		disk = '?', uptime = '?', atime = datetime('now','localtime') WHERE id = ?;`

	SQL_UPDATE_NODE_UPTIME = `UPDATE nodes SET uptime = '?', atime = datetime('now','localtime')
		WHERE id = ?;`

	SQL_SELECT_NODE_ID = `SELECT id FROM WHERE addr = '?';`

	SQL_SELECT_NODE_METRICS = `SELECT name, type, detector, params, md5 FROM metrics
		WHERE id IN (SELECT metric FROM metric_bindings WHERE node = ?);`

	// SQL_SELECT_NEW_METRICS = `SELECT name, type, detector, params, md5 FROM metrics
	// 	WHERE id IN (SELECT metric FROM metric_bindings AS b, nodes AS n
	// 		WHERE b.atime > n.atime AND b.node = n.id AND n.id = ?);`

	SQL_CLEAR_HISTORY = `DELETE FROM metric_records where julianday(strftime('%Y-%m-%d',datetime('now','localtime'))) - julianday(strftime('%Y-%m-%d', ctime)) > 365;`
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

	body, err := ioutil.ReadAll(r.Body)
	check(err)

	var login jarvis.Login

	err = json.Unmarshal(body, &login)
	check(err)

	var nodeID int64

	row := v.db.QueryRow(SQL_SELECT_NODE_ID, login.ListenAddr)

	if row != nil {

		err = row.Scan(nodeID)
		check(err)

	} else {

		result, e := v.db.Exec(SQL_NEW_NODE, login.ListenAddr, login.ListenType)
		check(e)

		nodeID, err = result.LastInsertId()
		check(err)

		_, err = v.db.Exec(SQL_NEW_DEFAULT_METRICS, nodeID)

		if err != nil {
			log.Println("[WARN]", err)
		}
	}

	_, err = v.db.Exec(SQL_UPDATE_NODE, login.ListenType, login.Stat.OSVer, login.Stat.CPU,
		login.Stat.Core, login.Stat.Mem, login.Stat.Disk, login.Stat.Uptime, nodeID)
	check(err)

	rows, err := v.db.Query(SQL_SELECT_NODE_METRICS, nodeID)
	check(err)

	defer rows.Close()

	var loginRsp jarvis.LoginRsp

	loginRsp.ID = strconv.FormatInt(nodeID, 10)
	loginRsp.Metrics = make(map[string]jarvis.MetricConfig)

	for rows.Next() {

		var metric jarvis.MetricConfig
		var params string

		err = rows.Scan(&metric.Name, &metric.Type, &metric.Detector, &params, &metric.MD5)
		check(err)

		if len(params) > 0 {
			err = json.Unmarshal([]byte(params), &metric.Params)
			check(err)
		}

		loginRsp.Metrics[metric.Name] = metric
	}

	check(rows.Err())

	respData, err := json.Marshal(loginRsp)
	check(err)

	if _, err = w.Write(respData); err != nil {
		log.Println("[ERRO]", "Write response failed.")
	}
}

func (v *Vis) handlePing(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	check(err)

	var ping jarvis.Ping

	err = json.Unmarshal(body, &ping)
	check(err)

	rows, err := v.db.Query(SQL_SELECT_NODE_METRICS, ping.ID)
	check(err)

	defer rows.Close()

	var pingRsp jarvis.PingRsp

	pingRsp.Metrics = make(map[string]jarvis.MetricConfig)

	for rows.Next() {

		var metric jarvis.MetricConfig
		var params string

		metric.Params = make([]interface{}, 0)

		err = rows.Scan(&metric.Name, &metric.Type, &metric.Detector, &params, &metric.MD5)
		check(err)

		if len(params) > 0 {
			err = json.Unmarshal([]byte(params), &metric.Params)
			check(err)
		}

		pingRsp.Metrics[metric.Name] = metric
	}

	check(rows.Err())

	respData, err := json.Marshal(pingRsp)
	check(err)

	if _, err = w.Write(respData); err != nil {
		log.Println("[ERRO]", "Write response failed.")
	}

	_, err = v.db.Exec(SQL_UPDATE_NODE_UPTIME, ping.Uptime, ping.ID)
	check(err)
}

func (v *Vis) handleReport(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	check(err)

	var report jarvis.MetricReport

	err = json.Unmarshal(body, &report)

	check(err)

	for name, value := range report.Metrics {

		_, err = v.db.Exec(SQL_NEW_METRIC_RECORD, report.ID, value, name)

		if err != nil {
			log.Println("[WARN]", err, string(body[:]))
		}
	}

	if _, err = w.Write([]byte(jarvis.COMMON_RSP_OK)); err != nil {
		log.Println("[WARN]", "Write response failed.")
	}
}
