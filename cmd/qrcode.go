package cmd

import (
	"Bilibili-Project/cmd/qrcodeTerminal"
	"github.com/skip2/go-qrcode"
)

func PrintQRcode(content string) {
	obj := qrcodeTerminal.New()
	obj.Get(content).Print()
	qrcode.WriteFile(content, qrcode.Medium, 256, "QRcode.png")
}
