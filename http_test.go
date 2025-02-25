package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/atopos31/loomcache/cache"
)

var httpdb = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestHttp(t *testing.T) {
	loadCounts := make(map[string]int, len(httpdb))
	cache := cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := httpdb[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:8666"
	HttpServer := NewHttpServer(addr)
	go HttpServer.RunCache(cache)

	resp, err := http.Get("http://" + addr + DefaultBasePath + "/get/scores/Tom")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status code error: %d %s", resp.StatusCode, resp.Status)
	}
	// 读取body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "630" {
		t.Fatalf("expect 630, but got %v", string(b))
	}
}
