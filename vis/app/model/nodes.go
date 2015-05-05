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
