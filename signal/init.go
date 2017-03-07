package signal

import (
	"sync"
)

func init() {
	Peers = make(map[string]*Peer)
	Rooms = make(map[string]*Room)
	newPeerMtx = &sync.Mutex{}
	peersMtx = &sync.RWMutex{}
	roomsMtx = &sync.RWMutex{}
}
