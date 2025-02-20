package consistenthash

import (
	"strconv"
	"testing"
)

func TestConsistenthash(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// 2 4 6 12 14 16 22 24 26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2", // 2<=2
		"11": "2", // 11<=12 12->2
		"23": "4", // 23<=24 24->4
		"27": "2", // 27<=2 2->2
	}

	for test, expected := range testCases {
		if ret := hash.Get(test); ret != expected {
			t.Errorf("Asking for %s, should have yielded %s but got %s", test, expected, ret)
		}
	}

	// 8 18 28
	hash.Add("8")

	// 2 4 6 8 12 14 16 18 22 24 26 28
	testCases["27"] = "8" // 27<=28 28->8

	for test, expected := range testCases {
		if ret := hash.Get(test); ret != expected {
			t.Errorf("Asking for %s, should have yielded %s but got %s", test, expected, ret)
		}
	}
}
