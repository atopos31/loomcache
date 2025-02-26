package main

import (
	"net/http"
	"sync"
	"testing"
)

func TestSingleFlight(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := http.Get("http://localhost:8888/api?key=Tom")
			if err != nil {
				t.Error(err)
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				t.Error("server returned: ", res.Status)
			}
			t.Log(res.StatusCode)
		}()
	}
	wg.Wait()
}
