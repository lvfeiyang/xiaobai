package config

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/flog"
	"io/ioutil"
	"runtime"
)

type config struct {
	RedisUrl    string
	MongoUrl    string
	HtmlPath    string
}

var ConfigVal = &config{}

func Init() {
	var filePath string
	if "linux" == runtime.GOOS {
		filePath = "/root/guild/config"
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
