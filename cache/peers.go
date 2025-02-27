package cache

import "github.com/atopos31/loomcache/proto"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
// It's used to retrieve a specific key from a peer.
type PeerGetter interface {
	// Get(group string, key string) ([]byte, error)
	Get(req *proto.Request) (*proto.Response, error)
}
