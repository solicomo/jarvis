package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"jarvis/jar/detector"
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

type MetricChan struct {
	Metric chan string
	Config chan MetricConfig
}

type Stat struct {
	ID      string
	Metrics map[string]string
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

	if len(j.config.ID) == 0 {
		j.config.ID = "auto" //TODO:

		configData, err = json.MarshalIndent(j.config, "", "\t")

		if err != nil {
			return
		}

		err = ioutil.WriteFile(configFile, configData, 0644)
	}

	return
}

func (j *Jar) Run() {

	ticker := time.Tick(3 * time.Second)
	for _ = range ticker {

		metricCount := len(j.config.Metrics)

		var stat Stat
		stat.ID = j.config.ID
		stat.Metrics = make(map[string]string, metricCount)

		chans := make(map[string]MetricChan, metricCount)

		for name, config := range j.config.Metrics {
			chans[name] = MetricChan{make(chan string), make(chan MetricConfig)}

			go j.detect(chans[name].Config, chans[name].Metric)

			chans[name].Config <- config
			chans[name].Metric <- name
		}

		for name, metricChan := range chans {
			stat.Metrics[name] = <-metricChan.Metric
		}

		statData, err := json.MarshalIndent(stat, "", "\t")

		if err != nil {
			log.Println("[ERRO]", err)
			continue
		}

		_, err = http.Post(j.config.ServerType+"://"+j.config.ServerAddr+"/report", "application/json; charset=utf-8", bytes.NewReader(statData))
		_, err = os.Stderr.Write(statData)

		if err != nil {
			log.Println("[ERRO]", err)
		}
	}

}

func (j *Jar) detect(configChan chan MetricConfig, metricChan chan string) {

	metricConf := <-configChan
	metric := <-metricChan

	if strings.HasPrefix(metricConf.Detector, "call:") {
		funcName := strings.TrimPrefix(metricConf.Detector, "call:")

		var err error
		metric, err = detector.Call(funcName, metricConf.Params)
		if err != nil {
			metric = err.Error()
		}
	} else {
		metric = "Not supported yet."
	}

	metricChan <- metric
}
