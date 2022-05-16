package test1

import (
	"database/sql"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// User 模型定义
type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     time.Time
	MemberNumber sql.NullTime
	ActivateAt   sql.NullTime
	CreateAt     time.Time
	UpdateAt     time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// User1 模型定义
type User1 struct {
	gorm.Model
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullTime
	ActivateAt   sql.NullTime
}

type Product struct {
	Price int
	Code  string
}

func main() {
	db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))

	var users []User
	_ = db.Select("id", "name").Find(&users, 1).Error

	// curd
	// 操作数据库
	db.AutoMigrate(&Product{})
	db.Migrator().CreateTable(&Product{})
	time.Now()

	// 创建
	user := User{Name: "wuyou", Age: 18, Birthday: time.Now()}
	// 通过数据的指针创建
	result := db.Create(&user)

	// 返回主键 last insert id
	print(user.ID)
	// 返回 error
	print(result.Error)
	// 返回影响的行数
	print(result.RowsAffected)

	// 批量创建
	users = []User{{Name: "wuyou"}, {Name: "wuyou"}, {Name: "wuyou"}}
	db.Create(&users)
	db.CreateInBatches(users, 100)
	for _, user := range users {
		print(user.ID)
	}

	// 读取
	var product Product
	// 查询id为1
	db.First(&product, 1)
	// 查询code为L1212的product
	db.First(&product, "code = ?", "L1212")

	result = db.Find(&users, []int{1, 2, 3})
	// 找到的记录数
	print(result.RowsAffected)
	// First, Last, Take 查不到数据
	errors.Is(result.Error, gorm.ErrRecordNotFound)

	// 更新某个字段
	db.Model(&product).Update("Price", 2000)
	db.Model(&product).UpdateColumn("Price", 2000)

	// 更新多个字段
	db.Model(&product).Updates(Product{Price: 2000, Code: "L1212"})
	db.Model(&product).Updates(map[string]interface{}{"Price": 2000, "Code": "L1212"})

	// 批量更新
	db.Model(&Product{}).Where("price < ?", 2000).Updates(map[string]interface{}{"Price": 2000})

	// 删除 product
	db.Delete(&product)
}
