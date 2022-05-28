package downloadVideo

import (
	"Bilibili-Project/global/getVideoDetails"
	"Bilibili-Project/global/yaml"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

var (
	data     yaml.Data = yaml.Getdata()
	dash     Dash
	videoUrl string
	audioUrl string
	ID       string
)

func init() {
	if !yaml.PathExists("./videoOutput") {
		err := os.MkdirAll("./VideoOutput", 0644)
		if err != nil {
			log.Errorf("Create VideoOutput error: %v", err)
		}
	}
}

func DownloadVideo(id string, split bool) {
	ID = id
	// 判断是否为avid
	re, _ := regexp.Compile("^\\d*?$")
	if re.MatchString(id) {
		videoUrl, audioUrl = sendRequest(id, false)
	} else {
		videoUrl, audioUrl = sendRequest(id, true)
	}
	saveVideoflow()
	if split {
		splitVideo()
		log.Info("分离音视频成功")
	}

}

func sendRequest(id string, id_type bool) (videoUrl string, audioUrl string) {
	URL := "https://api.bilibili.com/x/player/playurl"
	// 获取视频信息
	VI := &getVideoDetails.VideoInfo{}
	VI.GetVideoInfo(id)
	if VI.Data.Bvid == "" {
		log.Errorf("GetVideoInfo error: Data=nil")
		log.Fatalln()
	}
	cookie := http.Cookie{
		Name:  "SESSDATA",
		Value: data.SESSDATA,
	}
	RO := grequests.RequestOptions{
		Cookies: []*http.Cookie{
			&cookie,
		},
		Params: map[string]string{
			"cid":   strconv.Itoa(VI.Data.Cid),
			"fnval": "80",
			"fourk": "1",
			"fnver": "0",
			"qn":    "0",
		},
	}
	// 将id赋给Params
	if id_type {
		RO.Params["bvid"] = id
	} else {
		RO.Params["avid"] = id
	}
	resp, err := grequests.Get(URL, &RO)
	if err != nil {
		log.Errorf("sendRequest Error: %v", err)
	}
	_ = json.Unmarshal([]byte(resp.String()), &dash)
	//log.Info(dash.Data.Dash.Video[0].Id)
	videoUrl = dash.Data.Dash.Video[0].BaseUrl
	//log.Info(dash.Data.Dash.Audio[0].Id)
	audioUrl = dash.Data.Dash.Audio[0].BaseUrl
	return videoUrl, audioUrl
}

func saveVideoflow() {
	RO := grequests.RequestOptions{
		Headers: map[string]string{
			"referer": "https://www.bilibili.com",
		},
	}
	videoResp, err := grequests.Get(videoUrl, &RO)
	if err != nil {
		log.Errorf("saveVideoflow: %v", err)
	}
	audioResp, err := grequests.Get(audioUrl, &RO)
	if err != nil {
		log.Errorf("saveVideoflow: %v", err)
	}
	videoFlow := videoResp.Bytes()
	audioFlow := audioResp.Bytes()
	saveM4s(videoFlow, audioFlow)
	ffmpegCmd()
}

func saveM4s(videoFlow []byte, audioFlow []byte) {
	err := ioutil.WriteFile("./VideoOutput/Download_video.m4s", videoFlow, 0644)
	if err != nil {
		log.Errorf("saveM4s error: %v", err)
	}
	err = ioutil.WriteFile("./VideoOutput/Download_audio.m4s", audioFlow, 0644)
	if err != nil {
		log.Errorf("saveM4s error: %v", err)
	}
}

func ffmpegCmd() {
	if yaml.PathExists(fmt.Sprintf("./VideoOutput/%s.mp4", ID)) {
		_ = os.Remove(fmt.Sprintf("./VideoOutput/%s.mp4", ID))
	}
	cmd := exec.Command("ffmpeg", "-i", "./VideoOutput/Download_video.m4s", "-i", "./VideoOutput/Download_audio.m4s", "-codec", "copy", fmt.Sprintf("./VideoOutput/%s.mp4", ID))
	//fmt.Println(cmd)
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("cmd error: %v", err)
		return
	}
	_ = os.Remove("./VideoOutput/Download_video.m4s")
	_ = os.Remove("./VideoOutput/Download_audio.m4s")
	log.Info(fmt.Sprintf("保存成功 path: ./VideoOutput/%s.mp4", ID))
}

func splitVideo() {
	cmd := exec.Command("ffmpeg", "-i", fmt.Sprintf("./VideoOutput/%s.mp4", ID), "-vn", fmt.Sprintf("./VideoOutput/%s.flac", ID))
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("split To mp3 error: %v", err)
		return
	}
	cmd = exec.Command("ffmpeg", "-i", fmt.Sprintf("./VideoOutput/%s.mp4", ID), "-vn", fmt.Sprintf("./VideoOutput/%s.mp3", ID))
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("split To flac error: %v", err)
		return
	}
	cmd = exec.Command("ffmpeg", "-i", fmt.Sprintf("./VideoOutput/%s.mp4", ID), "-an", fmt.Sprintf("./VideoOutput/%s_onlyvideo.mp4", ID))
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("split To mp4_onlyvideo error: %v", err)
		return
	}

}
