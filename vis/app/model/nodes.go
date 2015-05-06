package model

var nodes Nodes
var groups Groups

func GetNodes() *Nodes {
	return &nodes
}

func GetGroups() *Groups {
	return &groups
}

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

func (self *Node) Save() {

}

func (self *Nodes) All() (nodes map[int64]Node, err error) {

}

func (self *Nodes) AllInGroup(gid int64) (nodes map[int64]Node, err error) {

}

func (self *Nodes) Add(gid int64, name, t, addr string) (node Node, err error) {

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

func (self *Groups) All() {

}
