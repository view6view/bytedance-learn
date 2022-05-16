package test1

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Account struct {
}

type Company struct {
}

type Language struct {
	Name string
}

var languageZH Language
var languageEN Language
var languageDE Language

type User struct {
	gorm.Model
	Name      string
	Account   Account
	Pets      []*Pet
	Toys      []*Toy `gorm:"polymorphic:Owner;"`
	CompanyID *int
	Company   Company
	ManagerID *uint
	Manager   *User
	Team      []User     `gorm:"foreignkey:ManagerID;"`
	Languages []Language `gorm:"many2many:UserSpeak;"`
	Friends   []*User    `gorm:"many2many:user_friends;"`
}

type Pet struct {
	gorm.Model
	UserID *uint
	Toy    Toy `gorm:"polymorphic:Owner;"`
}

type Toy struct {
	ID        uint
	Name      string
	OwnerId   string
	OwnerType string
	createAt  time.Time
}

func main() {
	db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))
	var user User
	var users []User
	db.Where("role <> ?", "manager").Where("age > ?", 35).
		Limit(100).Order("age desc").Find(&user)

	// 不同数据库甚至不同版本的数据库支持的 SQL 不同
	// SELECT * FROM `users` LOCK IN SHARE MODE // MySQL < 8,MariaDB
	// SELECT * FROM `users` FOR SHARE OF `users` // MySQL 8
	db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).Find(&users)
}
