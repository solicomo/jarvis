package model

const (
	SQL_SELECT_NODE_ID         = `SELECT id FROM nodes WHERE addr = ?;`
	SQL_INSERT_NODE            = `INSERT INTO nodes (name, addr, type) VALUES (?, ?, ?);`
	SQL_INSERT_DEFAULT_METRICS = `INSERT INTO metric_bindings (node, metric, interval, params, atime, ctime) 
		SELECT ?, id, interval, params, datetime('now','localtime'), datetime('now','localtime') 
		FROM default_metrics;`
)

type Node struct {
	ID     int64
	GID    int64
	Name   string
	Type   string
	Addr   string
	OS     string
	CPU    string
	Core   string
	Mem    string
	Disk   string
	Uptime string
	Ctime  string
	Atime  string
}

type Group struct {
	ID   int64
	PID  int64
	Name string
}

type Groups struct {
}

type Nodes struct {
}

var nodes Nodes
var groups Groups

func GetNodes() *Nodes {
	return &nodes
}

func GetGroups() *Groups {
	return &groups
}

func (self *Node) Save() {

}

func (self *Nodes) All() (nodes map[int64]Node, err error) {

}

func (self *Nodes) AllInGroup(gid int64) (nodes map[int64]Node, err error) {

}

func (self *Nodes) Add(t, addr string) (node Node, err error) {

	result, err := db.Exec(SQL_INSERT_NODE, addr, addr, t)

	if err != nil {
		return
	}

	node.ID, err = result.LastInsertId()

	if err != nil {
		return
	}

	db.Exec(SQL_INSERT_DEFAULT_METRICS, node.ID)

	node, err = self.Get(node.ID)

	return
}

func (self *Nodes) Rename(id int64, name string) (err error) {

}

func (self *Nodes) ChGroup(id, gid int64) (err error) {

}

func (self *Nodes) Update(id int64, os, cpu, core, mem, disk, uptime string) (err error) {

}

func (self *Nodes) Del(id int64) (err error) {

}

func (self *Nodes) Get(id int64) (node Node, err error) {

}

func (self *Nodes) GetIDFor(addr string) (id int64, err error) {
	err = self.db.QueryRow(SQL_SELECT_NODE_ID, addr).Scan(&id)
	return
}

func (self *Groups) All() {

}
