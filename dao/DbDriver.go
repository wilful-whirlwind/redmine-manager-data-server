package dao

import "database/sql"

type DbDriver struct {
	name         string
	password     string
	host         string
	port         string
	databaseName string
	charset      string
}

func (driver DbDriver) GetDriver() (*sql.DB, error) {
	return sql.Open("mysql", driver.name+":"+driver.password+"@("+driver.host+":"+driver.port+")/"+driver.databaseName+"?charset="+driver.charset)
}

func Driver() (*sql.DB, error) {
	return DbDriver{
		name:         "root",
		password:     "k@SEk1ni",
		host:         "34.146.4.61",
		port:         "3306",
		databaseName: "redmine-manager",
		charset:      "utf8mb4",
	}.GetDriver()
}
