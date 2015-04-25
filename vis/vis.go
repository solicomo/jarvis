package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"sync"
)

type Config struct {
	MonitorType string
	MonitorAddr string
	PortalType  string
	PortalAddr  string
}

type Vis struct {
	root    string
	appName string
	config  Config
}

func (v *Vis) Init(root, appName string) (err error) {

	v.root = root
	v.appName = appName

	err = v.initConfig()
	if err != nil {
		return err
	}

	return
}

func (v *Vis) initConfig() (err error) {

	configFile := path.Join(v.root, v.appName+".json")

	configData, err := ioutil.ReadFile(configFile)

	if err != nil {
		return
	}

	err = json.Unmarshal(configData, &v.config)

	return
}

func (v *Vis) Run() {

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		v.runMonitor()
	}()

	go func() {
		defer wg.Done()
		v.runPortal()
	}()

	wg.Wait()
}
