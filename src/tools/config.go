package tools

import (
	"github.com/go-ini/ini"
)

// 配置文件
var cnf *ini.File

func InitCnf() {
	cnf, _ = ini.Load("conf.ini")
}

// 读取数据库配置
func readConf() (username, password, host, port, schema string) {
	username = cnf.Section("mysql").Key("username").String()
	password = cnf.Section("mysql").Key("password").String()
	host = cnf.Section("mysql").Key("host").String()
	port = cnf.Section("mysql").Key("port").String()
	schema = cnf.Section("mysql").Key("schema").String()

	return
}
