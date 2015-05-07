package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	_ "github.com/mattn/go-sqlite3"

	"jarvis"
	"jarvis/vis/app/model"
)

const (
	SQL_SELECT_NODES_INFO  = `SELECT id, name, type, addr, os, cpu, core, mem, disk, uptime FROM nodes WHERE gid = ? ORDER BY id;`
	SQL_SELECT_NODE_GROUPS = `SELECT id, pid, name FROM groups ORDER by level, id;`
	SQL_SELECT_GROUPS_IN   = `SELECT id, pid, name FROM groups ORDER by level, id WHERE pid = ?;`
)

func (self *Vis) runPortal() {

	err := self.loadNodeGroups()

	if err != nil {
		log.Println("[ERRO]", "load node groups failed:", err)
	}

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  path.Join(self.root, "app/view/simple"),
		Extensions: []string{".gohtml", ".tmpl", ".html"},
	}))

	m.Get("/dashboard/:group/:gname", martiniSafeHandler("dashboard/layout", self.handleDashboardGroup))
	m.Get("/dashboard/**", martiniSafeHandler("dashboard/layout", self.handleDashboardOverviews))

	m.Get("/admin/group/:group/:gname", martiniSafeHandler("admin/nodes", self.handleAdminNodes))
	m.Get("/admin/metrics/default", martiniSafeHandler("admin/metrics", self.handleAdminMetricsDefault))
	m.Get("/admin/metrics/**", martiniSafeHandler("admin/metrics", self.handleAdminMetrics))
	m.Get("/admin/**", martiniSafeHandler("admin/nodes", self.handleAdminNodes))

	m.Get("/**", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard/overviews", http.StatusTemporaryRedirect)
	})

	m.RunOnAddr(self.config.PortalAddr)
}

func (self *Vis) handleDashboardOverviews(req *http.Request, params martini.Params, data map[string]interface{}) {

	self.nodeGroupsMutex.RLock()
	defer self.nodeGroupsMutex.RUnlock()

	data["Status"] = "200"
	data["Title"] = "Overviews | Dashboard"
	data["Groups"] = self.nodeGroups

	type Overview struct {
		Name        string
		Caption     string
		Description string
	}

	overviews := []Overview{
		{Name: "入网总数", Caption: "200万", Description: ""},
		{Name: "在线总数", Caption: "80万", Description: ""},
		{Name: "今日入网", Caption: "9000", Description: ""},
		{Name: "透传总数", Caption: "20万", Description: ""},
		{Name: "前置机数", Caption: "60台", Description: ""},
		{Name: "分中心数", Caption: "20个", Description: ""},
		{Name: "主中心数", Caption: "3个", Description: ""},
		{Name: "故障总数", Caption: "200万", Description: ""},
	}

	data["Overviews"] = overviews
}

func (v *Vis) handleDashboardGroup(req *http.Request, params martini.Params, data map[string]interface{}) {

	var group int64

	gname, _ := params["gname"]
	gid, ok := params["group"]

	group = 1

	if ok {
		group, _ = strconv.ParseInt(gid, 10, 0)
	}

	v.nodeGroupsMutex.RLock()
	defer v.nodeGroupsMutex.RUnlock()

	data["Status"] = "200"
	data["Title"] = gname + " | Dashboard"
	data["Groups"] = v.nodeGroups
	data["CurGroup"] = group
	data["CurGroupName"] = gname
	data["Nodes"] = v.loadNodesInGroup(group)
}

func (self *Vis) handleAdminNodes(req *http.Request, params martini.Params, data map[string]interface{}) {

	var group int64

	gname, _ := params["gname"]
	gid, ok := params["group"]

	group = 1

	if ok {
		group, _ = strconv.ParseInt(gid, 10, 0)
	}

	self.nodeGroupsMutex.RLock()
	defer self.nodeGroupsMutex.RUnlock()

	data["Status"] = "200"
	data["Title"] = gname + " | Nodes"
	data["Groups"] = self.nodeGroups
	data["CurGroup"] = group
	data["CurGroupName"] = gname
	data["Subs"] = self.loadSubsInGroup(group)
}

