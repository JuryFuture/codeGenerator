package tools

import (
	"config"
	"fmt"
	"testing"
)

func TestGenerateJavaV2(t *testing.T) {
	config.InitCnf()
	fmt.Println(config.Cnf)
	GenerateJavaV2()
}
