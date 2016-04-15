var room;
var leaveToNext = false;
var me = Date.now().toString();
var myStream;
var p, d;

$('#start').click(function(){
    websocket.send(JSON.stringify({
        For: 'join'
    }));
    initStream();
});

$('#next').click(function(){
    websocket.send(JSON.stringify({
        For: 'leave',
        Room: room
    }));
    leaveToNext = true;
});

$('#stop').click(function(){
    websocket.send(JSON.stringify({
        For: 'leave',
        Room: room
    }));
});

$('#send').click(function(){
    $('#chat').html($('#chat').html()+$('#text').val()+"<br/>");
    d.send($('#text').val());
    $('#text').val('');
});

function initStream() {
    navigator.mediaDevices.getUserMedia({
        //audio: true,
        video: {
            frameRate: { ideal: 10, max: 15 },
            width: {max: 300},
            height: {max: 300}
        }
    }).then(function(s) {
        myStream = s;
    }).catch(function(err) {
        console.log(err);
    });
}

function pc() {
    return new RTCPeerConnection({
        iceServers: [
            {url: "stun:sloth.nixisall.com:3478"},
            {
                url: "turn:sloth.nixisall.com:3478",
                username: "fuck",
                credential: "gfw"
            }
        ]
    });
}

var websocket = new WebSocket("wss://law.txthinking.com:444/signal/" + me);
websocket.onopen = function(e) {
    console.log('websocket opened');
}
websocket.onerror = function(e) {
    console.log(e);
}

websocket.onclose = function(e) {
    $('#you').removeAttr('src');
    $('#me').removeAttr('src');
    p.close();
    p.close();
    d.close();
    yourData.close();
    console.log('websocket closed');
}

websocket.onmessage = function(e) {
    var o = JSON.parse(e.data)
    console.log(o)
    switch (o.For) {
        case 'create':
            room = o.Room;
            if (myStream) {
                $('#me').attr('src', URL.createObjectURL(myStream));
            }
            break;
        case 'join':
            room = o.Room;
            if (myStream) {
                $('#me').attr('src', URL.createObjectURL(myStream));
            }
            break;
        case 'join_older':
            p = pc();
            if (myStream) {
                p.addStream(myStream);
            }
            p.onaddstream = function(e) {
                $('#you').attr('src', URL.createObjectURL(e.stream));
            }
            p.onicecandidate = function(e) {
                if (e.candidate) {
                    websocket.send(JSON.stringify({
                        Room: room,
                        From: me,
                        To: o.Data,
                        For: 'icecandidate',
                        Data: e.candidate
                    }));
                }
            }
            d = p.createDataChannel(o.Data);
            d.onopen = function(e) {
                console.log(o.Data, "dc opened");
            }
            d.onmessage = function(e) {
                $('#chat').html($('#chat').html()+e.data+"<br/>");
            }
            d.onclose = function(e) {
                console.log(o.Data, "dc closed");
            }
            p.createOffer().then(function(osd) {
                p.setLocalDescription(osd, function() {
                    websocket.send(JSON.stringify({
                        Room: room,
                        From: me,
                        To: o.Data,
                        For: 'offer',
                        Data: osd
                    }));
                });
            });
            p.onremovestream = function(e) {
                $('#you').removeAttr('src');
            }
            p.oniceconnectionstatechange = function(e) {
                if (p.iceConnectionState === 'connected') {}
                if (p.iceConnectionState === 'completed') {}
                if (p.iceConnectionState === 'disconnected') {
                    $('#you').removeAttr('src');
                }
                if (p.iceConnectionState === 'closed') {}
                console.log("ics", o.Data, p.iceConnectionState)
            }
            p.onsignalingstatechange = function(e) {
                console.log("ss", o.Data, p.signalingState)
            }
            break;
        case 'join_newer':
            p = pc();
            if (myStream) {
                p.addStream(myStream);
            }
            p.onaddstream = function(e) {
                $('#you').attr('src', URL.createObjectURL(e.stream));
            }
            p.onicecandidate = function(e) {
                if (e.candidate) {
                    websocket.send(JSON.stringify({
                        Room: room,
                        From: me,
                        To: o.Data,
                        For: 'icecandidate',
                        Data: e.candidate
                    }));
                }
            }
            p.onremovestream = function(e) {
                $('#you').removeAttr('src');
            }
            p.oniceconnectionstatechange = function(e) {
                if (p.iceConnectionState === 'connected') {}
                if (p.iceConnectionState === 'completed') {}
                if (p.iceConnectionState === 'closed') {}
                if (p.iceConnectionState === 'disconnected') {
                    $('#you').removeAttr('src');
                }
                console.log("ics", o.Data, p.iceConnectionState)
            }
            p.onsignalingstatechange = function(e) {
                console.log("ss", o.Data, p.signalingState)
            }
            p.ondatachannel = function(e) {
                d = e.channel;
                d.onopen = function(e) {
                    console.log(o.Data, "dc opened");
                }
                d.onmessage = function(e) {
                    $('#chat').html($('#chat').html()+e.data+"<br/>");
                }
                d.onclose = function(e) {
                    console.log(o.Data, "dc closed");
                }
            }
            break;
        case 'leave':
            room = "";
            p.close();
            d.close();
            if (leaveToNext) {
                websocket.send(JSON.stringify({
                    For: 'join'
                }));
            }
            break;
        case 'icecandidate':
            p.addIceCandidate(new RTCIceCandidate(o.Data));
            break;
        case 'offer':
            p.setRemoteDescription(new RTCSessionDescription(o.Data), function() {
                p.createAnswer(function(asd) {
                    p.setLocalDescription(asd, function() {
                        websocket.send(JSON.stringify({
                            Room: room,
                            From: me,
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
            p.setRemoteDescription(new RTCSessionDescription(o.Data));
            break;
        case 'notice':
            alert(o.Data);
            break;
        case 'online':
            $('#chat').html($('#chat').html()+"在线人数:"+o.Data+"<br/>");
        default:
            break;
    }
}
