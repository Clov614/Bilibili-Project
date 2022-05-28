package getVideoDetails

import (
	"Bilibili-Project/global/yaml"
	"encoding/json"
	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type View struct {
	getVideoinfoUrl string
}

var (
	data yaml.Data = yaml.Getdata()
	view View      = View{
		getVideoinfoUrl: "https://api.bilibili.com/x/web-interface/view",
	}
)

func (VI *VideoInfo) GetVideoInfo() {
	cookie := http.Cookie{
		Name:  "SESSDATA",
		Value: data.SESSDATA,
	}
	RO := grequests.RequestOptions{
		Cookies: []*http.Cookie{
			&cookie,
		},
	}
	resp, err := grequests.Get(view.getVideoinfoUrl, &RO)
	if err != nil {
		log.Errorf("GetVideoInfo Error: %v", err)
	}
	err = json.Unmarshal([]byte(resp.String()), VI)
	if err != nil {
		log.Errorf("GetVideoInfo Error: %v", err)
	}
}
