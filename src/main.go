package main

import (
	"github.com/go-ini/ini"
	"tools"
)

func main() {
	cnf, _ := ini.Load("conf.ini")
	var docType, _ = cnf.Section("").Key("flag").Int()
	tools.InitCnf()

	switch docType {
	case 1:
		tools.GenerateJava()
	case 2:
		tools.GenerateDoc()
	}
}
