package message

import (
	"encoding/json"
	"github.com/lvfeiyang/proxy/common/session"
	"github.com/lvfeiyang/xiaobai/common/db"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

type EventInfoReq struct {
	Id string
}
type EventInfoRsp struct {
	Time    string
	Address string
	Title   string
	Image   string
	Desc    string
}

func (req *EventInfoReq) GetName() (string, string) {
	return "event-info-req", "event-info-rsp"
}
func (req *EventInfoReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *EventInfoRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *EventInfoReq) Handle(sess *session.Session) ([]byte, error) {
	e := db.Event{}
	if bson.IsObjectIdHex(req.Id) {
		(&e).GetById(bson.ObjectIdHex(req.Id))
	}
	rsp := &EventInfoRsp{e.Time, e.Address, e.Title, ImgUrlAddQn(e.Image), e.Desc}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
func ImgUrlAddQn(img string) string {
	domainMapUrl := map[string]string{
		"xiaobai": "http://ov4dqx58l.bkt.clouddn.com",
	}
	re := regexp.MustCompile("(.*?)/")
	imgreg := re.FindStringSubmatch(img)
	if imgreg != nil {
		if url, ok := domainMapUrl[imgreg[1]]; ok {
			return string(re.ReplaceAll([]byte(img), []byte(url+"/"))) + "?imageView2/4/w/300/h/300"
		}
	}
	return img
}

type EventSaveReq struct {
	Id      string
	Time    string
	Address string
	Title   string
	Image   string
	Desc    string
}
type EventSaveRsp struct {
	Result bool
}

func (req *EventSaveReq) GetName() (string, string) {
	return "event-save-req", "event-save-rsp"
}
func (req *EventSaveReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *EventSaveRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *EventSaveReq) Handle(sess *session.Session) ([]byte, error) {
	reqToDb := func(req *EventSaveReq, e *db.Event) {
		if "" != req.Time {
			e.Time = req.Time
		}
		if "" != req.Address {
			e.Address = req.Address
		}
		if "" != req.Title {
			e.Title = req.Title
		}
		if "" != req.Image {
			e.Image = req.Image
		}
		if "" != req.Desc {
			e.Desc = req.Desc
		}
		return
	}

	if bson.IsObjectIdHex(req.Id) {
		e := &db.Event{Id: bson.ObjectIdHex(req.Id)}
		reqToDb(req, e)
		if err := e.UpdateById(); err != nil {
			return nil, err
		}
	} else {
		e := &db.Event{}
		reqToDb(req, e)
		if err := e.Save(); err != nil {
			return nil, err
		}
	}
	rsp := &EventSaveRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}

type EventDeleteReq struct {
	Id string
}
type EventDeleteRsp struct {
	Result bool
}

func (req *EventDeleteReq) GetName() (string, string) {
	return "event-delete-req", "event-delete-rsp"
}
func (req *EventDeleteReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *EventDeleteRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *EventDeleteReq) Handle(sess *session.Session) ([]byte, error) {
	if bson.IsObjectIdHex(req.Id) {
		db.DelEventById(bson.ObjectIdHex(req.Id))
	}
	rsp := &EventDeleteRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
