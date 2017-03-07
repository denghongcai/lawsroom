package signal

import (
    "net/http"
)

type Hooker interface {
    BeforeConnect(*http.Request) error
    BeforeMessage(*Peer, *Message) error
    AfterNewPeer(p *Peer)
    AfterPeerQuit(p *Peer)
}

