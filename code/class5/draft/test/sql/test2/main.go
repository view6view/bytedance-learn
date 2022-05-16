package main

import (
	"database/sql"
	"database/sql/driver"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLDriver struct {
}

func (m MySQLDriver) Open(name string) (driver.Conn, error) {
	//TODO implement me
	panic("implement me")
}

func main() {
	sql.Open("mysql", "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&&parseTime=True&loc=Local")
}

// 注册 Driver
func init() {
	sql.Register("mysql", &MySQLDriver{})
}
