function Room(config) {
    this.me = config.me;
    this.signalServer = config.signalServer;
    this.iceServers = config.iceServers;
    this.id;
    this.signal;
    this.stream;
    this.peers = {};
    this.channels = {};
    this.handles = {};
}

Room.prototype.on = function(evt, handle) {
    this.handles[evt] = handle;
}
Room.prototype.in = function() {
    this.signal = new WebSocket(this.signalServer + encodeURIComponent(this.me));
    this.signal.onopen = this._signal_open.bind(this);
    this.signal.onclose = this._signal_close.bind(this);
    this.signal.onerror = this._signal_error.bind(this);
    this.signal.onmessage = this._signal_message.bind(this);
}

Room.prototype._signal_open = function(e) {
    log('r:e:signal_open');
    if(typeof this.handles["signal_open"] === 'function'){
        this.handles["signal_open"](e);
    }
}
Room.prototype._signal_close = function(e) {
    log('r:e:signal_close');
    this._clean();
    if(typeof this.handles["signal_close"] === 'function'){
        this.handles["signal_close"](e);
    }
}
Room.prototype._signal_error = function(e) {
    log('r:e:signal_error');
    this._clean();
    if(typeof this.handles["signal_error"] === 'function'){
        this.handles["signal_error"](e);
    }
}
Room.prototype._signalSend = function(message) {
    log('r:o:'+JSON.stringify(message));
    this.signal.send(JSON.stringify(message));
}
Room.prototype._clean = function() {
    this.id = undefined;
    this.stream = undefined;
    for(var id in this.peers){
        this.peers[id].c.close();
        delete this.peers[id];
    }
    for(var id in this.channels){
        if(this.channels[id].readyState === 'open'){
            this.channels[id].close();
        }
        delete this.channels[id];
    }
}
Room.prototype._newPeerConnection = function() {
    return new RTCPeerConnection({iceServers: this.iceServers});
}

Room.prototype.setStream = function(stream) {
    this.stream = stream;
}
Room.prototype.create = function(id) {
    this._signalSend({
        For: 'create',
        Room: id
    });
}
Room.prototype.join = function(id) {
    this._signalSend({
        For: 'join',
        Room: id
    });
}
Room.prototype.leave = function() {
    this._signalSend({
        For: 'leave',
        Room: this.id
    });
}
Room.prototype.send = function(data) {
    for(var id in this.channels){
        if(this.channels[id].readyState === 'open'){
            this.channels[id].send(data);
        }
    }
}
Room.prototype.peersCount = function() {
    var i = 0;
    for(var id in this.peers){
        if(this.peers[id].readyState === 'connected' ||
                this.peers[id].iceConnectionState  === 'completed'){
            i++;
        }
    }
    return i;
}
Room.prototype.channelsCount = function() {
    var i = 0;
    for(var id in this.channels){
        if(this.channels[id].readyState === 'open'){
            i++;
        }
    }
    return i;
}

