package main

import (
	"net/http"
	"path"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	_ "github.com/mattn/go-sqlite3"

	"jarvis"
)

const (
	SQL_SELECT_NODES_INFO  = `SELECT name, type, addr, os, cpu, core, mem, disk, uptime FROM nodes WHERE group = ?;`
	SQL_SELECT_NODE_GROUPS = `SELECT id, pid, name FROM groups ORDER by level, id;`
)

func (v *Vis) runPortal() {

	err := v.loadNodeGroups()

	if err != nil {
		log.Println("[ERRO]", "load node groups failed:", err)
	}

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  path.Join(v.root, "app/views/simple"),
		Extensions: []string{".gohtml", ".tmpl", ".html"},
	}))

	m.Get("/dashboard", martiniSafeHandler("dashboard", v.handleDashboardIndex))
	m.Get("/dashboard/:group/:gname", martiniSafeHandler("dashboard", v.handleDashboardGroup))

	m.RunOnAddr(v.config.PortalAddr)
}

func (v *Vis) loadNodeGroups() (err error) {

	rows, err := v.db.Query(SQL_SELECT_NODE_GROUPS)

	if err != nil {
		return
	}

	defer rows.Close()

	v.nodeGroups = make(map[int64]NodeGroup)

	for rows.Next() {

		var group NodeGroup

		err = rows.Scan(&group.ID, &group.PID, &group.Name)
		if err != nil {
			return
		}

		// 当前只支持 2 级
		if group.PID == 0 {
			v.nodeGroups[group.ID] = group
			v.nodeGroups[group.ID].Subs = make(map[int64]NodeGroup)
		} else {

			subs, ok := v.nodeGroups[group.PID]
			if ok {
				subs[group.ID] = group
			}
		}
	}

	err = rows.Err()
}

func (v *Vis) handleDashboardIndex(req *http.Request, params martini.Params, data map[string]interface{}) {

	data["Status"] = "200"
	data["Title"] = "Dashboard"
	data["Groups"] = v.nodeGroups

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

	data["Status"] = "200"
	data["Title"] = gname + " | Dashboard"
	data["Groups"] = v.nodeGroups
	data["CurSubGroup"] = group

	pg := "1"
	for id, g := range v.nodeGroups {
		if _, ok := g.Subs[group]; ok {
			pg = id
			break
		}
	}

	data["CurGroup"] = pg

	type Node struct {
		Info    jarvis.NodeInfo
		Metrics interface{}
	}

	nodes := make([]Node, 0)

	rows, err := v.db.Query(SQL_SELECT_NODES_INFO, group)
	check(err)

	defer rows.Close()

	for rows.Next() {

		var node Node

		err = rows.Scan(&node.Info.Name, &node.Info.Type, &node.Info.Addr, &node.Info.OS,
			&node.Info.CPU, &node.Info.Core, &node.Info.Mem, &node.Info.Disk, &node.Info.Uptime)
		check(err)

		nodes = append(nodes, node)
	}

	err = rows.Err()
	check(err)

	data["Nodes"] = nodes
}
