// Copy stream link to share with viewers
function copyToClipboard(text) {
    if (window.clipboardData && window.clipboardData.setData) {
        clipboardData.setData("Text", text);
        return Swal.fire({
            position: 'center',
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
                position: 'center',
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

// Handles WebRTC connection for peers
function connect(stream) {
    document.getElementById('peers').style.display = 'block'
    document.getElementById('chat').style.display = 'flex'
    document.getElementById('noperm').style.display = 'none'
    pc = new RTCPeerConnection({
        iceServers: [{
            'urls': 'stun:stun.l.google.com:19302',
        },
            {
                'urls': 'turn:relay.metered.ca:80',
                'username': '1b176fb3d756c3300bba247a',
                'credential': 'CD/hGxq9WXgZ/UZu',
            },
        ]
    })
    pc.ontrack = function (event) {
        if (event.track.kind === 'audio') {
            return
        }

        col = document.createElement("div")
        col.className = "column is-6 peer"
        let el = document.createElement(event.track.kind)
        el.srcObject = event.streams[0]
        el.setAttribute("controls", "true")
        el.setAttribute("autoplay", "true")
        el.setAttribute("playsinline", "true")
        col.appendChild(el)
        document.getElementById('noone').style.display = 'none'
        document.getElementById('videos').appendChild(col)

        event.track.onmute = function (event) {
            el.play()
        }

        event.streams[0].onremovetrack = ({track}) => {
            if (el.parentNode) {
                el.parentNode.remove()
            }
            if (document.getElementById('videos').childElementCount <= 3) {
                document.getElementById('noone').style.display = 'grid'
            }
        }
    }

    stream.getTracks().forEach(track => pc.addTrack(track, stream))

    let ws = new WebSocket(RoomWebsocketAddr)
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

    ws.onclose = function (evt) {
        console.log("websocket has closed")
        pc.close();
        pc = null;
        pr = document.getElementById('videos')
        while (pr.childElementCount > 3) {
            pr.lastChild.remove()
        }
        document.getElementById('noone').style.display = 'none'
        document.getElementById('nocon').style.display = 'flex'
        setTimeout(function () {
            connect(stream);
        }, 1000);
    }

    ws.onmessage = function (evt) {
        let msg = JSON.parse(evt.data)
        if (!msg) {
            return console.log('failed to parse msg')
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

    ws.onerror = function (evt) {
        console.log("error: " + evt.data)
    }
}

// Webcam & Microphone permissions + settings
navigator.mediaDevices.getUserMedia({
    video: {
        width: {max: 1280},
        height: {max: 720},
        aspectRatio: 4 / 3,
        frameRate: 30,
    },
    audio: {
        sampleSize: 16,
        channelCount: 2,
        echoCancellation: true
    }
})
    .then(stream => {
        document.getElementById('localVideo').srcObject = stream
        connect(stream)
    }).catch(err => console.log(err))

let screenStream = null;
let localStream = null;
let audioTrack = null;
let pc = null;

// Toggle screen sharing on/off
document.getElementById("share-screen-btn").addEventListener("click", async () => {
    try {
        const localVideo = document.getElementById("localVideo");
        const displayMediaOptions = {
            video: true,
            audio: true,
        };
        if (!screenStream) {
            screenStream = await navigator.mediaDevices.getDisplayMedia(displayMediaOptions);
            const videoTracks = screenStream.getVideoTracks();
            console.log("videoTracks: " + videoTracks);
            await pc.getSenders().find(sender => sender.track.kind === 'video').replaceTrack(videoTracks[0], videoTracks[0].clone());
            localVideo.srcObject = screenStream;
            document.getElementById("share-screen-btn").classList.remove("btn-danger");
            document.getElementById("share-screen-btn").classList.add("btn-primary");

            // Disable audio track from localStream
            if (localStream) {
                audioTrack = localStream.getAudioTracks()[0];
                console.log("audioTrack: " + audioTrack);
                audioTrack.enabled = false;
            }
        } else {
            const localVideoStream = await navigator.mediaDevices.getUserMedia({video: true, audio: true});
            const sender = pc.getSenders().find(sender => sender.track.kind === 'video');
            const localVideoTrack = localVideoStream.getVideoTracks()[0];
            const localAudioTrack = localVideoStream.getAudioTracks()[0];
            localStream = new MediaStream([localVideoTrack, localAudioTrack]);
            await sender.replaceTrack(localVideoTrack);
            localVideo.srcObject = localStream;
            document.getElementById("share-screen-btn").classList.remove("btn-primary");
            document.getElementById("share-screen-btn").classList.add("btn-danger");
            screenStream.getTracks().forEach(track => track.stop());
            screenStream = null;

            audioTrack = localAudioTrack;
        }
    } catch (e) {
        console.error("Error sharing screen: ", e);
    }
})

// Toggle microphone on/off
document.getElementById("mute-audio-btn").addEventListener("click", () => {
    let localStream = document.getElementById("localVideo").srcObject;
    if (localStream) {
        let audioTrack = localStream.getAudioTracks()[0];
        if (audioTrack) {
            let enabled = audioTrack.enabled;
            if (enabled) {
                audioTrack.enabled = false;
                document.getElementById("mute-audio-btn").innerHTML = '<i class="fa-solid fa-microphone-slash"></i>';
                document.getElementById("mute-audio-btn").classList.remove("btn-primary");
                document.getElementById("mute-audio-btn").classList.add("btn-danger");
            } else {
                audioTrack.enabled = true;
                document.getElementById("mute-audio-btn").innerHTML = '<i class="fa-solid fa-microphone"></i>';
                document.getElementById("mute-audio-btn").classList.remove("btn-danger");
                document.getElementById("mute-audio-btn").classList.add("btn-primary");
            }
        }
    }
})

// Toggle camera on/off
document.getElementById("mute-video-btn").addEventListener("click", () => {
    let localStream = document.getElementById("localVideo").srcObject;
    if (localStream) {
        let videoTrack = localStream.getVideoTracks()[0];
        let enabled = videoTrack.enabled;
        if (enabled) {
            videoTrack.enabled = false;
            document.getElementById("mute-video-btn").innerHTML = '<i class="fa fa-video-slash"></i>';
            document.getElementById("mute-video-btn").classList.remove("btn-primary");
            document.getElementById("mute-video-btn").classList.add("btn-danger");
        } else {
            videoTrack.enabled = true;
            document.getElementById("mute-video-btn").innerHTML = '<i class="fa fa-video"></i>';
            document.getElementById("mute-video-btn").classList.remove("btn-danger");
            document.getElementById("mute-video-btn").classList.add("btn-primary");
        }
    }
})