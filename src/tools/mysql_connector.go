package tools

import (
	"database/sql"
	"fmt"
	"strings"
)

// 初始化链接
func connect(userName, password, host, port, schema string) *sql.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", userName, password, host, port, schema)
	fmt.Println(url)
	db, err := sql.Open("mysql", url)
	check(err)

	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(1)
	db.Ping()

	return db
}

// 读取数据库中所有的表
func readTables(tables []string) (tableNames, tableComments []string) {
	db := connect(readConf())

	defer db.Close()

	schema := cnf.Section("mysql").Key("schema").String()
	sql := "SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.`TABLES` WHERE table_schema = '" + schema + "' AND table_type = 'base table'"
	if len(tables) > 0 {
		sql += " AND TABLE_NAME in ('" + strings.Join(tables, "','") + "')"
	}
	rows, _ := db.Query(sql)

	tableNames = make([]string, 0)
	tableComments = make([]string, 0)

	for rows.Next() {
		var tableName, tableComment string
		rows.Scan(&tableName, &tableComment)
		tableNames = append(tableNames, tableName)
		tableComments = append(tableComments, tableComment)
	}

	return
}

// 读取表的所有列
func readColumns(table string, op int) (columnNames, dataTypes, columnComments, extras []string) {
	db := connect(readConf())

	defer db.Close()

	schema := cnf.Section("mysql").Key("schema").String()
	var rows *sql.Rows
	if op == 1 {
		rows, _ = db.Query("select COLUMN_NAME,DATA_TYPE,COLUMN_COMMENT,EXTRA from information_schema.columns where table_schema='" + schema + "' and table_name='" + table + "'")
	} else if op == 2 {
		rows, _ = db.Query("select COLUMN_NAME,COLUMN_TYPE,COLUMN_COMMENT,EXTRA from information_schema.columns where table_schema='" + schema + "' and table_name='" + table + "'")
	}
	columnNames = make([]string, 0)
	dataTypes = make([]string, 0)
	columnComments = make([]string, 0)
	extras = make([]string, 0)

	for rows.Next() {
		var columnName, dataType, columnComment, extra string
		rows.Scan(&columnName, &dataType, &columnComment, &extra)

		columnNames = append(columnNames, columnName)
		dataTypes = append(dataTypes, dataType)
		columnComments = append(columnComments, columnComment)
		extras = append(extras, extra)
	}

	return
}
