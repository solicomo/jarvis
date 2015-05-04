package model

var db *sql.DB

func InitDB(driverName, dataSourceName string) (err error) {
	db, err = sql.Open(driverName, dataSourceName)
	return
}
