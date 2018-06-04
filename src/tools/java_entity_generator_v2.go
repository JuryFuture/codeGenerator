package tools

import (
	"config"
	"entity"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var DIRNAME = "./class/"

var TEMPLATE_FILE = "template/"

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

func generateClassV2(table, comment string) {
	className := getClassName(table)
	fields := ""
	methods := ""
	columnNames, dataTypes, columnComments, extras := readColumns(table, 1)
	for i := 0; i < len(columnNames); i++ {
		fields += generatorField(columnNames[i], dataTypes[i], columnComments[i], extras[i]) + "\n"
		methods += generatorMethod(columnNames[i]) + "\n"
	}
	toString := generatorToString(className, columnNames)

	classInfo := &entity.ClassInfo{packageName, className, dateTime, year, comment, author, date, table, fields, methods, toString}

	tmpl, _ := template.ParseFiles(TEMPLATE_FILE + "class")

	javaFileName := DIRNAME + className + ".java"
	javaFile, _ := os.Create(javaFileName)
	defer javaFile.Close()

	tmpl.Execute(javaFile, classInfo)
}
