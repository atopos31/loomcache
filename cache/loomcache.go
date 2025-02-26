package cache

import (
	"errors"
	"log"
	"sync"

	"github.com/atopos31/loomcache/singleflight"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker // get peer and get value
	loader    *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("must set getter!")
	}
	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
		loader: &singleflight.Group{},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is empty")
	}
	// 1 get key from cache
	if v, ok := g.mainCache.get(key); ok {
		log.Printf("cache hit, key: %v", key)
		return v, nil
	}

	viewi, err := g.loader.Do(key, func() (any, error) {
		return g.load(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return ByteView{}, err
}

func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		// 2 get key from peer
		if peer, ok := g.peers.PickPeer(key); ok {
			value, err := peer.Get(g.name, key)
			if err == nil {
				log.Printf("cache hit,form peer:%s key: %v", peer, key)
				return ByteView{b: cloneBytes(value)}, nil
			}
			log.Printf("cache miss,from peer:%s key: %v", peer, key)
		}
	}
	// 3 get key from local
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
