package login

import (
	"Bilibili-Project/global/yaml"
	"encoding/json"
	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type LOGINInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		IsLogin       bool   `json:"isLogin"`
		EmailVerified int    `json:"email_verified"`
		Face          string `json:"face"`
		LevelInfo     struct {
			CurrentLevel int `json:"current_level"`
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
			NextExp      int `json:"next_exp"`
		} `json:"level_info"`
		Mid            int     `json:"mid"`
		MobileVerified int     `json:"mobile_verified"`
		Money          float64 `json:"money"`
		Moral          int     `json:"moral"`
		Official       struct {
			Role  int    `json:"role"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
			Type  int    `json:"type"`
		} `json:"official"`
		OfficialVerify struct {
			Type int    `json:"type"`
			Desc string `json:"desc"`
		} `json:"officialVerify"`
		Pendant struct {
			Pid          int    `json:"pid"`
			Name         string `json:"name"`
			Image        string `json:"image"`
			Expire       int    `json:"expire"`
			ImageEnhance string `json:"image_enhance"`
		} `json:"pendant"`
		Scores       int    `json:"scores"`
		Uname        string `json:"uname"`
		VipDueDate   int64  `json:"vipDueDate"`
		VipStatus    int    `json:"vipStatus"`
		VipType      int    `json:"vipType"`
		VipPayType   int    `json:"vip_pay_type"`
		VipThemeType int    `json:"vip_theme_type"`
		VipLabel     struct {
			Path       string `json:"path"`
			Text       string `json:"text"`
			LabelTheme string `json:"label_theme"`
		} `json:"vip_label"`
		VipAvatarSubscript int    `json:"vip_avatar_subscript"`
		VipNicknameColor   string `json:"vip_nickname_color"`
		Wallet             struct {
			Mid           int `json:"mid"`
			BcoinBalance  int `json:"bcoin_balance"`
			CouponBalance int `json:"coupon_balance"`
			CouponDueTime int `json:"coupon_due_time"`
		} `json:"wallet"`
		HasShop        bool   `json:"has_shop"`
		ShopURL        string `json:"shop_url"`
		AllowanceCount int    `json:"allowance_count"`
		AnswerStatus   int    `json:"answer_status"`
	} `json:"data"`
}

var (
	data yaml.Data
	URL  string = "http://api.bilibili.com/nav"
	path        = "Data/Cache/LoginInfo.yaml"

	LI LOGINInfo = LOGINInfo{}
)

func (LI *LOGINInfo) getLoginInfo() {
	data = yaml.Getdata()
	cookie := http.Cookie{
		Name:  "SESSDATA",
		Value: data.SESSDATA,
	}
	RO := grequests.RequestOptions{
		Cookies: []*http.Cookie{
			&cookie,
		},
	}
	resp, err := grequests.Get(URL, &RO)
	if err != nil {
		log.Fatalln("getLoginInfo Error: ", err)
	}
	err = json.Unmarshal([]byte(resp.String()), &LI)
	if err != nil {
		log.Fatalln("getLoginInfo Error:", err)
	}

	// 将请求到的json 写入 LoginInfo struct 中
	err = json.Unmarshal([]byte(resp.String()), &LI)
	if err != nil {
		log.Fatalln("json写入LoginInfo Error: ", err)
	}

	yaml.WriteYaml(&LI, path)
	log.Info("获取LoginInfo成功")
	time.Sleep((time.Second) * 2)
}

func LoginInfo() {
	if !yaml.PathExists(path) {
		LI.getLoginInfo()
	}
}
