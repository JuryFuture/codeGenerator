package main

import (
	"config"
	"fmt"
	"github.com/go-ini/ini"
	"tools"
)

func main() {
	cnf, _ := ini.Load("conf.ini")
	var docType, _ = cnf.Section("").Key("flag").Int()
	config.InitCnf()

	switch docType {
	case 1:
		tools.GenerateJavaV2()
	case 2:
		tools.GenerateDoc()
	}

	// 等待控制台
	var a string
	fmt.Scanf("%s", a)
}
