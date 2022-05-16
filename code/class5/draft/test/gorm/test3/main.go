package test3

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Order struct {
}

type Profile struct {
}

type Company struct {
	Alive bool
}

type User struct {
	Orders  []Order
	Profile Profile
}

func main() {
	db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))
	var users []User
	var user User
	// 查询用户的时候并找出其订单，个人信息（1 + 1条SQL）
	db.Preload("Orders").Preload("Profile").Find(&users)
	// SELECT * FROM users
	// SELECT * FROM orders WHERE user_id IN (1,2,3,4); // 一对多
	// SELECT * FROM profiles WHERE user_id IN (1,2,3,4); // 一对一

	// 使用 Join SQL 加载 （单条JOIN SQL）
	db.Joins("Company").Joins("Manager").Find(&user, 1)
	db.Joins("Company", db.Where(&Company{Alive: true})).Find(&users)

	// 预加载关联全部 （只加载一级关联）
	db.Preload(clause.Associations).Find(&users)

	// 多级预加载
	db.Preload("Orders.OrderItems.Product").Find(&users)
	// 多级预加载 + 预加载全部一级关联
	db.Preload("Orders.OrderItems.Product").Preload(clause.Associations).Find(&users)

	// 查询用户的时候找出其未取消的订单
	db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
	db.Preload("Orders", "state =  ?", "paid").Preload("Orders.OrderItems").Find(&users)

	db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
		return db.Order("orders.amount DESC")
	}).Find(&users)
}
