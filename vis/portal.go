package main

import (
	"path"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

const (
	SQL_SELECT_NODES_INFO = `SELECT name, type, addr, os, cpu, core, mem, disk, uptime FROM nodes WHERE addr = ?;`
)

func (v *Vis) runPortal() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  path.Join(v.root, "app/views/simple"),
		Extensions: []string{".tmpl", ".html"},
	}))

	m.Get("/dashboard/:group", martiniSafeHandler("dashboard", v.handleDashboard)

	m.RunOnAddr(v.config.PortalAddr)
}

func (v *Vis) handleDashboard(req *http.Request, params martini.Params, data *map[string]interface{}) {

	group := params["group"]

	data["Status"] = "200"
	data["Title"] = group + " | Dashboard"

	data["Nodes"] = make(map[string]interface{})
	data["Nodes"]["Info"] = make([]jarvis.NodeInfo, 1)

	rows, err := v.db.Query(SQL_SELECT_NODE_METRICS, nodeID)
	check(err)

	defer rows.Close()

	for rows.Next() {

		var info jarvis.NodeInfo

		err = rows.Scan(&info.Name, &info.Type, &info.Addr, &info.OS,
			&info.CPU, &info.Core, &info.Mem, &info.Disk, &info.Uptime)
		check(err)
		
		data["Nodes"]["Info"] = append(data["Nodes"]["Info"], info)
	}

	err = rows.Err()
	cehck(err)
}
