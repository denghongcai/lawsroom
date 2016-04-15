package main

import(
    "git.txthinking.com/txthinking/signal"
    "net/http"
    "strconv"
    util "github.com/txthinking/ant"
)

type Dog struct {
}

func (d *Dog) BeforeConnect(r *http.Request) error {
    return nil
}

func (d *Dog) AfterNewPeer(p *signal.Peer) {
    d.Broadcast(signal.Message{
        For: "online",
        Data: len(signal.Peers),
    })
}

func (d *Dog) AfterPeerQuit(p *signal.Peer) {
    d.Broadcast(signal.Message{
        For: "online",
        Data: len(signal.Peers),
    })
}

func (d *Dog) BeforeMessage(in *signal.InMessage) error {
    switch in.Message.For {
    case signal.FOR_JOIN:
        for id, room := range signal.Rooms {
            if !room.IsFull() {
                in.Message.Room = id
                return nil
            }
        }
        in.Message.For = signal.FOR_CREATE
        in.Message.Room = util.SHA1(strconv.Itoa(int(util.RandomNumber())))
    case signal.FOR_LEAVE:
    }
    return nil
}

func (d *Dog) Broadcast(m signal.Message) {
    for _, peer := range signal.Peers {
        go peer.Send(m)
    }
}