func (self *Vis) handleAdminMetrics(req *http.Request, params martini.Params, data map[string]interface{}) {

	var group int64

	gname, _ := params["gname"]
	gid, ok := params["group"]

	group = 1

	if ok {
		group, _ = strconv.ParseInt(gid, 10, 0)
	}

	self.nodeGroupsMutex.RLock()
	defer self.nodeGroupsMutex.RUnlock()

	data["Status"] = "200"
	data["Title"] = gname + " | Nodes"
	data["Groups"] = self.nodeGroups
	data["CurGroup"] = group
	data["CurGroupName"] = gname
	data["Subs"] = self.loadSubsInGroup(group)
}

func (self *Vis) handleAdminMetricsDefault(req *http.Request, params martini.Params, data map[string]interface{}) {

	var group int64

	gname, _ := params["gname"]
	gid, ok := params["group"]

	group = 1

	if ok {
		group, _ = strconv.ParseInt(gid, 10, 0)
	}

	self.nodeGroupsMutex.RLock()
	defer self.nodeGroupsMutex.RUnlock()

	data["Status"] = "200"
	data["Title"] = gname + " | Nodes"
	data["Groups"] = self.nodeGroups
	data["CurGroup"] = group
	data["CurGroupName"] = gname
	data["Subs"] = self.loadSubsInGroup(group)
}

func (self *Vis) loadNodeGroups() (err error) {

	rows, err := self.db.Query(SQL_SELECT_NODE_GROUPS)

	if err != nil {
		return
	}

	defer rows.Close()

	self.nodeGroupsMutex.Lock()
	defer self.nodeGroupsMutex.Unlock()

	self.nodeGroups = make(map[int64]NodeGroup)

	for rows.Next() {

		var group NodeGroup

		err = rows.Scan(&group.ID, &group.PID, &group.Name)
		if err != nil {
			return
		}

		// 当前只支持 2 级
		if group.PID == 0 {
			group.Subs = make(map[int64]NodeGroup)
			self.nodeGroups[group.ID] = group
		} else {

			g, ok := self.nodeGroups[group.PID]
			if ok {
				g.Subs[group.ID] = group
			}
		}
	}

	err = rows.Err()

	//
	ungroup, ok := self.nodeGroups[1]
	if ok {
		if len(ungroup.Subs) == 0 {
			delete(self.nodeGroups, 1)
		}
	}

	return
}

func (v *Vis) loadNodesInGroup(group int64) interface{} {

	type Metric struct {
		ID     int64
		Name   string
		Value  string
		Values map[string]string
	}

	type Node struct {
		Info    jarvis.NodeInfo
		Metrics map[int64]Metric
	}

	nodes := make(map[int64]Node)

	rows, err := v.db.Query(SQL_SELECT_NODES_INFO, group)
	check(err)

	defer rows.Close()

	for rows.Next() {

		var node Node

		err = rows.Scan(&node.Info.ID, &node.Info.Name, &node.Info.Type, &node.Info.Addr, &node.Info.OS,
			&node.Info.CPU, &node.Info.Core, &node.Info.Mem, &node.Info.Disk, &node.Info.Uptime)
		check(err)

		nodes[node.Info.ID] = node
	}

	err = rows.Err()
	check(err)

	metricsRecords := model.GetMetricsRecords()

	for id, node := range nodes {

		records, err := metricsRecords.CurrentFor(id)
		check(err)

		node.Metrics = make(map[int64]Metric)

		for mid, m := range records {

			metric := Metric{ID: mid, Name: m.Name, Value: m.Value}

			metric.Values = make(map[string]string)
			json.Unmarshal([]byte(metric.Value), &metric.Values)

			node.Metrics[metric.ID] = metric
		}

		nodes[id] = node
	}

	return nodes
}

func (self *Vis) loadSubsInGroup(group int64) interface{} {

	type Node struct {
		ID    int64
		GID   int64
		Name  string
		Addr  string
		Atime string
	}

	type Sub struct {
		ID    int64
		PID   int64
		Name  string
		Nodes map[int64]Node
	}

	subs := make(map[int64]Sub)

	return subs
}
