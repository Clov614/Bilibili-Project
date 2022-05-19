package yaml

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type Data struct {
	Cookie string
}

var (
	data = Data{}
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 生成默认的空白Data.yaml
func init() {
	log.SetFormatter(&nested.Formatter{
		NoColors:        true,
		ShowFullLevel:   true,
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
	})
	log.SetOutput(colorable.NewColorableStdout())

	if !PathExists("./Data/Cache") {
		err := os.MkdirAll("./Data/Cache", 0766)
		if err != nil {
			log.Fatalln("创建Data目录失败:", err)
		}
	}

	if !PathExists("./Data/Cache/Data.yaml") {
		createYaml(&data)
		log.Info("生成默认Data.yaml成功")
	}
}

func SaveCookie(cookie string) {
	data = Data{
		Cookie: cookie,
	}
	createYaml(&data)

}

//将Data写入/Data/Cache/Data.yaml
func createYaml(data *Data) {
	dataStr, err := yaml.Marshal(data)
	if err != nil {
		log.Fatalln("转换Data to yaml error:", err)
	}

	err = ioutil.WriteFile("./Data/Cache/Data.yaml", dataStr, 0644)
	if err != nil {
		log.Fatalln("写入Data.yaml error:", err)
	}

}

func readYaml() {
	file, err := os.ReadFile("./Data/Cache/Data.yaml")
	if err != nil {
		log.Fatalln("读取Data.yaml error:", err)
	}
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Fatalln("Data.yaml to data error: ", err)
	}
}

func Getdata() Data {
	readYaml()
	return data
}
