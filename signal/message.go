package signal

type Message struct {
	Room string
	From string
	To   string
	For  string
	Data interface{}
}

const (
	FOR_CREATE       = "create"
	FOR_JOIN         = "join"
	FOR_JOIN_OLDER   = "join_older"
	FOR_JOIN_NEWER   = "join_newer"
	FOR_LEAVE        = "leave"
	FOR_ICECANDIDATE = "icecandidate"
	FOR_OFFER        = "offer"
	FOR_ANSWER       = "answer"
	FOR_NOTICE       = "notice"
)

func (p *Peer) CreateRoom(m Message) {
	if p.HasRoom() {
		m.For = FOR_NOTICE
		m.Data = "already_in_a_room"
		p.Send(m)
		return
	}
	if m.Room == "" {
		m.For = FOR_NOTICE
		m.Data = "miss_room_id"
		p.Send(m)
		return
	}
	var r *Room
	var err error
	r, err = NewRoom(m.Room)
	if err != nil {
		m.For = FOR_NOTICE
		m.Data = err.Error()
		p.Send(m)
		return
	}
	err = r.Add(p)
	if err != nil {
		m.For = FOR_NOTICE
		m.Data = err.Error()
		p.Send(m)
		return
	}
	p.InRoom(r)
	p.Send(m)
	return
}

func (p *Peer) JoinRoom(m Message) {
	if m.Room == "" {
		m.For = FOR_NOTICE
		m.Data = "miss_room_id"
		p.Send(m)
		return
	}
	if p.HasRoom() {
		m.For = FOR_NOTICE
		m.Data = "already_in_a_room"
		p.Send(m)
		return
	}
	var r *Room
	var ok bool
	roomsMtx.RLock()
	r, ok = Rooms[m.Room]
	roomsMtx.RUnlock()
	if !ok {
		m.For = FOR_NOTICE
		m.Data = "room_not_exists"
		p.Send(m)
		return
	}
	if err := r.Add(p); err != nil {
		m.For = FOR_NOTICE
		m.Data = err.Error()
		p.Send(m)
		return
	}
	p.InRoom(r)
	p.Send(m)
	r.peersMtx.Lock()
	for _, peer := range r.Peers {
		m.For = FOR_JOIN_NEWER
		m.Data = p.ID
		if peer != p {
			peer.Send(m)
		}
		m.For = FOR_JOIN_OLDER
		m.Data = peer.ID
		if peer != p {
			p.Send(m)
		}
	}
	r.peersMtx.Unlock()
}

func (p *Peer) LeaveRoom(m Message) {
	if m.Room == "" {
		m.For = FOR_NOTICE
		m.Data = "miss_room_id"
		p.Send(m)
		return
	}
	var r *Room
	var ok bool
	roomsMtx.RLock()
	r, ok = Rooms[m.Room]
	roomsMtx.RUnlock()
	if !ok {
		m.For = FOR_NOTICE
		m.Data = "room_not_exists"
		p.Send(m)
		return
	}
	if err := r.Remove(p); err != nil {
		m.For = FOR_NOTICE
		m.Data = err.Error()
		p.Send(m)
		return
	}
	r.DestroyOrNot()
	p.OutRoom()
	p.Send(m)
}

func (p *Peer) Forwarded(m Message) {
	if m.Room == "" {
		m.For = FOR_NOTICE
		m.Data = "miss_room_id"
		p.Send(m)
		return
	}
	if m.To == "" {
		m.For = FOR_NOTICE
		m.Data = "miss_to"
		p.Send(m)
		return
	}
	var r *Room
	var ok bool
	var to *Peer
	roomsMtx.Lock()
	r, ok = Rooms[m.Room]
	roomsMtx.Unlock()
	if !ok {
		m.For = FOR_NOTICE
		m.Data = "room_not_exists"
		p.Send(m)
		return
	}
	peersMtx.RLock()
	to, ok = Peers[m.To]
	peersMtx.RUnlock()
	if !ok {
		m.For = FOR_NOTICE
		m.Data = "to_not_exists"
		p.Send(m)
		return
	}
	if !r.Has(p) {
		m.For = FOR_NOTICE
		m.Data = "you_not_in_this_room"
		p.Send(m)
		return
	}
	if !r.Has(to) {
		m.For = FOR_NOTICE
		m.Data = "to_not_in_this_room"
		p.Send(m)
		return
	}
	m.From = p.ID
	to.Send(m)
}
