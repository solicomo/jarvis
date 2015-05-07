package model

const (
	SQL_CLEAR_METRICS_RECORDS_BEFORE = `DELETE FROM metric_records WHERE 
		julianday(datetime('now','localtime')) - julianday(ctime) > ?;`

	SQL_SELECT_CURRENT_METRICS = `SELECT metric, name, value, ctime FROM current_metrics_view WHERE node = ?;`

	SQL_INSERT_METRIC_RECORD = `INSERT INTO metric_records (node, metric, value, ctime)
		VALUES (?, ?, ?, datetime('now','localtime'));`

	SQL_INSERT_CURRENT_METRIC = `INSERT INTO current_metrics (node, metric, value, ctime)
		VALUES (?, ?, ?, datetime('now','localtime'));`

	SQL_UPDATE_CURRENT_METRIC = `UPDATE current_metrics SET value = ?, ctime = datetime('now','localtime') 
		WHERE node = ? AND metric = ?;`
)

var metrics_records MetricsRecords

func GetMetricsRecords() *MetricsRecords {
	return &metrics_records
}

type MetricRecord struct {
	Node   int64
	Metric int64
	Name   string
	Value  string
	Ctime  string
}

type MetricsRecords struct {
}

func (self *MetricsRecords) ClearBefore(days int) (err error) {
	_, err = db.Exec(SQL_CLEAR_METRICS_RECORDS_BEFORE, days)
	return
}

func (self *MetricsRecords) CurrentFor(node int64) (metrics map[int64]MetricRecord, err error) {

	crows, err := db.Query(SQL_SELECT_CURRENT_METRICS, id)

	if err != nil {
		return
	}

	defer crows.Close()

	metrics = make(map[int64]MetricRecord)

	for crows.Next() {

		metric := MetricRecord{Node: node}

		err = crows.Scan(&metric.Metric, &metric.Name, &metric.Value, &metric.Ctime)

		if err != nil {
			return
		}

		metrics[metric.Metric] = metric
	}

	err = crows.Err()

	return
}

func (self *MetricsRecords) Add(node, metric int64, value string) (err error) {

	_, err = db.Exec(SQL_INSERT_METRIC_RECORD, node, metric, value)

	if err != nil {
		return
	}

	r, err := db.Exec(SQL_UPDATE_CURRENT_METRIC, value, node, metric)

	up := true

	if err != nil {
		up = false
	} else {
		c, e := r.RowsAffected()
		if e != nil || c < 1 {
			up = false
			err = e
		}
	}

	if !up {
		_, err = db.Exec(SQL_INSERT_CURRENT_METRIC, node, metric, value)
	}

	return
}
