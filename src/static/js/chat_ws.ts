enum MsgType{
    Connect,
    TextMsg
}

interface TextMsg{
    Text: string;
}

interface Msg {
    Type: number;
    Uid: string;
    ChatId: string;
    Msg: TextMsg;
}

let uid: string;
let chatId: string;

export function RunWs(){
    let area = document.getElementById("chat-textarea") as HTMLTextAreaElement;
    const socket = new WebSocket("ws://localhost:8000/chat-ws");
    socket.addEventListener("open", (event) => {
        console.log("Connect.");
    });

    socket.addEventListener("message", (event) => {
        const msg: Msg = JSON.parse(event.data);
        switch (msg.Type){
            case MsgType.TextMsg:
                if (msg.Uid == uid){
                    document.getElementById("chat-content").innerHTML += ` <div class="chat-content-msg chat-content-msg-my-msg">
                        <div class="chat-content-msg-text">
                            ${msg.Msg.Text}
                        </div>
                        <div class="chat-content-msg-date">1/22/3333</div>
                    </div>`;
                } else {
                    document.getElementById("chat-content").innerHTML += ` <div class="chat-content-msg">
                        <div class="chat-content-msg-text">
                            ${msg.Msg.Text}
                        </div>
                        <div class="chat-content-msg-date">1/22/3333</div>
                    </div>`;
                }
               break;
            case MsgType.Connect:
                uid = msg.Uid;
                chatId = msg.ChatId;
                break;
        }
    });

    socket.addEventListener("close", (event) => {
        console.log("Close.");
    });

    socket.addEventListener("error", (event) => {
        console.error("Error:", event);
    });

    document.getElementById("send").onclick = function () {
        let m: Msg = {
            Type: MsgType.TextMsg,
            Uid: uid,
            ChatId: chatId,
            Msg: {
                Text: area.value,
            }
        }
        socket.send(JSON.stringify(m));
        area.value = "";
    }
}