package db

import "testing"

func TestToMap(t *testing.T) {
	e := &Event{Address: "test"}
	m := ToMap(e)
	if m["address"] != "test" {
		t.Error("map fail")
	}
}
