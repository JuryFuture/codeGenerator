package tools

import (
	"fmt"
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

// 包名
var packageName string

// 作者
var author string

// 表名前缀
var tablePrefix string

// 字段名前缀
var columnPrefix string

// 年
var year = time.Now().Format("2006")

// 日期
var date = time.Now().Format("2006/01/02")

// 时间
var dateTime = time.Now().Format("2006/01/02 15:04:05")

func init() {
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

// 表名转换为类名
func getClassName(table string) (class string) {
	if tablePrefix != "" {
		table = strings.Replace(table, tablePrefix, "", 1)
	}

	subNames := strings.Split(table, "_")
	for i := 0; i < len(subNames); i++ {
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
	content = strings.Replace(content, "{{.TableName}}", table, -1)
	content = strings.Replace(content, "{{.TableComment}}", comment, -1)
	content = strings.Replace(content, "{{.PackageName}}", packageName, -1)
	content = strings.Replace(content, "{{.Date}}", date, -1)
	content = strings.Replace(content, "{{.DateTime}}", dateTime, -1)
	content = strings.Replace(content, "{{.Author}}", author, -1)
	content = strings.Replace(content, "{{.Year}}", year, -1)

	className := getClassName(table)
	content = strings.Replace(content, "{{.ClassName}}", className, -1)

	columnNames, dataTypes, columnComments, extras := readColumns(table, 1)

	fields := ""
	methods := ""
	for i := 0; i < len(columnNames); i++ {
		fields += generatorField(columnNames[i], dataTypes[i], columnComments[i], extras[i]) + "\n"
		methods += generatorMethod(columnNames[i]) + "\n"
	}
	content = strings.Replace(content, "{{.Fields}}", fields, -1)
	content = strings.Replace(content, "{{.Methods}}", methods, -1)

	toString := generatorToString(className, columnNames)
	content = strings.Replace(content, "{{.ToString}}", toString, -1)

	dirName := "./class"
	os.Mkdir(dirName, 0777)

	fileName := dirName + "/" + className + ".java"
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
	filed = strings.Replace(filed, "{{.Comment}}", columnComment, -1)
	if extra == "auto_increment" {
		extra = "\n    @Id\n    @GeneratedValue(strategy = GenerationType.AUTO)"
	} else {
		extra = ""
	}
	filed = strings.Replace(filed, "{{.Extra}}", extra, -1)
	filed = strings.Replace(filed, "{{.Column}}", columnName, -1)

	fieldType := typeMap[dataType]
	reg := regexp.MustCompile(".*[i/I]d")
	if dataType == DATE_TYPE_INT && reg.MatchString(columnName) {
		fieldType = "long"
	}
	filed = strings.Replace(filed, "{{.FieldType}}", fieldType, -1)

	others := ""
	if dataType == DATE_TYPE_DATETIME || dataType == DATE_TYPE_TIMESTAMP {
		others = "\n    @Temporal(TemporalType.TIMESTAMP)"
	}
	filed = strings.Replace(filed, "{{.Others}}", others, -1)

	fieldName := getFieldName(columnName)
	filed = strings.Replace(filed, "{{.FieldName}}", fieldName, -1)

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
	method = strings.Replace(method, "{{.FieldType}}", fieldType, -1)

	fieldName := getFieldName(columnName)
	method = strings.Replace(method, "{{.FieldName}}", fieldName, -1)

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
	str = strings.Replace(str, "{{.ClassName}}", className, -1)

	content := ""
	for i := 0; i < len(columnNames); i++ {
		fieldName := getFieldName(columnNames[i])
		content += " + \"," + fieldName + "=\" + " + fieldName
	}
	content = "\"" + content[5:]
	str = strings.Replace(str, "{{.Content}}", content, -1)
	return
}

func GenerateJava() {
	packageName = cnf.Section("class").Key("package").String()
	author = cnf.Section("class").Key("author").String()

	tablePrefix = cnf.Section("prefix").Key("tablePrefix").String()
	columnPrefix = cnf.Section("prefix").Key("columnPrefix").String()

	fmt.Println(tablePrefix, columnPrefix)

	tables := cnf.Section("mysql").Key("tables").String()
	tableNames, tableComments := readTables(strings.Split(tables, ","))

	for i := 0; i < len(tableNames); i++ {
		generateClass(tableNames[i], tableComments[i])
	}

	fmt.Printf("共%d个表\n", len(tableNames))
}
