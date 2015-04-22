package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

	ticker := time.Tick(3 * time.Minute)
	for _ = range ticker {

		var stat Stat
		stat.ID = j.config.ID
		stat.Metrics = make(map[string]Metric)

		for name, metric := range j.config.Metrics {
			stat.Metrics[name] = Metric{"", make(chan string)}

			go func(sch chan string) {
				detector := <-sch
				var value string

				//TODO:
				value = detector

				sch <- value
			}(stat.Metrics[name].Chan)

			stat.Metrics[name].Chan <- metric.Detector
		}

		for _, metric := range stat.Metrics {
			metric.Value = <-metric.Chan
			fmt.Println(metric.Value)
		}

		statData, err := json.MarshalIndent(stat, "", "\t")

		if err != nil {
			log.Println("[ERRO]", err)
			continue
		}

		_, err = http.Post(j.config.ServerType+"://"+j.config.ServerAddr+"/report", "application/json; charset=utf-8", bytes.NewReader(statData))
		//_, err = os.Stderr.Write(statData)

		if err != nil {
			log.Println("[ERRO]", err)
		}
	}

}
