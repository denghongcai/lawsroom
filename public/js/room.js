function Room(config) {
    this.id;
    this.ws;
    this.myStream;
    this.me = config.me;
    this.signalServer = config.signalServer;
    this.iceServers = config.iceServers;
    this.peerConns = {};
    this.dataChans = {};
    this.handles = {};
}

Room.prototype.on = function(evt, handle) {
    this.handles[evt] = handle;
}
Room.prototype.in = function() {
    this.ws = new WebSocket(this.signalServer + this.me);
    this.ws.onopen = this._wsopen.bind(this);
    this.ws.onclose = this._wsclose.bind(this);
    this.ws.onerror = this._wserror.bind(this);
    this.ws.onmessage = this._wsmessage.bind(this);
}

Room.prototype._wsopen = function(e) {
    if(typeof this.handles["ws_open"] === 'function'){
        this.handles["ws_open"](e);
    }
}
Room.prototype._wsclose = function(e) {
    this._clean();
    if(typeof this.handles["ws_close"] === 'function'){
        this.handles["ws_close"](e);
    }
}
Room.prototype._wserror = function(e) {
    this._clean();
    if(typeof this.handles["ws_error"] === 'function'){
        this.handles["ws_error"](e);
    }
}
Room.prototype._wsSend = function(message) {
    console.log('o', message);
    this.ws.send(JSON.stringify(message));
}
Room.prototype._clean = function() {
    this.id = undefined;
    this.myStream = undefined;
    for(var id in this.peerConns){
        this.peerConns[id].close();
        delete this.peerConns[id];
    }
    for(var id in this.dataChans){
        if(this.dataChans[id].readyState === 'open'){
            this.dataChans[id].close();
        }
        delete this.dataChans[id];
    }
}
Room.prototype._newPeerConn = function() {
    return new RTCPeerConnection({iceServers: this.iceServers});
}

Room.prototype.setMyStream = function(stream) {
    this.myStream = stream;
}
Room.prototype.create = function(id) {
    this._wsSend({
        For: 'create',
        Room: id
    });
}
Room.prototype.join = function(id) {
    this._wsSend({
        For: 'join',
        Room: id
    });
}
Room.prototype.leave = function() {
    this._wsSend({
        For: 'leave',
        Room: this.id
    });
}
Room.prototype.send = function(data) {
    for(var id in this.dataChans){
        if(this.dataChans[id].readyState === 'open'){
            this.dataChans[id].send(data);
        }
    }
}
Room.prototype.peerConnCount = function() {
    var i = 0;
    for(var id in this.peerConns){
        i++;
    }
    return i;
}
Room.prototype.dataChanCount = function() {
    var i = 0;
    for(var id in this.dataChans){
        if(this.dataChans[id].readyState === 'open'){
            i++;
        }
    }
    return i;
}

Room.prototype._wsmessage = function(e) {
    var o = JSON.parse(e.data);
    console.log('i', o);
    switch (o.For) {
    case "create":
        this.id = o.Room;
        if(typeof this.handles["message_create"] === 'function'){
            this.handles["message_create"](o);
        }
        break;
    case "join":
        this.id = o.Room;
        if(typeof this.handles["message_join"] === 'function'){
            this.handles["message_join"](o);
        }
        break;
    case "join_older":
        this._join_older(o);
        break;
    case "join_newer":
        this._join_newer(o);
        break;
    case "leave":
        this._clean();
        if(typeof this.handles["message_leave"] === 'function'){
            this.handles["message_leave"](o);
        }
        break;
    case 'icecandidate':
        this.peerConns[o.From].addIceCandidate(new RTCIceCandidate(o.Data));
        break;
    case 'offer':
        var self = this;
        self.peerConns[o.From].setRemoteDescription(new RTCSessionDescription(o.Data), function() {
            self.peerConns[o.From].createAnswer(function(asd) {
                self.peerConns[o.From].setLocalDescription(asd, function() {
                    self.ws.send(JSON.stringify({
                        Room: self.id,
                        From: self.me,
                        To: o.From,
                        For: 'answer',
                        Data: asd
                    }));
                });
            }, function(e) {
                console.log(e)
            });
        });
        break;
    case 'answer':
        this.peerConns[o.From].setRemoteDescription(new RTCSessionDescription(o.Data));
        break;
    case 'notice':
        if(typeof this.handles["message_notice"] === 'function'){
            this.handles["message_notice"](o);
        }
        break;
    default:
        break;
    }
}

