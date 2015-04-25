package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"jarvis"
	"jarvis/jar/detector"
)

type Config struct {
	ID          string
	ListenType  string
	ListenAddr  string
	MonitorType string
	MonitorAddr string
	Metrics     map[string]jarvis.MetricConfig
}

type Metric struct {
	Name  string
	Value string
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

	configFile := path.Join(j.root, j.appName+".json")

	configData, err := ioutil.ReadFile(configFile)

	if err != nil {
		return
	}

	err = json.Unmarshal(configData, &j.config)

	return
}

func (j *Jar) Run() {

	ticker := time.Tick(3 * time.Second)

	for _ = range ticker {

		metricCount := len(j.config.Metrics)
		metricConfigChan := make(chan jarvis.MetricConfig, metricCount)
		metricChan := make(chan Metric, metricCount)

		for name, config := range j.config.Metrics {

			config.Name = name

			go j.detect(metricConfigChan, metricChan)

			metricConfigChan <- config
		}

		var stat jarvis.Stat
		stat.ID = j.config.ID
		stat.Metrics = make(map[string]string, metricCount)

		for metric := range metricChan {
			stat.Metrics[metric.Name] = metric.Value
		}

		statData, err := json.MarshalIndent(stat, "", "\t")

		if err != nil {
			log.Println("[ERRO]", err)
			continue
		}

		reportURL := j.config.MonitorType + "://" + j.config.MonitorAddr + jarvis.REPORT_URL

		_, err = http.Post(reportURL, "application/json; charset=utf-8", bytes.NewReader(statData))

		if err != nil {
			log.Println("[ERRO]", err)
		}
	}

}

func (j *Jar) detect(configChan chan jarvis.MetricConfig, metricChan chan Metric) {

	metricConf := <-configChan
	var metric Metric
	metric.Name = metricConf.Name

	if metricConf.Type == "call" {

		var err error
		metric.Value, err = detector.Call(metricConf.Detector, metricConf.Params)
		if err != nil {
			metric.Value = err.Error()
		}
	} else {
		metric.Value = "Not supported yet."
	}

	metricChan <- metric
}
