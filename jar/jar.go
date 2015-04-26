package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"sync"
	"time"

	"jarvis"
	"jarvis/jar/detector"
)

type Config struct {
	ID           string
	ListenType   string
	ListenAddr   string
	MonitorType  string
	MonitorAddr  string
	MetricsMutex sync.RWMutex
	Metrics      map[string]jarvis.MetricConfig
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

	for {
		err := j.login()

		if err == nil {
			break
		}

		time.Sleep(10 * time.Second)
	}

	go j.ping()

	for _ = range time.Tick(10 * time.Minute) {

		go j.report()
	}
}

func (j *Jar) login() (err error) {

	var logi jarvis.Login

	logi.ListenType = j.config.ListenType
	logi.ListenAddr = j.config.ListenAddr

	// Can not build on Mac
	// var d detector.Detector

	// logi.Stat.OSVer, _ = d.OSVer([]string{"s", "r", "m", "n"})
	// logi.Stat.CPU, _ = d.CPUName()
	// logi.Stat.Core, _ = d.CPUCore()
	// logi.Stat.Mem, _ = d.MemSize()
	// logi.Stat.Disk, _ = d.DiskSize()
	// logi.Stat.Uptime, _ = d.Uptime()

	logi.Stat.OSVer, _ = detector.Call("OSVer", []interface{}{"s", "r", "m", "n"})
	logi.Stat.CPU, _ = detector.Call("CPUName", []interface{}{})
	logi.Stat.Core, _ = detector.Call("CPUCore", []interface{}{})
	logi.Stat.Mem, _ = detector.Call("MemSize", []interface{}{})
	logi.Stat.Disk, _ = detector.Call("DiskSize", []interface{}{})
	logi.Stat.Uptime, _ = detector.Call("Uptime", []interface{}{})

	resp, err := j.postTo(jarvis.URL_LOGIN, logi)

	if err != nil {
		return
	}

	var logRsp jarvis.LoginRsp

	err = json.Unmarshal(resp, &logRsp)

	if err != nil {
		log.Println("[ERRO]", err)
		return
	}

	j.config.ID = logRsp.ID
	j.config.Metrics = logRsp.Metrics

	return
}

func (j *Jar) ping() {

	for _ = range time.Tick(1 * time.Minute) {

		var ping jarvis.Ping

		ping.ID = j.config.ID
		ping.Uptime, _ = detector.Call("Uptime", []interface{}{})

		resp, err := j.postTo(jarvis.URL_PING, ping)

		if err != nil {
			continue
		}

		var pingRsp jarvis.PingRsp

		err = json.Unmarshal(resp, &pingRsp)

		if err != nil {
			log.Println("[ERRO]", err)
			continue
		}

		j.config.MetricsMutex.Lock()
		j.config.Metrics = pingRsp.Metrics
		j.config.MetricsMutex.Unlock()
	}
}

func (j *Jar) report() {

	j.config.MetricsMutex.RLock()

	metricCount := len(j.config.Metrics)
	metricConfigChan := make(chan jarvis.MetricConfig, metricCount)
	metricChan := make(chan Metric, metricCount)

	for name, config := range j.config.Metrics {

		config.Name = name

		go j.detect(metricConfigChan, metricChan)

		metricConfigChan <- config
	}

	j.config.MetricsMutex.RUnlock()

	var report jarvis.MetricReport
	report.ID = j.config.ID
	report.Metrics = make(map[string]string, metricCount)

	for metric := range metricChan {
		report.Metrics[metric.Name] = metric.Value
	}

	j.postTo(jarvis.URL_REPORT, report)
}

func (j *Jar) postTo(url string, data interface{}) (resp []byte, err error) {

	postData, err := json.Marshal(data)

	if err != nil {
		log.Println("[ERRO]", err)
		return
	}

	postURL := j.config.MonitorType + "://" + j.config.MonitorAddr + url

	r, err := http.Post(postURL, "application/json; charset=utf-8", bytes.NewReader(postData))

	defer r.Body.Close()

	if err != nil {
		log.Println("[ERRO]", err)
		return
	}

	resp, err = ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("[ERRO]", err)
	}

	return
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
