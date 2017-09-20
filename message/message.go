package message

import (
	"github.com/lvfeiyang/proxy/message"
	"net/http"
)

var MhMap map[string]message.MsgHandleIF

func Init() {
	MhMap = map[string]message.MsgHandleIF{
		"qiniu-token-req":  &QiniuTokenReq{},
		"event-info-req":   &EventInfoReq{},
		"event-save-req":   &EventSaveReq{},
		"event-delete-req": &EventDeleteReq{},
		"wx-config-req":    &WxConfigReq{},
	}
	return
}

type LocMessage message.Message

func (msg *LocMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	message.GeneralServeHTTP((*message.Message)(msg), w, r, MhMap)
}

/*func MsgMapHandle(name string) message.MsgHandleIF {

}*/
