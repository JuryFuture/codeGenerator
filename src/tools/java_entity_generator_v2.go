package tools

import (
	"config"
	"entity"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var DIRNAME = "./class/"

var TEMPLATE_FILE = "template/"

var ANNOTATION_ID = "\n    @Id\n    @GeneratedValue(strategy = GenerationType.AUTO)"

var ANNOTATION_DATE = "\n    @Temporal(TemporalType.TIMESTAMP)"

func GenerateJavaV2() {
	packageName = config.Cnf.Section("class").Key("package").String()
	author = config.Cnf.Section("class").Key("author").String()

	tablePrefix = config.Cnf.Section("prefix").Key("tablePrefix").String()
	columnPrefix = config.Cnf.Section("prefix").Key("columnPrefix").String()

	fmt.Println(tablePrefix, columnPrefix)

	tables := config.Cnf.Section("mysql").Key("tables").String()
	tableNames, tableComments := readTables(strings.Split(tables, ","))

	os.Mkdir(DIRNAME, 777)

	for i := 0; i < len(tableNames); i++ {
		generateClassV2(tableNames[i], tableComments[i])
	}

	fmt.Printf("共%d个表\n", len(tableNames))
}

// 生成类文件
func generateClassV2(table, comment string) {
	className := getClassName(table)
	fields := make([]entity.FieldInfo, 0)
	var methods []entity.MethodInfo
	columnNames, dataTypes, columnComments, extras := readColumns(table, 1)
	for i := 0; i < len(columnNames); i++ {
		filed := getFieldInfo(columnNames[i], dataTypes[i], columnComments[i], extras[i])
		fields = append(fields, filed)

		method := getMethodInfo(columnNames[i])
		methods = append(methods, method)
	}

	classInfo := &entity.ClassInfo{packageName, className, dateTime, year, comment, author, date, table, fields, methods}

	tmpl, _ := template.ParseFiles(TEMPLATE_FILE + "class_v2")

	javaFileName := DIRNAME + className + ".java"
	javaFile, _ := os.Create(javaFileName)
	defer javaFile.Close()

	tmpl.Execute(javaFile, classInfo)
}

// 生成属性
func getFieldInfo(columnName, dataType, columnComment, extra string) (filedInfo entity.FieldInfo) {
	fieldType := typeMap[dataType]
	reg := regexp.MustCompile(".*[i/I]d")
	if dataType == DATE_TYPE_INT && reg.MatchString(columnName) {
		fieldType = "long"
	}
	fieldNameTypeMap[columnName] = fieldType

	fieldName := getFieldName(columnName)

	if extra == "auto_increment" {
		extra = ANNOTATION_ID
	} else {
		extra = ""
	}

	others := ""
	if dataType == DATE_TYPE_DATE || dataType == DATE_TYPE_DATETIME || dataType == DATE_TYPE_TIMESTAMP {
		others = ANNOTATION_DATE
	}

	filedInfo = entity.FieldInfo{fieldType, fieldName, columnName, columnComment, extra, others}

	return
}

// 生成getter/setter
func getMethodInfo(columnName string) (methodInfo entity.MethodInfo) {
	fieldType := fieldNameTypeMap[columnName]

	fieldName := getFieldName(columnName)

	upperFieldName := strings.ToUpper(fieldName[:1]) + fieldName[1:]

	methodInfo = entity.MethodInfo{fieldType, fieldName, upperFieldName}

	return
}
