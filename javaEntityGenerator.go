package main

import (
	"database/sql"
	"fmt"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	DATE_TYPE_INT       = "int"
	DATE_TYPE_VARCHAR   = "varchar"
	DATE_TYPE_TINYINT   = "tinyint"
	DATE_TYPE_DATETIME  = "datetime"
	DATE_TYPE_TIMESTAMP = "timestamp"
)

var typeMap = make(map[string]string)

var fieldNameTypeMap = make(map[string]string)

// 配置文件
var cnf *ini.File = nil

// 包名
var packageName string = ""

// 作者
var author string = ""

// 表名前缀
var tablePrefix string = ""

// 字段名前缀
var columnPrefix string = ""

// 年
var year string = time.Now().Format("2006")

// 日期
var date string = time.Now().Format("2006/01/02")

// 时间
var dateTime string = time.Now().Format("2006/01/02 15:04")

func init() {
	cnf, _ = ini.Load("conf.ini")

	packageName = cnf.Section("class").Key("package").String()
	author = cnf.Section("class").Key("author").String()

	tablePrefix = cnf.Section("prefix").Key("tablePrefix").String()
	columnPrefix = cnf.Section("prefix").Key("columnPrefix").String()

	fmt.Println(tablePrefix, columnPrefix)

	typeMap[DATE_TYPE_INT] = "int"
	typeMap[DATE_TYPE_VARCHAR] = "String"
	typeMap[DATE_TYPE_TINYINT] = "int"
	typeMap[DATE_TYPE_DATETIME] = "Date"
	typeMap[DATE_TYPE_TIMESTAMP] = "Date"
	fmt.Println(typeMap)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// 读取数据库配置
func readConf() (username, password, host, port, schema string) {
	username = cnf.Section("mysql").Key("username").String()
	password = cnf.Section("mysql").Key("password").String()
	host = cnf.Section("mysql").Key("host").String()
	port = cnf.Section("mysql").Key("port").String()
	schema = cnf.Section("mysql").Key("schema").String()

	return
}

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
func readTables(table string) (tableNames, tableComments []string) {
	db := connect(readConf())

	defer db.Close()

	schema := cnf.Section("mysql").Key("schema").String()
	sql := "SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.`TABLES` WHERE table_schema = '" + schema + "' AND table_type = 'base table'"
	if table != "" {
		sql += " AND TABLE_NAME = '" + table + "'"
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
func readColumns(table string) (columnNames, dataTypes, columnComments, extras []string) {
	db := connect(readConf())

	defer db.Close()

	schema := cnf.Section("mysql").Key("schema").String()
	rows, _ := db.Query("select COLUMN_NAME,DATA_TYPE,COLUMN_COMMENT,EXTRA from information_schema.columns where table_schema='" + schema + "' and table_name='" + table + "'")

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

// 表名转换为类名
func getClassName(table string) (class string) {
	if tablePrefix != "" {
		table = strings.Replace(table, tablePrefix, "", -1)
	}

	subNames := strings.Split(table, "_")
	class = strings.ToUpper(subNames[0][:1]) + subNames[0][1:]

	for i := 1; i < len(subNames); i++ {
		class += strings.ToUpper(subNames[i][:1]) + subNames[i][1:]
	}

	return
}

// 字段转换为属性名
func getFieldName(column string) (field string) {
	if columnPrefix != "" {
		column = strings.Replace(column, columnPrefix, "", -1)
	}

	subNames := strings.Split(column, "_")
	field = subNames[0]

	for i := 1; i < len(subNames); i++ {
		field += strings.ToUpper(subNames[i][:1]) + subNames[i][1:]
	}

	return
}

// 生成类信息
func generateClass(table, comment string) {
	file, _ := os.Open("template/class")
	defer file.Close()

	data, _ := ioutil.ReadFile(file.Name())

	content := string(data)
	content = strings.Replace(content, "{tableName}", table, -1)
	content = strings.Replace(content, "{tableComment}", comment, -1)
	content = strings.Replace(content, "{packageName}", packageName, -1)
	content = strings.Replace(content, "{date}", date, -1)
	content = strings.Replace(content, "{dateTime}", dateTime, -1)
	content = strings.Replace(content, "{author}", author, -1)
	content = strings.Replace(content, "{year}", year, -1)

	className := getClassName(table)
	content = strings.Replace(content, "{className}", className, -1)

	columnNames, dataTypes, columnComments, extras := readColumns(table)

	fields := ""
	methods := ""
	for i := 0; i < len(columnNames); i++ {
		fields += generatorField(columnNames[i], dataTypes[i], columnComments[i], extras[i]) + "\n"
		methods += generatorMethod(columnNames[i]) + "\n"
	}
	content = strings.Replace(content, "{fields}", fields, -1)
	content = strings.Replace(content, "{methods}", methods, -1)

	toString := generatorToString(className, columnNames)
	content = strings.Replace(content, "{toString}", toString, -1)

	fileName := "class/" + className + ".java"
	class, _ := os.Create(fileName)

	defer class.Close()

	class.WriteString(content)
}

// 生成属性
func generatorField(columnName, dataType, columnComment, extra string) (filed string) {
	file, _ := os.Open("template/field")
	defer file.Close()

	data, _ := ioutil.ReadFile(file.Name())

	filed = string(data)
	filed = strings.Replace(filed, "{comment}", columnComment, -1)
	if extra == "auto_increment" {
		extra = "\n    @Id\n    @GeneratedValue(strategy = GenerationType.AUTO)"
	} else {
		extra = ""
	}
	filed = strings.Replace(filed, "{extra}", extra, -1)
	filed = strings.Replace(filed, "{column}", columnName, -1)

	fieldType := typeMap[dataType]
	reg := regexp.MustCompile(".*[i/I]d")
	if dataType == DATE_TYPE_INT && reg.MatchString(columnName) {
		fieldType = "long"
	}
	filed = strings.Replace(filed, "{fieldType}", fieldType, -1)

	others := ""
	if dataType == DATE_TYPE_DATETIME || dataType == DATE_TYPE_TIMESTAMP {
		others = "\n    @Temporal(TemporalType.TIMESTAMP)"
	}
	filed = strings.Replace(filed, "{others}", others, -1)

	fieldName := getFieldName(columnName)
	filed = strings.Replace(filed, "{fieldName}", fieldName, -1)

	fieldNameTypeMap[columnName] = fieldType

	return
}

// 生成方法
func generatorMethod(columnName string) (method string) {
	file, _ := os.Open("template/method")
	defer file.Close()

	data, _ := ioutil.ReadFile(file.Name())

	method = string(data)
	fieldType := fieldNameTypeMap[columnName]
	method = strings.Replace(method, "{fieldType}", fieldType, -1)

	fieldName := getFieldName(columnName)
	method = strings.Replace(method, "{fieldName}", fieldName, -1)

	upperFieldName := strings.ToUpper(fieldName[:1]) + fieldName[1:]
	method = strings.Replace(method, "{upperFieldName}", upperFieldName, -1)

	return
}

// 生成toString
func generatorToString(className string, columnNames []string) (str string) {
	file, _ := os.Open("template/toString")
	defer file.Close()

	data, _ := ioutil.ReadFile(file.Name())

	str = string(data)
	str = strings.Replace(str, "{className}", className, -1)

	content := ""
	for i := 0; i < len(columnNames); i++ {
		fieldName := getFieldName(columnNames[i])
		content += " + \"," + fieldName + "=\" + " + fieldName
	}
	content = "\"" + content[5:]
	str = strings.Replace(str, "{content}", content, -1)
	return
}

func main() {
	table := cnf.Section("mysql").Key("table").String()
	tableNames, tableComments := readTables(table)

	for i := 0; i < len(tableNames); i++ {
		generateClass(tableNames[i], tableComments[i])
	}

	fmt.Printf("共%d个表\n", len(tableNames))
}