Room.prototype._join_older = function(o) {
    var self = this;
    var peerConn = self._newPeerConn();
    if (self.myStream) {
        peerConn.addStream(self.myStream);
    }
    peerConn.onaddstream = function(e) {
        if(typeof self.handles["remote_stream_add"] === 'function'){
            self.handles["remote_stream_add"](o.Data, e.stream, e);
        }
    }
    peerConn.onremovestream = function(e) {
        if(typeof self.handles["remote_stream_remove"] === 'function'){
            self.handles["remote_stream_remove"](o.Data, e);
        }
    }
    peerConn.onicecandidate = function(e) {
        if (e.candidate) {
            self._wsSend({
                Room: self.id,
                From: self.me,
                To: o.Data,
                For: 'icecandidate',
                Data: e.candidate
            });
        }
    }
    var dataChan = peerConn.createDataChannel(o.Data);
    dataChan.onopen = function(e) {
        self.dataChans[o.Data] = dataChan;
        if(typeof self.handles["data_channel_open"] === 'function'){
            self.handles["data_channel_open"](o.Data, e);
        }
    }
    dataChan.onmessage = function(e) {
        if(typeof self.handles["data_channel_message"] === 'function'){
            self.handles["data_channel_message"](o.Data, e.data, e);
        }
    }
    dataChan.onclose = function(e) {
        delete self.dataChans[o.Data];
        if(typeof self.handles["data_channel_close"] === 'function'){
            self.handles["data_channel_close"](o.Data, e);
        }
    }
    peerConn.createOffer().then(function(osd) {
        peerConn.setLocalDescription(osd, function() {
            self._wsSend({
                Room: self.id,
                From: self.me,
                To: o.Data,
                For: 'offer',
                Data: osd
            });
        });
    });
    peerConn.oniceconnectionstatechange = function(e) {
        console.log("new", peerConn.iceConnectionState);
        if (peerConn.iceConnectionState === 'connected') {
        }
        if (peerConn.iceConnectionState === 'disconnected') {
            delete self.peerConns[o.Data];
            if(typeof self.handles["peer_conn_close"] === 'function'){
                self.handles["peer_conn_close"](o.Data, e);
            }
        }
        if (peerConn.iceConnectionState === 'completed') {
            if(typeof self.handles["peer_conn_open"] === 'function'){
                self.handles["peer_conn_open"](o.Data, e);
            }
        }
        if (peerConn.iceConnectionState === 'closed') {}
    }
    peerConn.onsignalingstatechange = function(e) {
    }
    self.peerConns[o.Data] = peerConn;
}

Room.prototype._join_newer = function(o) {
    var self = this;
    var peerConn = self._newPeerConn();
    if (self.myStream) {
        peerConn.addStream(self.myStream);
    }
    peerConn.onaddstream = function(e) {
        if(typeof self.handles["remote_stream_add"] === 'function'){
            self.handles["remote_stream_add"](o.Data, e.stream, e);
        }
    }
    peerConn.onremovestream = function(e) {
        if(typeof self.handles["remote_stream_remove"] === 'function'){
            self.handles["remote_stream_remove"](o.Data, e);
        }
    }
    peerConn.onicecandidate = function(e) {
        if (e.candidate) {
            self._wsSend({
                Room: self.id,
                From: self.me,
                To: o.Data,
                For: 'icecandidate',
                Data: e.candidate
            });
        }
    }
    peerConn.oniceconnectionstatechange = function(e) {
        console.log('old', peerConn.iceConnectionState)
        if (peerConn.iceConnectionState === 'connected') {
            if(typeof self.handles["peer_conn_open"] === 'function'){
                self.handles["peer_conn_open"](o.Data, e);
            }
        }
        if (peerConn.iceConnectionState === 'disconnected') {
            delete self.peerConns[o.Data];
            if(typeof self.handles["peer_conn_close"] === 'function'){
                self.handles["peer_conn_close"](o.Data, e);
            }
        }
        if (peerConn.iceConnectionState === 'completed') {}
        if (peerConn.iceConnectionState === 'closed') {}
    }
    peerConn.onsignalingstatechange = function(e) {
    }
    peerConn.ondatachannel = function(e) {
        var dataChan = e.channel;
        dataChan.onopen = function(e) {
            self.dataChans[o.Data] = dataChan;
            if(typeof self.handles["data_channel_open"] === 'function'){
                self.handles["data_channel_open"](o.Data, e);
            }
        }
        dataChan.onmessage = function(e) {
            if(typeof self.handles["data_channel_message"] === 'function'){
                self.handles["data_channel_message"](o.Data, e.data, e);
            }
        }
        dataChan.onclose = function(e) {
            delete self.dataChans[o.Data];
            if(typeof self.handles["data_channel_close"] === 'function'){
                self.handles["data_channel_close"](o.Data, e);
            }
        }
    }
    self.peerConns[o.Data] = peerConn;
}


