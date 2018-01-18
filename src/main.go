package main

import (
	"flag"
	"tools"
	"fmt"
)

var docType = 1

func main() {
	flag.IntVar(&docType, "docType", 1, "docType")

	tools.InitCnf()

	switch docType {
	case 1:
		tools.GenerateJava()
	case 2:
		tools.GenerateDoc()
	}
}
