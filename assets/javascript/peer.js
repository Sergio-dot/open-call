function copyToClipboard(text) {
    if (window.clipboardData && window.clipboardData.setData) {
        clipboardData.setData("Text", text);
        return Swal.fire({
            position: 'top-end',
            text: 'Copied to clipboard',
            showConfirmButton: false,
            timer: 1000,
            width: '150px',
        })
    } else if (document.queryCommandSupported && document.queryCommandSupported("copy")) {
        let textarea = document.createElement("textarea");
        textarea.textContent = text;
        textarea.style.position = "fixed";
        document.body.appendChild(textarea);
        textarea.select();
        try {
            document.execCommand("copy");
            return Swal.fire({
                position: 'top-end',
                text: 'Copied to clipboard',
                showConfirmButton: false,
                timer: 1000,
                width: '150px',
            })
        } catch (ex) {
            console.warn("Copy to clipboard failed:", ex);
            return false;
        } finally {
            document.body.removeChild(textarea);
        }
    }
}

function connect(stream) {
    document.getElementById("no-perm").style.display = 'none'

    // creates a new peer connection
    let pc = new RTCPeerConnection({
        iceServers: [{
            'urls': 'stun:turn.videochat:3478',
        },
            {
                'urls': 'turn:turn.videochat:3478',
                'username': 'user',
                'credential': 'user',
            }
        ]
    })

    // handles track event
    pc.ontrack = function (event) {
        if (event.track.kind === 'audio') {
            return
        }

        // creates new video container in the DOM
        let newCol = document.createElement("div")
        newCol.className = "col"
        let el = document.createElement(event.track.kind)
        el.srcObject = event.streams[0]
        el.setAttribute("controls", "true")
        el.setAttribute("autoplay", "true")
        el.setAttribute("playsinline", "true")
        newCol.appendChild(el)
        document.getElementById('no-one').style.display = 'none'
        document.getElementById('no-con').style.display = 'none'
        document.getElementById('videos').appendChild(newCol)

        // handles mute event
        event.track.onmute = function (event) {
            el.play()
        }

        // handles remove track event
        event.streams[0].onremovetrack = ({track}) => {
            if (el.parentNode) {
                el.parentNode.remove()
            }

            if (document.getElementById('videos').childElementCount <= 3) {
                document.getElementById('no-one').style.display = 'grid'
                document.getElementById('no-streamer').style.display = 'grid'
            }
        }
    }

    stream.getTracks().forEach(track => pc.addTrack(track, stream))

    // creates websocket for new peer
    let ws = new WebSocket(RoomWebsocketAddr)

    // handles ice candidates
    pc.onicecandidate = e => {
        if (!e.candidate) {
            return
        }

        ws.send(JSON.stringify({
            event: 'candidate',
            data: JSON.stringify(e.candidate)
        }))
    }

    ws.addEventListener('error', function (event) {
        console.log('error: ', event)
    })

    // handles websocket closure
    ws.onclose = function (evt) {
        console.log("websocket has been closed")
        pc.close();
        pc = null;
        pr = document.getElementById('videos')
        while (pr.childElementCount > 3) {
            pr.lastChild.remove()
        }
        document.getElementById('no-one').style.display = 'none'
        document.getElementById('no-con').style.display = 'flex'
        setTimeout(function () {
            connect(stream);
        }, 1000);
    }

    // handles websocket messages
    ws.onmessage = function (evt) {
        let msg = JSON.parse(evt.data)
        if (!msg) {
            return console.log("failed to parse message")
        }

        switch (msg.event) {
            case 'offer':
                let offer = JSON.parse(msg.data)
                if (!offer) {
                    return console.log('failed to parse answer')
                }
                pc.setRemoteDescription(offer)
                pc.createAnswer().then(answer => {
                    pc.setLocalDescription(answer)
                    ws.send(JSON.stringify({
                        event: 'answer',
                        data: JSON.stringify(answer)
                    }))
                })
                return

            case 'candidate':
                let candidate = JSON.parse(msg.data)
                if (!candidate) {
                    return console.log('failed to parse candidate')
                }

                pc.addIceCandidate(candidate)
        }
    }

    // handles websocket errors
    ws.onerror = function (evt) {
        console.log('error: ' + evt.data)
    }
}

navigator.mediaDevices.getUserMedia({
    video: {
        width: {
            max: 1280
        },
        height: {
            max: 720
        },
        aspectRatio: 4 / 3,
        frameRate: 30,
    },
    audio: {
        sampleSize: 16,
        channelCount: 2,
        echoCancellation: true
    }
}).then(stream => {
    document.getElementById('localVideo').srcObject = stream
    connect(stream)
}).catch(err => console.log(err))