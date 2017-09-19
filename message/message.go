package message

import(
	"github.com/lvfeiyang/proxy/message"
	"net/http"
)

type LocMessage message.Message
func (msg *LocMessage) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	mhMap := map[string]message.MsgHandleIF{
		"qiniu-token-req": &QiniuTokenReq{},
		"event-info-req": &EventInfoReq{},
		"event-save-req": &EventSaveReq{},
		"event-delete-req": &EventDeleteReq{},
		"wx-config-req": &WxConfigReq{},
	}
	message.GeneralServeHTTP((*message.Message)(msg), w, r, mhMap)
}
/*func MsgMapHandle(name string) message.MsgHandleIF {

}*/