Room.prototype._signal_message = function(e) {
    log('r:i:'+e.data);
    var o = JSON.parse(e.data);
    switch (o.For) {
    case "create":
        this.id = o.Room;
        if(typeof this.handles["message_create"] === 'function'){
            log('r:e:message_create');
            this.handles["message_create"](o);
        }
        break;
    case "join":
        this.id = o.Room;
        if(typeof this.handles["message_join"] === 'function'){
            log('r:e:message_join');
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
            log('r:e:message_leave');
            this.handles["message_leave"](o);
        }
        break;
    case 'icecandidate':
        // remoteDescription should be set (which should be done at the moment of receiving offer)
        if(this.peers[o.From].hasRSDP){
            this.peers[o.From].c.addIceCandidate(new RTCIceCandidate(o.Data));
        }else{
            this.peers[o.From].candidates.push(o.Data);
        }
        break;
    case 'offer':
        var self = this;
        self.peers[o.From].c.setRemoteDescription(new RTCSessionDescription(o.Data), function() {
            self.peers[o.From].hasRSDP = true;
            for(;;){
                var cddt = self.peers[o.From].candidates.shift();
                if(!cddt){
                    break;
                }
                self.peers[o.From].c.addIceCandidate(new RTCIceCandidate(cddt));
            }
            self.peers[o.From].c.createAnswer(function(asd) {
                self.peers[o.From].c.setLocalDescription(asd, function() {
                    self._signalSend({
                        Room: self.id,
                        From: self.me,
                        To: o.From,
                        For: 'answer',
                        Data: asd
                    });
                }, function(e){
                    log('r:e:on got offer, set local dsp error');
                });
            }, function(e) {
                log('r:e:on got offer, create answer error');
            });
        }, function(e){
            log('r:e:on got offer, set remote dsp error');
        });
        break;
    case 'answer':
        var self = this;
        self.peers[o.From].c.setRemoteDescription(new RTCSessionDescription(o.Data), function(){
            self.peers[o.From].hasRSDP = true;
            for(;;){
                var cddt = self.peers[o.From].candidates.shift();
                if(!cddt){
                    break;
                }
                self.peers[o.From].c.addIceCandidate(new RTCIceCandidate(cddt));
            }
        }, function(e){
            log('r:e:on got answer, set remote dsp error');
        });
        break;
    case 'notice':
        if(typeof this.handles["message_notice"] === 'function'){
            log('r:e:message_notice');
            this.handles["message_notice"](o);
        }
        break;
    default:
        break;
    }
}

Room.prototype._join_older = function(o) {
    var self = this;
    var c = self._newPeerConnection();
    self.peers[o.Data] = {c: undefined, hasRSDP: false, candidates: []};
    self.peers[o.Data].c = c;
    if (self.stream) {
        c.addStream(self.stream);
    }
    c.onaddstream = function(e) {
        if(typeof self.handles["stream_add"] === 'function'){
            log('r:e:'+o.Data+':stream_add');
            self.handles["stream_add"](o.Data, e.stream, e);
        }
    }
    c.onremovestream = function(e) {
        if(typeof self.handles["stream_remove"] === 'function'){
            log('r:e:'+o.Data+':stream_remove');
            self.handles["stream_remove"](o.Data, e);
        }
    }
    c.onicecandidate = function(e) {
        if (e.candidate) {
            self._signalSend({
                Room: self.id,
                From: self.me,
                To: o.Data,
                For: 'icecandidate',
                Data: e.candidate
            });
        }
    }
    var dataChan = c.createDataChannel(o.Data);
    dataChan.onopen = function(e) {
        self.channels[o.Data] = dataChan;
        if(typeof self.handles["channel_open"] === 'function'){
            log('r:e:'+o.Data+':channel_open');
            self.handles["channel_open"](o.Data, e);
        }
    }
    dataChan.onmessage = function(e) {
        if(typeof self.handles["channel_message"] === 'function'){
            log('r:e:'+o.Data+':channel_message');
            self.handles["channel_message"](o.Data, e.data, e);
        }
    }
    dataChan.onclose = function(e) {
        if(typeof self.handles["channel_close"] === 'function'){
            log('r:e:'+o.Data+':channel_close');
            self.handles["channel_close"](o.Data, e);
        }
    }
    c.createOffer(function(osd) {
        c.setLocalDescription(osd, function() {
            self._signalSend({
                Room: self.id,
                From: self.me,
                To: o.Data,
                For: 'offer',
                Data: osd
            });
        }, function(e){
            log('r:error:'+o.Data+':set local dsp error on create offer');
        });
    }, function(e){
        log('r:error:'+o.Data+':create offer error');
    });
    c.oniceconnectionstatechange = function(e) {
        log('r:ics:offer:'+o.Data +':'+c.iceConnectionState);
        if (c.iceConnectionState === 'connected') {
            if(typeof self.handles["peer_open"] === 'function'){
                log('r:e:'+o.Data+':peer_open');
                self.handles["peer_open"](o.Data, e);
            }
        }
        if (c.iceConnectionState === 'completed') {
            if(typeof self.handles["peer_open"] === 'function'){
                log('r:e:'+o.Data+':peer_open');
                self.handles["peer_open"](o.Data, e);
            }
        }
        if (c.iceConnectionState === 'disconnected') {
            if(typeof self.handles["peer_close"] === 'function'){
                log('r:e:'+o.Data+':peer_close');
                self.handles["peer_close"](o.Data, e);
            }
        }
        if (c.iceConnectionState === 'closed') {
            if(typeof self.handles["peer_close"] === 'function'){
                log('r:e:'+o.Data+':peer_close');
                self.handles["peer_close"](o.Data, e);
            }
        }
        if (c.iceConnectionState === 'failed') {
            if(typeof self.handles["peer_close"] === 'function'){
                log('r:e:'+o.Data+':peer_close');
                self.handles["peer_close"](o.Data, e);
            }
        }
    }
    c.onsignalingstatechange = function(e) {
        log('r:ss:offer:'+o.Data+':' + c.signalingState);
    }
}

