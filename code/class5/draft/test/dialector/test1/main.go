package test1

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=true&loc=Local"
	_, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	_, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	_, _ = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	// import "xxx.io/caches"
	//_, _ = gorm.Open(caches.New(caches.Config{
	//	Fallback: mysql.Open(dsn),
	//	Store:    lru.New(lru.Config{}),
	//}))
}
