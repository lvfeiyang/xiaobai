package message

import (
	"encoding/hex"
	"encoding/json"
	"github.com/lvfeiyang/xiaobai/common/flog"
	"github.com/lvfeiyang/xiaobai/common/session"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Message struct {
	Name      string
	Data      string
	SessionId uint64
}

func (msg *Message) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg.Name = r.URL.Path[len("/msg/"):] + "-req"
	var err error
	if headSessId := r.Header.Get("SessionId"); "" == headSessId {
		msg.SessionId = 0
	} else {
		msg.SessionId, err = strconv.ParseUint(headSessId, 10, 64)
		if err != nil {
			flog.LogFile.Println(err)
		}
	}
	if 0 == strings.Compare("application/json", r.Header.Get("Content-Type")) {
		defer r.Body.Close()
		buff, err := ioutil.ReadAll(r.Body)
		if err != nil {
			flog.LogFile.Println(err)
		}
		msg.Data = string(buff)

		sendMsg := msg.HandleMsg()
		w.Header().Set("Content-Type", "application/json")
		if 0 == strings.Compare("error-msg", sendMsg.Name) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Write([]byte(sendMsg.Data))
	} else {
		// IDEA: form表单需整合为json
		return
	}
}

func (msg *Message) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}
func (msg *Message) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

type msgHandleIF interface {
	Decode(msgData []byte) error
	Handle(sess *session.Session) ([]byte, error)
	GetName() (string, string)
}

func deCrypto(msgData []byte, sess *session.Session) ([]byte, error) {
	recvEn := make([]byte, hex.DecodedLen(len(msgData)))
	n, err := hex.Decode(recvEn, msgData)
	if err != nil {
		return nil, err
	}
	recv, err := AesDe(recvEn[:n], NewKey(sess.N))
	if err != nil {
		return nil, err
	}
	return recv, nil
}
func handleOneMsg(req msgHandleIF, msgData []byte, sess *session.Session) *Message {
	sendMsg := &Message{Name: "error-msg", Data: UnknowError()}
	reqName, rspName := req.GetName()

	if req.Decode(msgData) != nil {
		sendMsg = &Message{Name: "error-msg", Data: DecodeError(reqName)}
	} else {
		var rspData []byte
		var err interface{}
		// req.SessionId = msgSessId
		rspData, err = req.Handle(sess)
		if err != nil {
			if _, ok := err.(*ErrorMsg); ok {
				sendMsg = &Message{Name: "error-msg", Data: string(rspData)}
			} else {
				flog.LogFile.Println(err)
			}
		} else {
			sendMsg = &Message{Name: rspName, Data: string(rspData)}
		}
	}
	return sendMsg
}

// IDEA: 可改为 rpc 做转发分担 不断线升级
func (msg *Message) HandleMsg() *Message {
	// IDEA: 用接口定义去掉 switch case
	sess := &session.Session{SessId: msg.SessionId}
	if 0 != msg.SessionId {
		if err := sess.Get(msg.SessionId); err != nil {
			errData, _ := NormalError(ErrGetSessionFail)
			return &Message{Name: "error-msg", Data: string(errData)}
		}
	}
	switch msg.Name {
	case "qiniu-token-req":
		return handleOneMsg(&QiniuTokenReq{}, []byte(msg.Data), sess)
	case "event-info-req":
		return handleOneMsg(&EventInfoReq{}, []byte(msg.Data), sess)
	case "event-save-req":
		return handleOneMsg(&EventSaveReq{}, []byte(msg.Data), sess)
	case "event-delete-req":
		return handleOneMsg(&EventDeleteReq{}, []byte(msg.Data), sess)
	case "wx-config-req":
		return handleOneMsg(&WxConfigReq{}, []byte(msg.Data), sess)
	default:
		return &Message{Name: "error-msg", Data: UnknowMsg()}
	}
}
