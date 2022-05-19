package login

import (
	"Bilibili-Project/login/loginQRcode"
	log "github.com/sirupsen/logrus"
	"os"
)

func Login() {
	log.Info("请扫码登录:")
	loginQRcode.LoginQR()
	err := os.Remove("QRcode.png")
	if err != nil {
		log.Fatalln("删除QRcode.png出错:", err)
	}
	
}
