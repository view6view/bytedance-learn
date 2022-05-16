package test2

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))

	// 注册新的 Callback
	db.Callback().Create().Register("MyPlugin", func(db *gorm.DB) {})

	// 删除 Callback
	db.Callback().Create().Remove("gorm:begin_transaction")

	// 替换 Callback
	db.Callback().Create().Replace("gorm:before_create", func(db *gorm.DB) {})

	// 查询注册的 Callback
	db.Callback().Create().Get("gorm:create")

	// 指定 Callback 顺序
	db.Callback().Create().
		Before("gorm:create").
		After("MyPlugin").
		Register("MyPlugin2", func(db *gorm.DB) {})

	// 注册到所有服务之前
	db.Callback().Create().Before("*").Register("MyPlugin:newCallBack", func(db *gorm.DB) {})

	// 注册时检查条件
	enableTransaction := func(db *gorm.DB) bool { return !db.SkipDefaultTransaction }
	db.Callback().Create().Match(enableTransaction).Register("gorm:begin_transaction", func(db *gorm.DB) {})
}
