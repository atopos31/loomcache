package proto

import (
	"encoding/json"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestBytes(t *testing.T) {
	req := &Request{
		Group: "test",
		Key:   "test",
	}
	t.Log("test data:", req)
	pbbytes, err := proto.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	jsonbytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("pb bytes:", len(pbbytes))
	t.Log("json bytes:", len(jsonbytes))
	// pb 降低百分之多少 相对于json?
	t.Log("pb short bytes:", (float64(len(jsonbytes)-len(pbbytes)))/float64(len(jsonbytes)))
}
