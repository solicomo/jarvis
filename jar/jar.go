package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type MetricConfig struct {
	Detector string
	Params   []interface{}
	MD5      string
}

type Config struct {
	ID         string
	ListenType string
	ListenAddr string
	ServerType string
	ServerAddr string
	Metrics    map[string]MetricConfig
}

type Metric struct {
	Value string
	Chan  chan string `json:"-"`
}

type Stat struct {
	ID      string
	Metrics map[string]Metric
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

	configData, err := ioutil.ReadFile(configFile)

	if err != nil {
		return
	}

	err = json.Unmarshal(configData, &j.config)

	if err != nil {
		return
	}

	if j.config.ID == "auto" {

	}

	configData, err = json.MarshalIndent(j.config, "", "\t")

	if err != nil {
		return
	}

	err = ioutil.WriteFile(configFile, configData, 0644)

	return
}

func (j *Jar) Run() {

	for {
		var stat Stat
		stat.ID = j.config.ID
		stat.Metrics = make(map[string]Metric)

		for name, metric := range j.config.Metrics {
			stat.Metrics[name].Chan = make(chan string)

			go func(sch chan string) {
				detector := <-sch
				var value string
				//TODO:
				value = "val"
				sch <- value
			}(stat.Metrics[name].Chan)

			stat.Metrics[name].Chan <- metric.Detector
		}

		for name, metric := range stat.Metrics {
			metric.Value = <-metric.Chan
		}

		statData, err := json.MarshalIndent(stat, "", "\t")

		if err != nil {
			log.Println("[ERRO]", err)
			continue
		}

		_, err = http.Post(j.config.ServerType+"://"+j.config.ServerAddr+"/report", "application/json; charset=utf-8", bytes.NewReader(statData))

		if err != nil {
			log.Println("[ERRO]", err)
		}
	}

}
