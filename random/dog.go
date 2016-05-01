package main

import(
    "git.txthinking.com/txthinking/signal"
    "net/http"
    "errors"
)

type Dog struct {
}

func (d *Dog) BeforeConnect(r *http.Request) error {
    return nil
}

func (d *Dog) AfterNewPeer(p *signal.Peer) {
    d.Broadcast(signal.Message{
        For: "notice",
        Data: map[string]int{"online": len(signal.Peers)},
    })
}

func (d *Dog) AfterPeerQuit(p *signal.Peer) {
    delete(Pairs, p.ID)
    d.Broadcast(signal.Message{
        For: "notice",
        Data: map[string]int{"online": len(signal.Peers)},
    })
}

func (d *Dog) BeforeMessage(in *signal.InMessage) error {
    switch in.Message.For {
    case signal.FOR_JOIN:
        for id, room := range signal.Rooms {
            if room.IsFull() {
                continue
            }

            var lastPair string
            var ok bool
            var you string
            for you, _ = range room.Peers {
                break
            }

            lastPair, ok = Pairs[in.Peer.ID]
            if !ok {
                MakePair(in.Peer.ID, you)
                in.Message.Room = id
                return nil
            }
            if you != lastPair {
                MakePair(in.Peer.ID, you)
                in.Message.Room = id
                return nil
            }
        }
        return errors.New("no_pair_room")
    case signal.FOR_LEAVE:
    }
    return nil
}

func (d *Dog) Broadcast(m signal.Message) {
    for _, peer := range signal.Peers {
        go peer.Send(m)
    }
}

