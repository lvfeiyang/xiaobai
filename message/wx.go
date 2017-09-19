package message

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/lvfeiyang/proxy/common/config"
	"github.com/lvfeiyang/proxy/common/session"
	"github.com/lvfeiyang/proxy/common/wx"
	"math/rand"
	"strconv"
	"time"
)

type WxConfigReq struct {
	Url string
}
type WxConfigRsp struct {
	AppId     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

func (req *WxConfigReq) GetName() (string, string) {
	return "wx-config-req", "wx-config-rsp"
}
func (req *WxConfigReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *WxConfigRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *WxConfigReq) Handle(sess *session.Session) ([]byte, error) {
	rsp := &WxConfigRsp{config.ConfigVal.WxAppid, time.Now().Unix(), generateNoncestr(16), ""}
	sigstr := "jsapi_ticket=" + wx.JsTicket() + "&noncestr=" + rsp.NonceStr +
		"&timestamp=" + strconv.FormatInt(rsp.Timestamp, 10) + "&url=" + req.Url
	rsp.Signature = fmt.Sprintf("%x", sha1.Sum([]byte(sigstr)))
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
func generateNoncestr(len uint) string {
	const cr = "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	nonce := make([]byte, len)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range nonce {
		nonce[i] = cr[r.Intn(62)] //62 is len of cr
	}
	return string(nonce)
}
