package model

type Group struct {
	ID   int64
	PID  int64
	Name string
}

type Groups struct {
}

var groups Groups

func GetGroups() *Groups {
	return &groups
}

func (self *Groups) All() (map[int64]Group, err error) {

}

func (self *Groups) AllInGroup(gid int64)  (map[int64]Group, err error) {

}
