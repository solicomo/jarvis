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
	SQL_SELECT_NODES_INFO = `SELECT name, type, addr, os, cpu, core, mem, disk, uptime FROM nodes;`
)

func (v *Vis) runPortal() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  path.Join(v.root, "app/views/simple"),
		Extensions: []string{".gohtml", ".tmpl", ".html"},
	}))

	m.Get("/dashboard/:group", martiniSafeHandler("dashboard", v.handleDashboard))

	m.RunOnAddr(v.config.PortalAddr)
}

func (v *Vis) handleDashboard(req *http.Request, params martini.Params, data map[string]interface{}) {

	group := params["group"]

	data["Status"] = "200"
	data["Title"] = group + " | Dashboard"

	nodesInfo := make([]jarvis.NodeInfo, 0)

	rows, err := v.db.Query(SQL_SELECT_NODES_INFO)
	check(err)

	defer rows.Close()

	for rows.Next() {

		var info jarvis.NodeInfo

		err = rows.Scan(&info.Name, &info.Type, &info.Addr, &info.OS,
			&info.CPU, &info.Core, &info.Mem, &info.Disk, &info.Uptime)
		check(err)

		nodesInfo = append(nodesInfo, info)
	}

	err = rows.Err()
	check(err)

	nodes := make(map[string]interface{})
	nodes["Info"] = nodesInfo
	data["Nodes"] = nodes
}
