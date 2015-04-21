package jar

import (
	"encoding/json"
	"io/ioutil"
)

type Metric struct {
	Name     string
	Detector string
}

type Config struct {
	ID         string
	ListenType string
	ListenAddr string
	ServerType string
	ServerAddr string
	Metrics    map[string]string
}

type Jar struct {
	root    string
	appName string
	config  Config
}

func (j *Jar) Init(root, appName string) (err error) {
	j.root = root
	j.appName = appName

	err = j.initConfig()
	if err != nil {
		return err
	}

	return
}

func (j *Jar) initConfig() (err error) {
	configFile := j.root + "/config.json"

	configData, err := ioutil.ReadFile(configFile + "." + j.appName)

	if err != nil {
		configData, err = ioutil.ReadFile(configFile)

		if err != nil {
			return
		}
	}

	err = json.Unmarshal(configData, &j.config)

	if err != nil {
		return
	}

	configData, err = json.MarshalIndent(j.config, "", "\t")

	if err != nil {
		return
	}

	configFile += "." + j.appName
	err = ioutil.WriteFile(configFile, configData, 0644)

	return
}