Room.prototype._join_newer = function(o) {
    var self = this;
    var c = self._newPeerConnection();
    self.peers[o.Data] = {c: undefined, hasRSDP: false, candidates: []};
    self.peers[o.Data].c = c;
    if (self.stream) {
        c.addStream(self.stream);
    }
    c.onaddstream = function(e) {
        if(typeof self.handles["stream_add"] === 'function'){
            log('r:e:'+o.Data+':stream_add');
            self.handles["stream_add"](o.Data, e.stream, e);
        }
    }
    c.onremovestream = function(e) {
        if(typeof self.handles["stream_remove"] === 'function'){
            log('r:e:'+o.Data+':stream_remove');
            self.handles["stream_remove"](o.Data, e);
        }
    }
    c.onicecandidate = function(e) {
        if (e.candidate) {
            self._signalSend({
                Room: self.id,
                From: self.me,
                To: o.Data,
                For: 'icecandidate',
                Data: e.candidate
            });
        }
    }
    c.oniceconnectionstatechange = function(e) {
        log('r:ics:anwser:'+o.Data+':' + c.iceConnectionState);
        if (c.iceConnectionState === 'connected') {
            if(typeof self.handles["peer_open"] === 'function'){
                log('r:e:'+o.Data+':peer_open');
                self.handles["peer_open"](o.Data, e);
            }
        }
        if (c.iceConnectionState === 'completed') {
            if(typeof self.handles["peer_open"] === 'function'){
                log('r:e:'+o.Data+':peer_open');
                self.handles["peer_open"](o.Data, e);
            }
        }
        if (c.iceConnectionState === 'disconnected') {
            if(typeof self.handles["peer_close"] === 'function'){
                log('r:e:'+o.Data+':peer_close');
                self.handles["peer_close"](o.Data, e);
            }
        }
        if (c.iceConnectionState === 'closed') {
            if(typeof self.handles["peer_close"] === 'function'){
                log('r:e:'+o.Data+':peer_close');
                self.handles["peer_close"](o.Data, e);
            }
        }
        if (c.iceConnectionState === 'failed') {
            if(typeof self.handles["peer_close"] === 'function'){
                log('r:e:'+o.Data+':peer_close');
                self.handles["peer_close"](o.Data, e);
            }
        }
    }
    c.onsignalingstatechange = function(e) {
        log('r:ss:anwser:'+o.Data+':' + c.signalingState);
    }
    c.ondatachannel = function(e) {
        var dataChan = e.channel;
        dataChan.onopen = function(e) {
            self.channels[o.Data] = dataChan;
            if(typeof self.handles["channel_open"] === 'function'){
                log('r:e:'+o.Data+':channel_open');
                self.handles["channel_open"](o.Data, e);
            }
        }
        dataChan.onmessage = function(e) {
            if(typeof self.handles["channel_message"] === 'function'){
                log('r:e:'+o.Data+':channel_message');
                self.handles["channel_message"](o.Data, e.data, e);
            }
        }
        dataChan.onclose = function(e) {
            if(typeof self.handles["channel_close"] === 'function'){
                log('r:e:'+o.Data+':channel_close');
                self.handles["channel_close"](o.Data, e);
            }
        }
    }
}


