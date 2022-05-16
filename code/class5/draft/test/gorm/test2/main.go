package test2

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	// 保存用户及其关联 （Upsert）
	db.Save(&User{
		Name:      "wuyou",
		Languages: []Language{{Name: "zh-CN"}, {Name: "en-US"}},
	})

	user := User{}
	languages := Language{}
	// 关联模式
	langAssociation := db.Model(&user).Association("Languages")
	// 查询关联
	langAssociation.Find(&languages)
	// 将汉语，英语添加到用户掌握的语言中
	langAssociation.Append([]Language{languageZH, languageEN})
	// 把用户掌握的语言替换为汉语，德语
	langAssociation.Replace([]Language{languageZH, languageDE})
	// 删除用户掌握的两个语言
	langAssociation.Delete(languageZH, languageEN)
	// 删除用户掌握的语言
	langAssociation.Clear()
	// 返回用户所掌握的语言的数量
	langAssociation.Count()

	var user1 User
	var user2 User
	var user3 User
	users := []User{user1, user2, user3}
	// 批量模式 Append, Replace
	langAssociation = db.Model(&users).Association("Languages")

	db.Model(&users).Association("Team").Append(&user1, &user2, &[]User{user1, user2, user3})
}
