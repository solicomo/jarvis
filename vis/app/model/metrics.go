package model

import (
	"errors"
	"strconv"
)

const (
	SQL_SELECT_METRICS      = `SELECT id, name, type, detector, md5 FROM metrics;`
	SQL_SELECT_METRIC_BY_ID = `SELECT id, name, type, detector, md5 FROM metrics WHERE id = ?;`
	SQL_SELECT_METRICS_DFT  = `SELECT id, name, params, interval
	                            FROM default_metrics AS d, metrics AS m
	                           WHERE d.id = m.id;`
	SQL_INSERT_METRIC      = `INSERT INTO metrics (name, type, detector, md5) VALUES(?,?,?,?);`
	SQL_UPDATE_METRIC_NAME = `UPDATE metrics SET name = ? WHERE id = ?;`
	SQL_UPDATE_METRIC      = `UPDATE metrics SET type = ?, detector = ?, md5 = ? WHERE id = ?;`
	SQL_DELETE_METRIC      = `DELETE FROM metrics WHERE id = ?;`
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
}

var metrics Metrics

func GetMetrics() *Metrics {
	return &metrics
}

func (self *Metric) Save() {

}

func (self *Metrics) All() (metrics map[int64]Metric, err error) {

	metrics = make(map[int64]Metric)

	rows, err := db.Query(SQL_SELECT_METRICS)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {

		var metric Metric

		err = rows.Scan(&metric.ID, &metric.Name, &metric.Type, &metric.Detector, &metric.MD5)

		if err != nil {
			return
		}

		metrics[metric.ID] = metric
	}

	err = rows.Err()

	return
}

func (self *Metrics) AllDefault() (dfts map[int64]DefaultMetric, err error) {

	defaults = make(map[int64]DefaultMetric)

	rows, err = db.Query(SQL_SELECT_METRICS_DEFAULT)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {

		var metric DefaultMetric

		err = rows.Scan(&metric.ID, &metric.Name, &metric.Params, &metric.Interval)

		if err != nil {
			return
		}

		self.defaults[metric.ID] = metric
	}

	err = rows.Err()

	return
}

func (self *Metrics) Add(name, t, detector, md5 string) (metric Metric, err error) {

	metric = Metric{Name: name, Type: t, Detector: detector, MD5: md5}

	result, err := db.Exec(SQL_INSERT_METRIC, name, t, detector, md5)

	if err != nil {
		return
	}

	metric.ID, err = result.LastInsertId()
	return
}

func (self *Metrics) Rename(id int64, name string) (err error) {

	result, err := db.Exec(SQL_UPDATE_METRIC_NAME, name, id)

	if err != nil {
		return
	}

	c, err := result.RowsAffected()

	if err != nil {
		return
	}

	if c < 1 {
		err = errors.New("No such metric: " + strconv.FormatInt(id, 10))
	}

	return
}

func (self *Metrics) Update(id int64, t, detector, md5 string) (err error) {

	result, err := db.Exec(SQL_UPDATE_METRIC, t, detector, md5, id)

	if err != nil {
		return
	}

	c, err := result.RowsAffected()

	if err != nil {
		return
	}

	if c < 1 {
		err = errors.New("No such metric: " + strconv.FormatInt(id, 10))
	}

	return
}

func (self *Metrics) Del(id int64) (err error) {

	result, err := db.Exec(SQL_DELETE_METRIC, id)

	if err != nil {
		return
	}

	c, err := result.RowsAffected()

	if err != nil {
		return
	}

	if c < 1 {
		err = errors.New("No such metric: " + strconv.FormatInt(id, 10))
	}

	return
}

func (self *Metrics) Get(id int64) (metric Metric, err error) {

	err = db.QueryRow(SQL_SELECT_METRIC_BY_ID, id).Scan(&metric.ID, &metric.Name,
		&metric.Type, &metric.Detector, &metric.MD5)

	return
}
