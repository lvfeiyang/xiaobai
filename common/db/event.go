package db

import "gopkg.in/mgo.v2/bson"

type Event struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	Time    string
	Address string
	Title   string
	Image   string
	Desc    string
}

const eventCName = "event"

func (e *Event) GetById(id bson.ObjectId) error {
	return FindOneById(eventCName, id, e)
}
func (e *Event) Save() error {
	e.Id = bson.NewObjectId()
	return Create(eventCName, e)
}
func (e *Event) UpdateById() error {
	ud := ToMap(e)
	if 0 == len(ud) {
		return nil
	} else {
		return UpdateOne(eventCName, e.Id, bson.M{"$set": ud})
	}
}
func FindAllEvents() ([]Event, error) {
	var es []Event
	err := FindMany(eventCName, bson.M{}, &es, "time")
	return es, err
}
func DelEventById(id bson.ObjectId) error {
	return DeleteOne(eventCName, id)
}
