package model

const (
	SQL_SELECT_NODE_ID         = `SELECT id FROM nodes WHERE addr = ?;`
	SQL_SELECT_NODE_BY_ID      = `SELECT id,gid,name,addr,type,os,cpu,core,mem,disk,uptime,ctime,atime FROM nodes WHERE id = ?;`
	SQL_SELECT_NODES_IN_GROUP  = `SELECT id,gid,name,addr,type,os,cpu,core,mem,disk,uptime,ctime,atime FROM nodes WHERE gid = ?;`
	SQL_SELECT_ALL_NODES       = `SELECT id,gid,name,addr,type,os,cpu,core,mem,disk,uptime,ctime,atime FROM nodes;`
	SQL_DELETE_NODE            = `DELETE nodes  WHERE id = ?;`
	SQL_UPDATE_NODE_NAME       = `UPDATE nodes SET name = ? WHERE id = ?;`
	SQL_UPDATE_NODE_GROUP      = `UPDATE nodes SET gid = ? WHERE id = ?;`
	SQL_UPDATE_NODE            = `UPDATE nodes SET gid=?,name=?,addr=?,type=?,os=?,cpu=?,core=?,mem=?,disk=?,uptime=?,atime = datetime('now','localtime') WHERE id = ?;`
	SQL_INSERT_NODE            = `INSERT INTO nodes (gid,name,addr,type,os,cpu,core,mem,disk,uptime,ctime,atime) VALUES (?,?,?,?,?,?,?,?,?,?,datetime('now','localtime'),datetime('now','localtime'));`
	SQL_INSERT_DEFAULT_METRICS = `INSERT INTO metric_bindings (node, metric, interval, params, atime, ctime) 
		SELECT ?, id, interval, params, datetime('now','localtime'), datetime('now','localtime') 
		FROM default_metrics;`
)

type Node struct {
	ID     int64
	GID    int64
	Name   string
	Addr   string
	Type   string
	OS     string
	CPU    string
	Core   string
	Mem    string
	Disk   string
	Uptime string
	Ctime  string
	Atime  string
}

type Nodes struct {
}

var nodes Nodes

func GetNodes() *Nodes {
	return &nodes
}

func (self *Node) Save() (err error) {

	// update
	if self.ID != 0 {
		_, err = db.Exec(SQL_UPDATE_NODE,
			self.GID,
			self.Name,
			self.Addr,
			self.Type,
			self.OS,
			self.CPU,
			self.Core,
			self.Mem,
			self.Disk,
			self.Uptime,
			self.ID)

		return
	}

	// insert
	result, err := db.Exec(SQL_INSERT_NODE,
		self.GID,
		self.Name,
		self.Addr,
		self.Type,
		self.OS,
		self.CPU,
		self.Core,
		self.Mem,
		self.Disk,
		self.Uptime)

	if err != nil {
		return
	}

	self.ID, err = result.LastInsertId()

	if err != nil {
		return
	}

	db.Exec(SQL_INSERT_DEFAULT_METRICS, self.ID)

	return
}

func (self *Nodes) All() (nodes map[int64]Node, err error) {

	nodes = make(map[int64]Node)

	rows, err := db.Query(SQL_SELECT_ALL_NOTES)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {

		var node Node

		err = rows.Scan(
			&node.ID,
			&node.GID,
			&node.Name,
			&node.Addr,
			&node.Type,
			&node.OS,
			&node.CPU,
			&node.Core,
			&node.Mem,
			&node.Disk,
			&node.Uptime,
			&node.Ctime,
			&node.Atime)

		if err != nil {
			return
		}

		nodes[node.ID] = node
	}

	err = rows.Err()

	return
}

func (self *Nodes) AllInGroup(gid int64) (nodes map[int64]Node, err error) {

	nodes = make(map[int64]Node)

	rows, err := db.Query(SQL_SELECT_NOTES_IN_GROUP, gid)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {

		var node Node

		err = rows.Scan(
			&node.ID,
			&node.GID,
			&node.Name,
			&node.Addr,
			&node.Type,
			&node.OS,
			&node.CPU,
			&node.Core,
			&node.Mem,
			&node.Disk,
			&node.Uptime,
			&node.Ctime,
			&node.Atime)

		if err != nil {
			return
		}

		nodes[node.ID] = node
	}

	err = rows.Err()

	return
}

func (self *Nodes) Add(t, addr string) (node Node, err error) {

	node.Addr = addr
	node.Type = t

	err = node.Save()

	return
}

func (self *Nodes) Rename(id int64, name string) (err error) {

	result, err = db.Exec(SQL_UPDATE_NODE_NAME, name, id)

	if err != nil {
		return
	}

	c, err := result.RowsAffected()

	if err != nil {
		return
	}

	if c < 1 {
		err = errors.New("No such node: " + strconv.FormatInt(id, 10))
	}

	return
}

func (self *Nodes) ChGroup(id, gid int64) (err error) {

	result, err = db.Exec(SQL_UPDATE_NODE_GROUP, gid, id)

	if err != nil {
		return
	}

	c, err := result.RowsAffected()

	if err != nil {
		return
	}

	if c < 1 {
		err = errors.New("No such node: " + strconv.FormatInt(id, 10))
	}

	return
}

func (self *Nodes) Update(id int64, os, cpu, core, mem, disk, uptime string) (err error) {

	node, err := self.Get(id)

	if err != nil {
		return
	}

	node.OS = os
	node.CPU = cpu
	node.Core = core
	node.Mem = mem
	node.Disk = disk
	node.Uptime = uptime

	err = node.Save()
}

func (self *Nodes) Del(id int64) (err error) {

	_, err = db.Exec(SQL_DELETE_NODE, id)
}

func (self *Nodes) Get(id int64) (node Node, err error) {

	err = self.db.QueryRow(SQL_SELECT_NODE_BY_ID, id).Scan(
		&node.ID,
		&node.GID,
		&node.Name,
		&node.Addr,
		&node.Type,
		&node.OS,
		&node.CPU,
		&node.Core,
		&node.Mem,
		&node.Disk,
		&node.Uptime,
		&node.Ctime,
		&node.Atime)

	return
}

func (self *Nodes) GetIDFor(addr string) (id int64, err error) {

	err = self.db.QueryRow(SQL_SELECT_NODE_ID, addr).Scan(&id)
	return
}
