package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func main() {
	connector, _ := mysql.NewConnector(&mysql.Config{
		User:      "gorm",
		Passwd:    "gorm",
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "gorm",
		ParseTime: true,
	})

	sql.OpenDB(connector)
}
