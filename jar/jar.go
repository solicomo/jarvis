package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"time"

	"jarvis"
	"jarvis/jar/detector"
)

type Config struct {
	ListenType  string
	ListenAddr  string
	MonitorType string
	MonitorAddr string
}

type Metric struct {
	ID    int64
	Value string
}

type Jar struct {
	ID      int64
	root    string
	appName string
	config  Config
	Metrics map[string]jarvis.MetricConfig
	Mutex   sync.RWMutex
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

	for _ = range time.Tick(10 * time.Second) {

		go j.report()
	}
}

func (j *Jar) login() (err error) {

	var logi jarvis.Login

	logi.Type = j.config.ListenType
	logi.Addr = j.config.ListenAddr

	// Can not build on Mac
	// var d detector.Detector

	// logi.OS, _ = d.OSVer("s", "r", "m")
	// logi.CPU, _ = d.CPUName()
	// logi.Core, _ = d.CPUCore()
	// logi.Mem, _ = d.MemSize()
	// logi.Disk, _ = d.DiskSize()
	// logi.Uptime, _ = d.Uptime()

	logi.OS, _ = detector.Call("OSVer", "s", "r", "m")
	logi.CPU, _ = detector.Call("CPUName")
	logi.Core, _ = detector.Call("CPUCore")
	logi.Mem, _ = detector.Call("MemSize")
	logi.Disk, _ = detector.Call("DiskSize")
	logi.Uptime, _ = detector.Call("Uptime")

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

	j.Mutex.Lock()
	j.ID = logRsp.ID
	j.Metrics = logRsp.Metrics
	j.Mutex.Unlock()

	log.Println("[INFO]", "login success, this is", logRsp.ID)
	return
}

func (j *Jar) ping() {

	for _ = range time.Tick(1 * time.Minute) {

		log.Println("[INFO]", "ping")

		var ping jarvis.Ping

		ping.ID = j.ID
		ping.Type = j.config.ListenType
		ping.Addr = j.config.ListenAddr
		ping.Uptime, _ = detector.Call("Uptime")

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

		j.Mutex.Lock()
		j.ID = pingRsp.ID
		j.Metrics = pingRsp.Metrics
		j.Mutex.Unlock()
	}
}

func (j *Jar) report() {

	log.Println("[INFO]", "report")

	j.Mutex.RLock()

	metricCount := len(j.Metrics)
	metricConfigChan := make(chan jarvis.MetricConfig, metricCount)
	metricChan := make(chan Metric, metricCount)

	for _, config := range j.Metrics {

		go j.detect(metricConfigChan, metricChan)

		metricConfigChan <- config
	}

	j.Mutex.RUnlock()

	var report jarvis.MetricReport
	report.ID = j.ID
	report.Metrics = make(map[string]string, metricCount)

	for i := 0; i < metricCount; i++ {
		metric := <-metricChan
		report.Metrics[strconv.FormatInt(metric.ID, 10)] = metric.Value
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

	if err != nil {
		log.Println("[ERRO]", err)
		return
	}
	defer r.Body.Close()

	resp, err = ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("[ERRO]", err)
	}

	if r.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("%v %v", r.StatusCode, url))
		log.Println("[ERRO]", err)
	}

	return
}

func (j *Jar) detect(configChan chan jarvis.MetricConfig, metricChan chan Metric) {

	metricConf := <-configChan

	var metric Metric
	metric.ID = metricConf.ID

	switch metricConf.Type {
	case "call":
		{
			var err error
			metric.Value, err = detector.Call(metricConf.Detector, metricConf.Params...)
			if err != nil {
				metric.Value = err.Error()
			}
		}
	case "remote":
		{
			metricConf.Detector = j.config.MonitorType + "://" + j.config.MonitorAddr + metricConf.Detector
		}
		fallthrough
	case "url":
		{
			detectorPath := path.Join(j.root, "cache/detector")
			detectorFile := path.Join(detectorPath, strconv.FormatInt(metricConf.ID, 10))

			fileData, err := ioutil.ReadFile(detectorFile)

			var h string

			if err == nil {
				sum := sha1.Sum(fileData)
				h = fmt.Sprintf("%x", sum)
			}

			if err != nil || (metricConf.MD5 != "*" && h != metricConf.MD5) {
				r, err := http.Get(metricConf.Detector)

				if err != nil {
					metric.Value = err.Error()
					log.Println("[ERRO]", "get", metricConf.Detector, err)
					break
				}
				defer r.Body.Close()

				if r.StatusCode != 200 {
					log.Println("[ERRO]", r.StatusCode, metricConf.Detector)
					break
				}

				resp, err := ioutil.ReadAll(r.Body)

				sum := sha1.Sum(resp)
				h = fmt.Sprintf("%x", sum)

				if h != metricConf.MD5 {
					metric.Value = fmt.Sprintf("file verify failed [%v:%v]: %v", h, metricConf.MD5, detectorFile)
					log.Println("[ERRO]", metric.Value)
					break
				}

				err = os.MkdirAll(detectorPath, 0755)

				if err != nil {
					metric.Value = err.Error()
					log.Println("[ERRO]", "mkdir", detectorPath, err)
					break
				}

				err = ioutil.WriteFile(detectorFile, resp, 0755)

				if err != nil {
					metric.Value = err.Error()
					log.Println("[ERRO]", "save", detectorFile, err)
					break
				}
			}

			params := make([]string, len(metricConf.Params))

			for _, p := range metricConf.Params {
				if s, ok := p.(string); ok {
					params = append(params, s)
				}
			}

			out, err := exec.Command(detectorFile, params...).Output()

			if err != nil {
				metric.Value = err.Error()
				log.Println("[ERRO]", "exec", detectorFile, err)
				break
			}

			metric.Value = string(out[:])
		}
	default:
		{
			metric.Value = "Not supported yet."
		}
	}

	metricChan <- metric
}
