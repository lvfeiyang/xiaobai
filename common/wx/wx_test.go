package wx

import (
	"github.com/lvfeiyang/xiaobai/common/config"
	"github.com/lvfeiyang/xiaobai/common/flog"
	"testing"
)

func TestAccessToken(t *testing.T) {
	flog.Init()
	config.Init()
	// token := AccessToken()
	// flog.LogFile.Println("token: ", token)
	ticket := JsTicket()
	flog.LogFile.Println("ticket: ", ticket)
}
