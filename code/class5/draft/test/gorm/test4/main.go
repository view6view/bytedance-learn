package test4

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Order struct {
}

type Account struct {
}

type CreditCard struct {
	BillingAddress interface{}
}

type User struct {
	ID          uint
	Name        string
	Account     Account      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreditCards []CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Orders      []Order      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func main() {
	db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))
	var user User
	// 如果需要使用GORM Migrate 数据库迁移数据库外键才行
	db.AutoMigrate(&User{})
	// 入宫未启用软删除，在删除 User 时会自动删除其依赖
	db.Delete(&User{})

	// 方法2：使用 Select 实现级联删除，不依赖数据库的约束以及软删除
	// 删除 user 时候，也删除 user 的 Orders、CreditCards 记录
	db.Select("Account").Delete(&user)

	// 删除 user时候，也删除 user 的Orders、CreditCards 记录
	db.Select("Orders", "CreditCards").Delete(&user)

	// 删除 user 时候， 也删除 user 的 Orders、CreditCards 记录，也删除订单的BillingAddress
	db.Select("Orders", "Orders.BillingAddress", "CreditCards").Delete(&user)

	// 删除 user 时候，也删除用户及其依赖的所有 has one/many、many2many 记录
	db.Select(clause.Associations).Delete(&user)
}
