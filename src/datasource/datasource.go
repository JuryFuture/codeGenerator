package datasource

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 初始化链接
func DataBase(userName, password, host, port, schema string) *sql.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", userName, password, host, port, schema)
	fmt.Println(url)
	db, err := sql.Open("mysql", url)
	check(err)

	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(1)
	db.Ping()

	return db
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
