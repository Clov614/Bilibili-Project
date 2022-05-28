package main

import (
	"Bilibili-Project/global/yaml"
	_ "Bilibili-Project/log"
	"Bilibili-Project/login"
	"Bilibili-Project/module/downloadVideo"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	data = yaml.Data{}
	wg   = sync.WaitGroup{}
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
	log.Info("初始化完成")
LABEL:
	for {
		var order string
		log.Infof("请输入指令 (/help 帮助菜单):")
		_, err := fmt.Scanf("%s \n", &order)
		if err != nil {
			log.Error(err)
		}
		switch order {
		case "/help":
			help()
		case "/dv":
			fallthrough
		case "/DV":
			getVideo()
			break LABEL
		default:
			log.Error("输入错误！！！")
		}
	}

}

func getVideo() {
	var id string
	var YN string
	var split bool
	log.Infof("是否分离音视频:(y/n):")
	_, err := fmt.Scanf("%s \n", &YN)
	if err != nil {
		log.Error(err)
	}
	switch YN {
	case "y":
		fallthrough
	case "Y":
		split = true
	case "n":
		fallthrough
	case "N":
		split = false
	default:
		log.Errorf("输入错误")
		getVideo()
	}
	log.Info("请输入bv或av号(输入Q退出): ")
	fmt.Scanf("%s", &id)
	switch id {
	case "q":
		return
	case "Q":
		return
	default:
		downloadVideo.DownloadVideo(id, split)
	}
}

func help() {
	strHelp := `
********* 帮助菜单 *********
      /dv 下载音视频

`
	log.Info(strHelp)
}
