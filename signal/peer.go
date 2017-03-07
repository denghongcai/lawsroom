package signal

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Peer struct {
	ID      string
	Conn    *websocket.Conn
	Room    *Room
	Message chan Message
	Exit    chan struct{}
}

var Peers map[string]*Peer
var newPeerMtx *sync.Mutex
var peersMtx *sync.RWMutex

func NewPeer(id string, conn *websocket.Conn) (*Peer, error) {
	defer newPeerMtx.Unlock()
	newPeerMtx.Lock()

	if _, ok := Peers[id]; ok {
		return nil, errors.New("peer_id_exists")
	}
	p := &Peer{
		ID:      id,
		Conn:    conn,
		Message: make(chan Message),
		Exit:    make(chan struct{}),
	}
	Peers[id] = p
	return p, nil
}

func (p *Peer) InRoom(r *Room) {
	p.Room = r
}

func (p *Peer) OutRoom() {
	p.Room = nil
}

func (p *Peer) HasRoom() bool {
	return p.Room != nil
}

func (p *Peer) Send(m Message) {
	for {
		select {
		case p.Message <- m:
			return
		case <-p.Exit:
			return
		}
	}
}

func (p *Peer) sendQueue() {
	ticker := time.NewTicker(45 * time.Second)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case m := <-p.Message:
			log.Println("writing", m.Room, m.From, m.To, m.For)
			if err := p.Conn.WriteJSON(m); err != nil {
				return
			}
		case <-ticker.C:
			if err := p.Conn.WriteControl(websocket.PingMessage,
				[]byte{},
				time.Now().Add(10*time.Second)); err != nil {
				return
			}
		}
	}
}

func (p *Peer) Quit() {
	defer peersMtx.Unlock()
	peersMtx.Lock()
	if p.Room != nil {
		_ = p.Room.Remove(p)
		p.Room.DestroyOrNot()
	}
	p.OutRoom()
	_ = p.Conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(10*time.Second))
	_ = p.Conn.Close()
	delete(Peers, p.ID)
	close(p.Exit)
	//close(p.Message)
}
