package tools

import (
	"io/ioutil"
	"os"
)

// 1.生成模板文件
func GenerateTemplateFile(table string) {
	file, _ := os.Open("../../template/class_v2")
	defer file.Close()

	template, _ := ioutil.ReadFile(file.Name())
	content := string(template)

	dirName := "./class"
	os.Mkdir(dirName, 0777)

	className := getClassName(table)
	fileName := dirName + "/" + className + ".java"
	class, _ := os.Create(fileName)

	defer class.Close()

	class.WriteString(content)
}

// 2.生成最终文件
