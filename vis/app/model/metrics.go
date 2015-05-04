package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var metrics Metrics

func InitMetrics() {
	metrics.init()
}

func GetMetrics() (metrics *Metrics) {
	return &metrics
}

const (
	SQL_SELECT_METRICS         = `SELECT id, name, type, detector, md5 FROM metrics;`
	SQL_SELECT_METRICS_DEFAULT = `SELECT id, name, params, interval
	                              FROM default_metrics AS d, metrics AS m
	                              WHERE d.id = m.id;`
	SQL_INSERT_METRIC      = `INSERT INTO metrics (name, type, detector, md5) VALUES(?,?,?,?);`
	SQL_UPDATE_METRIC_NAME = `UPDATE metrics SET name = ? WHERE id = ?;`
	SQL_UPDATE_METRIC      = `UPDATE metrics SET type = ?, detector = ?, md5 = ? WHERE id = ?;`
)

type Metric struct {
	ID       int64
	Name     string
	Type     string
	Detector string
	MD5      string
}

type DefaultMetric struct {
	ID       int64
	Name     string
	Params   string
	Interval int
}

type Metrics struct {
	metrics  map[int64]Metric
	defaults map[int64]DefaultMetric
	db       *sql.DB
}

func (self *Metrics) init() {

	self.metrics = make(map[int64]Metric)
	self.defaults = make(map[int64]DefaultMetric)

	rows, err := self.db.Query(SQL_SELECT_METRICS)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {

		var metric Metric

		err = rows.Scan(&metric.ID, &metric.Name, &metric.Type, &metric.Detector, &metric.MD5)

		if err != nil {
			panic(err)
		}

		self.metrics[metric.ID] = metric
	}

	err = rows.Err()

	if err != nil {
		panic(err)
	}

	rows, err = self.db.Query(SQL_SELECT_METRICS_DEFAULT)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {

		var metric DefaultMetric

		err = rows.Scan(&metric.ID, &metric.Name, &metric.Params, &metric.Interval)

		if err != nil {
			panic(err)
		}

		self.defaults[metric.ID] = metric
	}

	err = rows.Err()

	if err != nil {
		panic(err)
	}
}

func (self *Metrics) All() (metrics map[int64]Metric) {
	return self.metrics
}

func (self *Metrics) AllDefault() (dfts map[int64]DefaultMetric) {
	return self.defaults
}

func (self *Metrics) AddMetric(name, t, detector, md5 string) (metric Metric, err error) {

	metric = Metric{Name: name, Type: t, Detector: detector, MD5: md5}

	result, err := self.db.Exec(SQL_INSERT_METRIC, name, t, detector, md5)

	if err != nil {
		return
	}

	metric.ID, err = result.LastInsertId()
	return
}

func (self *Metrics) RenameMetric(id int64, name string) (err error) {

	result, err := self.db.Exec(SQL_UPDATE_METRIC_NAME, name, id)

	if err != nil {
		return
	}

	c, err := r.RowsAffected()

	if err != nil {
		return
	}

	if c < 1 {
		err = errors.New("No such metric: " + strconv.FormatInt(id, 10))
	}

	return
}

func (self *Metrics) UpdateMetric(id int64, t, detector, md5 string) (err error) {

	result, err := self.db.Exec(SQL_UPDATE_METRIC, t, detector, md5, id)

	if err != nil {
		return
	}

	c, err := r.RowsAffected()

	if err != nil {
		return
	}

	if c < 1 {
		err = errors.New("No such metric: " + strconv.FormatInt(id, 10))
	}

	return
}
