package test4

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
}

type Session struct {
	PrepareStmt bool
}

func main() {
	// 全局模式，所有 DB 操作都会预编译并缓存（缓存不包含参数部分）
	db, _ := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{PrepareStmt: true})
	var user User
	var users []User
	db.First(&user)

	// 会话模式，后续会话的操作都会预售并缓存
	tx := db.Session(&gorm.Session{PrepareStmt: true})
	tx.Find(&users)
	tx.Model(&user).Update("Age", 18)
	// 全局缓存的语句可能会被会话使用
	tx.First(&user, 2)

	stmtManager, _ := tx.ConnPool.(*gorm.PreparedStmtDB)
	// 关闭当前会话的预编译语句
	stmtManager.Close()
}
