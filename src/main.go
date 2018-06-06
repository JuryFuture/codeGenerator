package main

import (
	"config"
	"fmt"
	"tools"
)

func main() {
	config.InitCnf()
	var docType, _ = config.Cnf.Section("").Key("flag").Int()

	switch docType {
	case 1:
		tools.GenerateJavaV2()
	case 2:
		tools.GenerateDocV2()
	}

	// 等待控制台
	var a string
	fmt.Scanf("%s", a)
}
