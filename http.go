package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/atopos31/loomcache/cache"
	"github.com/atopos31/loomcache/consistenthash"
	"github.com/gin-gonic/gin"
)

const DefaultBasePath = "/loomcache/api"

var _ cache.PeerPicker = (*HttpServer)(nil)

type HttpServer struct {
	addr        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]cache.PeerGetter
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{
		addr:        addr,
		basePath:    DefaultBasePath,
		peers:       consistenthash.New(10, nil),
		httpGetters: make(map[string]cache.PeerGetter),
	}
}

// httpool Set
func (h *HttpServer) AddGetters(addr ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers.Add(addr...)
	for _, peer := range addr {
		h.httpGetters[peer] = &httpGetter{
			addr:    peer,
			baseURL: h.basePath,
		}
	}
}

func (p *HttpServer) PickPeer(key string) (cache.PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.addr {
		log.Printf("Pick peer %s for key %s\n", peer, key)
		return p.httpGetters[peer], true
	}
	return nil, false
}

func (h *HttpServer) RunCache(cache *cache.Group) {
	r := gin.Default()
	baser := r.Group(h.basePath)
	{
		baser.GET("/get/:group/:key", func(c *gin.Context) {
			group := c.Param("group")
			key := c.Param("key")
			if group == "" || key == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "group or key is empty",
				})
				return
			}
			if value, err := cache.Get(key); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
			} else {
				c.String(http.StatusOK, value.String())
			}
		})

	}
	r.Run(h.addr)
}

func (h *HttpServer) RunAPI(addr string, cache *cache.Group) {
	r := gin.Default()
	r.GET("/api", func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.JSON(400, gin.H{
				"error": "key is empty",
			})
			return
		}
		if value, err := cache.Get(key); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			c.String(http.StatusOK, value.String())
		}
	})
	r.Run(addr)
}

var _ cache.PeerGetter = (*httpGetter)(nil)

type httpGetter struct {
	addr    string
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"http://%s%v/get/%v/%v",
		h.addr,
		h.baseURL,
		group,
		key,
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	return io.ReadAll(res.Body)
}
