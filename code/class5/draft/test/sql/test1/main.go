package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Name string
}

func main() {
	// import driver 实现
	// 使用driver + DSN 初始化 DB 连接
	db, err := sql.Open("mysql", "user:password@tcp(127.0.01:3306)/hello")

	// 执行一条sql，通过rows取回返回的数据处理完毕，需要释放链接
	rows, err := db.Query("select id, name from users where id = ?", 1)
	if err != nil {
		// XXX
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)

		// 数据、错误处理
		if err != nil {
			// XXX
		}

		users = append(users, user)
	}

	// 错误处理
	if rows.Err() != nil {
		// XXX
	}
}
