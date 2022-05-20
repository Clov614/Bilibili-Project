package loginQRcode

import (
	"Bilibili-Project/cmd"
	"Bilibili-Project/global/yaml"
	"encoding/json"
	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"time"
)

type LoginDate struct {
	GetloginUrl  string
	LogininitUrl string
	res          RESULT_JSON
}

type RESULT_JSON struct {
	Code   int  `json:"code"`
	Status bool `json:"status"`
	Ts     int  `json:"ts"`
	Data   DATA `json:"data"`
}

type DATA struct {
	Url      string `json:"url"`
	OauthKey string `json:"oauthKey"`
}

//check comfirm status struct
type FalseJson struct {
	Status  bool   `json:"status"`
	Data    int    `json:"data"`
	Message string `json:"message"`
}

// include Cookie
type TrueJson struct {
	Code   int       `json:"code"`
	Status bool      `json:"status"`
	Ts     int       `json:"ts"`
	Data   TokenData `json:"data"`
}

type TokenData struct {
	Url string `json:"url"`
}

var (
	trueData  TrueJson
	falseData FalseJson
)

//check comfirm status struct

// getLoginUrl
func (LD *LoginDate) loginInit() {
	resp, err := grequests.Get(LD.GetloginUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	res := resp.String()

	err = json.Unmarshal([]byte(res), &LD.res)
	if err != nil {
		log.Fatalln(err)
	}
	cmd.PrintQRcode(LD.res.Data.Url)
}

// work in longinConfirm
func checkComfirmStatus(str_Json string) bool {
	err := json.Unmarshal([]byte(str_Json), &trueData)
	if err != nil {
		_ = json.Unmarshal([]byte(str_Json), &falseData)
		return falseData.Status

	}
	return trueData.Status
}

//check comfirm status struct and return CookieInitURL
func (LD *LoginDate) loginConfirm() {
	RO := grequests.RequestOptions{
		Params: map[string]string{
			"oauthKey": LD.res.Data.OauthKey,
		},
	}

	for {
		resp, err := grequests.Post(LD.LogininitUrl, &RO)
		if err != nil {
			log.Fatalln(err)
		}
		status := checkComfirmStatus(resp.String())
		// log.Info(status)
		if status == false {
			time.Sleep(time.Second)
			continue
		}
		if status {
			break
		}
	}

}

func loginJump() (Cookie string) {
	re := regexp.MustCompile("^https://passport.biligame.com/crossDomain\\?([\\s\\S]*?)&gourl=([\\s\\S]*?)$")
	result := re.FindStringSubmatch(trueData.Data.Url)

	Cookie = result[1]
	cookie := http.Cookie{
		Raw: result[1],
	}

	RO := grequests.RequestOptions{
		Cookies: []*http.Cookie{
			&cookie,
		},
	}
	resp, err := grequests.Get(trueData.Data.Url, &RO)
	if err != nil {
		log.Fatalln("loginJump:网络请求错误 ", err)
	}
	if !resp.Ok {
		log.Fatalln("获取Cookie错误，登录测试失败！")
	}
	return
}

//正则提取出SESSDATA
func getSESSDATA() (SESSDATA string) {
	re := regexp.MustCompile("&SESSDATA=([\\s\\S]*?)&bili_jct")
	result := re.FindStringSubmatch(trueData.Data.Url)
	SESSDATA = result[1]
	return
}

func LoginQR() {
	LD := LoginDate{
		GetloginUrl:  "http://passport.bilibili.com/qrcode/getLoginUrl",
		LogininitUrl: "http://passport.bilibili.com/qrcode/getLoginInfo",
		res:          RESULT_JSON{},
	}
	LD.loginInit()
	LD.loginConfirm()
	cookie := loginJump()
	yaml.SaveCookie(cookie, getSESSDATA())
	log.Info("成功生成Data.yaml")
}
