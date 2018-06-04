package config

import (
	"github.com/go-ini/ini"
)

// 配置文件
var Cnf *ini.File

func InitCnf() {
	Cnf, _ = ini.Load("conf.ini")
}

// 读取数据库配置
func ReadConf() (username, password, host, port, schema string) {
	username = Cnf.Section("mysql").Key("username").String()
	password = Cnf.Section("mysql").Key("password").String()
	host = Cnf.Section("mysql").Key("host").String()
	port = Cnf.Section("mysql").Key("port").String()
	schema = Cnf.Section("mysql").Key("schema").String()

	return
}
