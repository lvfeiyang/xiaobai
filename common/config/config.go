package config

import (
	"encoding/json"
	"github.com/lvfeiyang/xiaobai/common/flog"
	"io/ioutil"
	"runtime"
)

type config struct {
	RedisUrl string
	MongoUrl string
	HtmlPath string
	QiniuAK  string
	QiniuSK  string
	WxAppid  string
	WxSecret string
}

var ConfigVal = &config{}

func Init() {
	var filePath string
	if "linux" == runtime.GOOS {
		filePath = "/data/leon-wp/xiaobai/config"
	} else {
		filePath = "C:\\Users\\Administrator\\config" //lxm19
	}
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		flog.LogFile.Fatal(err)
	}
	err = json.Unmarshal(conf, ConfigVal)
	if err != nil {
		flog.LogFile.Fatal(err)
	}
}
