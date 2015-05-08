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
	SQL_UPDATE_NODE = `UPDATE nodes SET type = ?, os = ?, cpu = ?, core = ?, mem = ?,
		disk = ?, uptime = ?, atime = datetime('now','localtime') WHERE id = ?;`

	SQL_UPDATE_NODE_UPTIME = `UPDATE nodes SET uptime = ?, atime = datetime('now','localtime')
		WHERE id = ?;`

	SQL_SELECT_NODE_METRICS = `SELECT m.id, m.type, m.detector, b.params, m.md5 
		FROM metrics AS m, metric_bindings AS b 
		WHERE m.id = b.metric AND b.node = ?;`

	// SQL_SELECT_NEW_METRICS = `SELECT name, type, detector, params, md5 FROM metrics
	// 	WHERE id IN (SELECT metric FROM metric_bindings AS b, nodes AS n
	// 		WHERE b.atime > n.atime AND b.node = n.id AND n.id = ?);`
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

func (self *Vis) handleIndex(w http.ResponseWriter, r *http.Request) {

}

func (self *Vis) handleLogin(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	check(err)

	var login jarvis.Login

	err = json.Unmarshal(body, &login)
	check(err)

	var loginRsp jarvis.LoginRsp

	loginRsp.ID = self.getNodeID(login.Type, login.Addr, r.RemoteAddr)
	loginRsp.Metrics = self.getMetrics(loginRsp.ID)

	respData, err := json.Marshal(loginRsp)
	check(err)

	if _, err = w.Write(respData); err != nil {
		log.Println("[ERRO]", "Write response failed.")
	}

	_, err = self.db.Exec(SQL_UPDATE_NODE, login.Type, login.OS, login.CPU,
		login.Core, login.Mem, login.Disk, login.Uptime, loginRsp.ID)
	check(err)

}

func (self *Vis) handlePing(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	check(err)

	var ping jarvis.Ping

	err = json.Unmarshal(body, &ping)
	check(err)

	var pingRsp jarvis.PingRsp

	pingRsp.ID = self.getNodeID(ping.Type, ping.Addr, r.RemoteAddr)
	pingRsp.Metrics = self.getMetrics(pingRsp.ID)

	respData, err := json.Marshal(pingRsp)
	check(err)

	if _, err = w.Write(respData); err != nil {
		log.Println("[ERRO]", "Write response failed.")
	}

	_, err = self.db.Exec(SQL_UPDATE_NODE_UPTIME, ping.Uptime, ping.ID)
	check(err)
}

func (self *Vis) getNodeID(typ, addr, remote string) (node int64) {

	// node id
	switch {
	case strings.HasPrefix(addr, "127.0.0.1"):
		fallthrough
	case strings.HasPrefix(addr, "0.0.0.0"):
		fallthrough
	case strings.HasPrefix(addr, "localhost"):
		_, op, e := net.SplitHostPort(addr)

		if e == nil {
			h, _, e := net.SplitHostPort(remote)

			if e == nil {
				addr = h + ":" + op
			}
		}
	}

	nodes := model.GetNodes()
	err := nodes.GetIDFor(addr)

	if err == sql.ErrNoRows {

		n, e := nodes.Add(typ, addr)

		check(e)

		node = n.ID

	} else {
		check(err)
	}

	return
}

func (self *Vis) getMetrics(node int64) (metrics map[string]jarvis.MetricConfig) {

	// metrics
	rows, err := self.db.Query(SQL_SELECT_NODE_METRICS, node)
	check(err)

	defer rows.Close()

	for rows.Next() {

		var metric jarvis.MetricConfig
		var params string

		err = rows.Scan(&metric.ID, &metric.Type, &metric.Detector, &params, &metric.MD5)
		check(err)

		if len(params) > 0 {
			err = json.Unmarshal([]byte(params), &metric.Params)
			check(err)
		}

		metrics[strconv.FormatInt(metric.ID, 10)] = metric
	}

	check(rows.Err())

	return
}

func (self *Vis) handleReport(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	check(err)

	var report jarvis.MetricReport

	err = json.Unmarshal(body, &report)

	check(err)

	metricsRecords := model.GetMetricsRecords()

	for metric, value := range report.Metrics {

		mid, err := strconv.ParseInt(metric, 10, 0)

		if err != nil {
			log.Println("[WARN]", "report invalid metric id", metric, "from node", report.ID)
			continue
		}

		err = metricsRecords.Add(report.ID, mid, value)

		if err != nil {
			log.Println("[WARN]", "new metric record:", report.ID, id, value, err)
		}
	}

	if _, err = w.Write([]byte(jarvis.COMMON_RSP_OK)); err != nil {
		log.Println("[WARN]", "Write response failed.")
	}
}
