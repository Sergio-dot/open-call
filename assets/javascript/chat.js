const msg = document.getElementById("msg");
const log = document.getElementById("log");
const availableCmds = ["help", "howto", "about"]
let msgCount = 0;
let slideOpen = false;

function formatMessage(message) {
    const sections = message.split(/```([\s\S]*?)```/);

    return sections
        .map((section, index) => {
            if (index % 2 === 1) {
                // Format code section
                const code = section.trim();
                const highlightedCode = hljs.highlightAuto(code).value;
                return `<pre><code class="hljs">${highlightedCode}</code></pre>`;
            } else {
                // Format regular message section
                return formatText(section);
            }
        })
        .join("");
}

function formatText(text) {
    // Escape HTML entities to prevent XSS attacks
    const escapedText = escapeHTML(text.trim());

    // Wrap the text in a <p> tag
    return `<p>${escapedText}</p>`;
}

function escapeHTML(html) {
    const escapeChars = {
        "&": "&amp;",
        "<": "&lt;",
        ">": "&gt;",
        '"': "&quot;",
        "'": "&#39;",
    };

    return html.replace(/[&<>"']/g, (char) => escapeChars[char]);
}


function slideToggle() {
    var chat = document.getElementById('chat-content');
    if (slideOpen) {
        chat.style.display = 'none';
        slideOpen = false;
    } else {
        chat.style.display = 'block'
        document.getElementById('chat-alert').style.display = 'none';
        document.getElementById('chat-alert').classList.remove("d-flex");
        document.getElementById('msg').focus();
        msgCount = 0;
        document.getElementById('chat-alert').innerText = msgCount;
        slideOpen = true
    }
}

function appendLog(item) {
    log.appendChild(item);
    log.scrollTop = log.scrollHeight - log.clientHeight;
}

function currentTime() {
    var date = new Date;
    hour = date.getHours();
    minute = date.getMinutes();
    if (hour < 10) {
        hour = "0" + hour
    }
    if (minute < 10) {
        minute = "0" + minute
    }
    return hour + ":" + minute
}

document.getElementById("chat-form").onsubmit = function () {
    if (!chatWs) {
        return false;
    }
    if (!msg.value) {
        return false;
    }
    if (msg.value.startsWith("!")) {
        let messageString = msg.value.split("!");
        messageString = messageString[1];
        if (messageString === availableCmds[2]) {
            Swal.fire(
                'ChatBot says:',
                'OpenCall is an open source video chat platform. It is currently in development' +
                ' but we are doing our best to be in production as soon as possible so... Stay tuned!',
                'question'
            )
        } else if (messageString === availableCmds[1]) {
            Swal.fire(
                'ChatBot says:',
                'To start streaming you just need to create an account from our home page,' +
                ' after that you just need to create a room and share the link to your friends/colleagues',
                'question'
            )
        } else if (messageString === availableCmds[0]) {
            Swal.fire(
                'ChatBot says:',
                'Available commands are: ' +
                availableCmds[0] + ', ' +
                availableCmds[1] + ', ' +
                availableCmds[2],
                'question'
            )
        } else {
            Swal.fire(
                'ChatBot says:',
                'Unknown command. Please use !help to see available commands.',
                'error'
            )
        }
    } else {
        chatWs.send(msg.value);
    }

    msg.value = "";
    return false;
};

function connectChat() {
    chatWs = new WebSocket(ChatWebsocketAddr)

    chatWs.onclose = function (evt) {
        console.log("websocket has closed")
        document.getElementById('chat-button').disabled = true
        setTimeout(function () {
            connectChat();
        }, 1000);
    }

    chatWs.onmessage = function (evt) {
        const messages = evt.data.split('\n');
        if (slideOpen === false) {
            msgCount += 1;
            document.getElementById('chat-alert').style.display = 'flex';
            document.getElementById('chat-alert').classList.add("d-flex");
            document.getElementById('chat-alert').innerText = msgCount;
        }
        for (let i = 0; i < messages.length; i++) {
            const item = document.createElement("div");

            item.innerHTML = formatMessage(messages[i]);
            appendLog(item);
        }
    }

    chatWs.onerror = function (evt) {
        console.log("error: " + evt.data)
    }

    setTimeout(function () {
        if (chatWs.readyState === WebSocket.OPEN) {
            document.getElementById('chat-button').disabled = false
        }
    }, 1000);
}

connectChat();