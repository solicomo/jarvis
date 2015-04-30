package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"

	"jarvis"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SQL_NEW_METRIC_RECORD = `INSERT INTO metric_records (node, metric, value, ctime)
		VALUES (?, ?, ?, datetime('now','localtime'));`

	SQL_NEW_CURRENT_METRIC = `INSERT INTO current_metrics (node, metric, value, ctime)
		VALUES (?, ?, ?, datetime('now','localtime'));`

	SQL_UPDATE_CURRENT_METRIC = `UPDATE current_metrics SET value = ?, ctime = datetime('now','localtime') 
		WHERE node = ? AND metric = ?;`

	SQL_NEW_NODE = `INSERT INTO nodes (name, addr, type) VALUES (?, ?, ?);`

	SQL_NEW_DEFAULT_METRICS = `INSERT INTO metric_bindings (node, metric, interval, params, atime, ctime) 
		SELECT ?, id, interval, params, datetime('now','localtime'), datetime('now','localtime') 
		FROM default_metrics;`

	SQL_UPDATE_NODE = `UPDATE nodes SET type = ?, os = ?, cpu = ?, core = ?, mem = ?,
		disk = ?, uptime = ?, atime = datetime('now','localtime') WHERE id = ?;`

	SQL_UPDATE_NODE_UPTIME = `UPDATE nodes SET uptime = ?, atime = datetime('now','localtime')
		WHERE id = ?;`

	SQL_SELECT_NODE_ID = `SELECT id FROM nodes WHERE addr = ?;`

	SQL_SELECT_NODE_METRICS = `SELECT m.id, m.type, m.detector, b.params, m.md5 
		FROM metrics AS m, metric_bindings AS b 
		WHERE m.id = b.metric AND b.node = ?;`

	// SQL_SELECT_NEW_METRICS = `SELECT name, type, detector, params, md5 FROM metrics
	// 	WHERE id IN (SELECT metric FROM metric_bindings AS b, nodes AS n
	// 		WHERE b.atime > n.atime AND b.node = n.id AND n.id = ?);`

	SQL_CLEAR_HISTORY = `DELETE FROM metric_records WHERE 
		julianday(datetime('now','localtime')) - julianday(ctime) > 30;`
)

func (v *Vis) runMonitor() {

	mux := http.NewServeMux()

	mux.HandleFunc(jarvis.URL_INDEX, safeHandler(v.handleIndex))
	mux.HandleFunc(jarvis.URL_LOGIN, safeHandler(v.handleLogin))
	mux.HandleFunc(jarvis.URL_PING, safeHandler(v.handlePing))
	mux.HandleFunc(jarvis.URL_REPORT, safeHandler(v.handleReport))

	mux.Handle(jarvis.URL_DETECTOR, http.FileServer(http.Dir(path.Join(v.root, jarvis.DIR_PUBLIC))))

	server := &http.Server{Addr: v.config.MonitorAddr, Handler: mux}
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
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

	switch {
	case strings.HasPrefix(login.Addr, "127.0.0.1"):
		fallthrough
	case strings.HasPrefix(login.Addr, "0.0.0.0"):
		fallthrough
	case strings.HasPrefix(login.Addr, "localhost"):
		_, op, e := net.SplitHostPort(login.Addr)

		if e == nil {
			h, _, e := net.SplitHostPort(r.RemoteAddr)

			if e == nil {
				login.Addr = h + ":" + op
			}
		}
	}

	var nodeID int64

	err = v.db.QueryRow(SQL_SELECT_NODE_ID, login.Addr).Scan(&nodeID)

	if err == sql.ErrNoRows {

		result, e := v.db.Exec(SQL_NEW_NODE, login.Addr, login.Addr, login.Type)
		check(e)

		nodeID, err = result.LastInsertId()
		check(err)

		_, err = v.db.Exec(SQL_NEW_DEFAULT_METRICS, nodeID)

		if err != nil {
			log.Println("[WARN]", err)
		}

	} else {
		check(err)
	}

	_, err = v.db.Exec(SQL_UPDATE_NODE, login.Type, login.OS, login.CPU,
		login.Core, login.Mem, login.Disk, login.Uptime, nodeID)
	check(err)

	rows, err := v.db.Query(SQL_SELECT_NODE_METRICS, nodeID)
	check(err)

	defer rows.Close()

	var loginRsp jarvis.LoginRsp

	loginRsp.ID = nodeID
	loginRsp.Metrics = make(map[string]jarvis.MetricConfig)

	for rows.Next() {

		var metric jarvis.MetricConfig
		var params string

		err = rows.Scan(&metric.ID, &metric.Type, &metric.Detector, &params, &metric.MD5)
		check(err)

		if len(params) > 0 {
			err = json.Unmarshal([]byte(params), &metric.Params)
			check(err)
		}

		loginRsp.Metrics[strconv.FormatInt(metric.ID, 10)] = metric
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

	log.Println("[DEBU]", string(body[:]))
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

		err = rows.Scan(&metric.ID, &metric.Type, &metric.Detector, &params, &metric.MD5)
		check(err)

		if len(params) > 0 {
			err = json.Unmarshal([]byte(params), &metric.Params)
			check(err)
		}

		pingRsp.Metrics[strconv.FormatInt(metric.ID, 10)] = metric
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

	for id, value := range report.Metrics {

		_, err = v.db.Exec(SQL_NEW_METRIC_RECORD, report.ID, id, value)

		if err != nil {
			log.Println("[WARN]", "new metric record:", report.ID, id, value, err)
		}

		r, err := v.db.Exec(SQL_UPDATE_CURRENT_METRIC, value, report.ID, id)

		up := true

		if err != nil {
			up = false
		} else {
			c, e := r.RowsAffected()
			if e != nil || c < 1 {
				up = false
				err = e
			}
		}

		if !up {
			log.Println("[WARN]", "update current metric record:", report.ID, id, value, err)

			_, err = v.db.Exec(SQL_NEW_CURRENT_METRIC, report.ID, id, value)

			if err != nil {
				log.Println("[WARN]", "new current metric record:", report.ID, id, value, err)
			}
		}
	}

	if _, err = w.Write([]byte(jarvis.COMMON_RSP_OK)); err != nil {
		log.Println("[WARN]", "Write response failed.")
	}
}
