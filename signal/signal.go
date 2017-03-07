package signal

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Signal struct {
	Upgrader websocket.Upgrader
	Hooker   Hooker
}

func New(checkOrigin func(r *http.Request) bool, hook Hooker) *Signal {
	signal := &Signal{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     checkOrigin,
		},
		Hooker: hook,
	}
	return signal
}

func (s *Signal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Hooker != nil {
		if err := s.Hooker.BeforeConnect(r); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}

	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//n.Conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(a string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]
	if id == "" {
		_ = conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(10*time.Second))
		_ = conn.Close()
		return
	}
	p, err := NewPeer(id, conn)
	if err != nil {
		_ = conn.WriteJSON(Message{
			For:  FOR_NOTICE,
			Data: err.Error(),
		})
		_ = conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(10*time.Second))
		_ = conn.Close()
		return
	}
	if s.Hooker != nil {
		s.Hooker.AfterNewPeer(p)
	}

	go p.sendQueue()
	for {
		m := Message{}
		err := conn.ReadJSON(&m)
		if err != nil {
			p.Quit()
			if s.Hooker != nil {
				s.Hooker.AfterPeerQuit(p)
			}
			return
		}

		log.Println("reading", m.Room, m.From, m.To, m.For)
		if s.Hooker != nil {
			if err := s.Hooker.BeforeMessage(p, &m); err != nil {
				m.For = FOR_NOTICE
				m.Data = err.Error()
				p.Send(m)
				continue
			}
		}
		switch m.For {
		case FOR_CREATE:
			p.CreateRoom(m)
		case FOR_JOIN:
			p.JoinRoom(m)
		case FOR_LEAVE:
			p.LeaveRoom(m)
		case FOR_ICECANDIDATE, FOR_OFFER, FOR_ANSWER:
			p.Forwarded(m)
		default:
		}
	}
}
