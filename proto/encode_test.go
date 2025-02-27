package proto

import (
	"encoding/json"
	"testing"

	"google.golang.org/protobuf/proto"
)

var req = &Request{
	Group: "test",
	Key:   "test",
}

// ➜  proto git:(master) ✗ go test -bench=. -benchmem -benchtime=3s
// goos: linux
// goarch: amd64
// pkg: github.com/atopos31/loomcache/proto
// cpu: Intel(R) Xeon(R) Platinum
// BenchmarkJSONMarshal-2          14334495               263.0 ns/op            32 B/op          1 allocs/op
// BenchmarkJSONUnmarshal-2         3683509               974.5 ns/op           304 B/op          7 allocs/op
// BenchmarkProtoMarshal-2         25882170               151.5 ns/op            16 B/op          1 allocs/op
// BenchmarkProtoUnmarshal-2       14115490               253.0 ns/op            88 B/op          3 allocs/op
// PASS
// ok      github.com/atopos31/loomcache/proto     18.620s

func BenchmarkJSONMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// JSON 反序列化基准测试
func BenchmarkJSONUnmarshal(b *testing.B) {
	data, _ := json.Marshal(req)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var p Request
		err := json.Unmarshal(data, &p)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Protobuf 序列化基准测试
func BenchmarkProtoMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Protobuf 反序列化基准测试
func BenchmarkProtoUnmarshal(b *testing.B) {
	data, _ := proto.Marshal(req)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p := &Request{}
		err := proto.Unmarshal(data, p)
		if err != nil {
			b.Fatal(err)
		}
	}
}
