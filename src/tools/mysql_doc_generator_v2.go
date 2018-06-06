package tools

import (
	"config"
	"entity"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var DIRNAME_DOC = "./html/"

func GenerateDocV2() {
	fmt.Println("http://127.0.0.1:0811/index")

	http.HandleFunc("/index", handler)
	http.ListenAndServe(":0811", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	tables := config.Cnf.Section("mysql").Key("tables").String()
	tableList := make([]string, 0)
	if tables != "" {
		tableList = strings.Split(tables, ",")
	}
	tableNames, tableComments := readTables(tableList)

	var docTableInfos []entity.DocTableInfo
	for i := 0; i < len(tableNames); i++ {
		columnNames, columnTypes, columnComments, _ := readColumns(tableNames[i], 2)

		docTableInfo := new(entity.DocTableInfo)
		docTableInfo.TableName = tableNames[i]
		docTableInfo.TableComment = tableComments[i]
		docTableInfo.ColumnNum = len(columnNames) + 1

		docColumnInfos := make([]entity.DocColumnInfo, 0)
		for j := 0; j < len(columnNames); j++ {
			docColumnInfo := entity.DocColumnInfo{columnNames[j], columnTypes[j], columnComments[j]}

			docColumnInfos = append(docColumnInfos, docColumnInfo)
		}
		docTableInfo.Columns = docColumnInfos

		docTableInfos = append(docTableInfos, *docTableInfo)
	}

	os.Mkdir(DIRNAME_DOC, 777)

	tmpl, _ := template.ParseFiles(TEMPLATE_FILE + "doc")

	/*docFileName := DIRNAME_DOC + "tables.html"
	docFile, _ := os.Create(docFileName)
	defer docFile.Close()

	tmpl.Execute(docFile, docTableInfos)*/

	tmpl.Execute(w, docTableInfos)

	fmt.Printf("共%d个表\n", len(tableNames))
}
