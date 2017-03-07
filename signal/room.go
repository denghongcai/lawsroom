package signal

import (
	"errors"
	"sync"
)

var ROOM_CAPACITY = 5

type Room struct {
	ID       string
	Peers    map[string]*Peer
	peersMtx *sync.RWMutex
}

var Rooms map[string]*Room
var roomsMtx *sync.RWMutex

func NewRoom(id string) (*Room, error) {
	defer roomsMtx.Unlock()
	roomsMtx.Lock()
	if _, ok := Rooms[id]; ok {
		return nil, errors.New("room_exists")
	}
	r := &Room{
		ID:       id,
		Peers:    make(map[string]*Peer),
		peersMtx: &sync.RWMutex{},
	}
	Rooms[id] = r
	return r, nil
}

func (r *Room) Add(p *Peer) error {
	defer r.peersMtx.Unlock()
	r.peersMtx.Lock()
	if _, ok := r.Peers[p.ID]; ok {
		return errors.New("already_in_this_room")
	}
	if len(r.Peers) >= ROOM_CAPACITY {
		return errors.New("room_full")
	}
	r.Peers[p.ID] = p
	return nil
}

func (r *Room) IsFull() bool {
	//defer r.peersMtx.RUnlock()
	//r.peersMtx.RLock()
	if len(r.Peers) >= ROOM_CAPACITY {
		return true
	}
	return false
}

func (r *Room) Remove(p *Peer) error {
	defer r.peersMtx.Unlock()
	r.peersMtx.Lock()
	if _, ok := r.Peers[p.ID]; !ok {
		return errors.New("you_not_in_this_room")
	}
	delete(r.Peers, p.ID)
	return nil
}

func (r *Room) Has(p *Peer) bool {
	defer r.peersMtx.RUnlock()
	r.peersMtx.RLock()
	var ok bool
	_, ok = r.Peers[p.ID]
	return ok
}

func (r *Room) DestroyOrNot() {
	defer roomsMtx.Unlock()
	roomsMtx.Lock()
	if len(r.Peers) == 0 {
		delete(Rooms, r.ID)
	}
}
