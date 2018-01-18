package tools

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var TABLE_TEMPLATE = "<p>-----------------------------</p><span>表名称: {{table_name}}</span></br><span>行数: {{column_num}}</span></br><span>注释: {{table_comment}}</span><table cellspacing='0'><thead><tr><td>字段名</th><td>类型</td><td>说明</td></tr></thead>{{items}}</table></br>"
var TABLE_ELEMENT_TEMPLATE = "<tr><td>{{column_name}}<td>{{data_type}}</td><td>{{comment}}</td></th></tr>"

func init() {

}

func generateHtml(str string) {
	dirName := "./html"
	os.Mkdir(dirName, 0777)
	f, err := os.Create(dirName + "/" + "./tables.html")
	check(err)
	defer f.Close()

	f.WriteString(str)
}

func GenerateDoc() {
	tables := cnf.Section("mysql").Key("tables").String()
	tableNames, tableComments := readTables(strings.Split(tables, ","))

	var tableStr string = ""
	for i := 0; i < len(tableNames); i++ {
		tableTmplate := strings.Replace(TABLE_TEMPLATE, "{{table_name}}", tableNames[i], -1)
		tableTmplate = strings.Replace(tableTmplate, "{{table_comment}}", tableComments[i], -1)

		columnNames, columnTypes, columnComments, _ := readColumns(tableNames[i])

		tableTmplate = strings.Replace(tableTmplate, "{{column_num}}", strconv.Itoa(len(columnNames)+1), -1)

		items := ""
		for j := 0; j < len(columnNames); j++ {
			columnTmplate := strings.Replace(TABLE_ELEMENT_TEMPLATE, "{{column_name}}", columnNames[j], -1)
			columnTmplate = strings.Replace(columnTmplate, "{{data_type}}", columnTypes[j], -1)
			columnTmplate = strings.Replace(columnTmplate, "{{comment}}", columnComments[j], -1)

			items += columnTmplate
		}

		tableTmplate = strings.Replace(tableTmplate, "{{items}}", items, -1)

		tableStr += tableTmplate
	}

	fmt.Printf("共%d个表\n", len(tableNames))

	generateHtml(tableStr)
}
