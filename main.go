package main

import (
	"Bilibili-Project/global/yaml"
	_ "Bilibili-Project/log"
	"Bilibili-Project/login"
	log "github.com/sirupsen/logrus"
)

var (
	data = yaml.Data{}
)

func init() {
}

func main() {

	data = yaml.Getdata()
	if data.Cookie == "" {
		login.Login()
	} else {
		log.Info("Cache/Data存在,无需登录")
	}
	login.LoginInfo()

}
