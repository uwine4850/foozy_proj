import {observeMessages} from "./observe_messages";
import {runLazyLoadMsg, runLazyLoadNotReadMsg} from "./lazy_load_msg";
import {globalIncrementMsgNotification, NotificationConnData} from "./notification_ws";

export enum MsgType{
    Connect,
    TextMsg,
    ReadMsg,
    Error,
}

export interface Msg {
    Type: number;
    Uid: string;
    ChatId: string;
    Msg: Record<string, string>;
}

export interface ConnectData{
    Socket: WebSocket;
    Uid: string;
    ChatId: string;
}

let connectData: ConnectData = {
    Socket: null,
    Uid: null,
    ChatId: null
}

export function RunWs(notification: NotificationConnData): ConnectData{
    let area = document.getElementById("chat-textarea") as HTMLTextAreaElement;
    const socket = new WebSocket("ws://localhost:8000/chat-ws");
    connectData.Socket = socket;
    socket.addEventListener("open", (event) => {
        console.log("Connect.");
    });

    socket.addEventListener("message", (event) => {
        if (event.data == ""){
            return
        }
        const msg: Msg = JSON.parse(event.data);
        switch (msg.Type){
            case MsgType.Error:
                console.log(msg.Msg.Error)
                break;
            case MsgType.TextMsg:
                const chat_content = document.getElementById("chat-content");
                let classes = ""
                let notReadMy = ""
                if (msg.Uid == connectData.Uid){
                    classes = "chat-content-msg-my-msg";
                    notReadMy = '<div class="chat-msg-not-read-my"></div>'
                    sendNotification(msg, notification);
                } else {
                    classes = "chat-msg-not-read chat-msg-not-read-obs";
                }
                chat_content.innerHTML += ` 
                    <div data-msgid="${msg.Msg.Id}" class="chat-content-msg ${classes}">
                        ${notReadMy}
                        <div class="chat-content-msg-text">
                            ${msg.Msg.Text}
                        </div>
                        <div class="chat-content-msg-date">${msg.Msg.Date}</div>
                    </div>`;
                observeMessages(connectData);
                runLazyLoadMsg(connectData);
                runLazyLoadNotReadMsg(connectData);
               break;
            case MsgType.Connect:
                connectData.Uid = msg.Uid;
                connectData.ChatId = msg.ChatId;
                break;
            case MsgType.ReadMsg:
                const element = document.querySelector(`[data-msgid="${msg.Msg.Id}"]`);
                if (element.classList.contains("chat-content-msg-my-msg") && element.querySelectorAll(".chat-msg-not-read-my").length > 0){
                    element.querySelectorAll(".chat-msg-not-read-my")[0].remove();
                }
                if (element.classList.contains("chat-msg-not-read")){
                    element.classList.remove("chat-msg-not-read");
                }
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
            Uid: connectData.Uid,
            ChatId: connectData.ChatId,
            Msg: {"Text": area.value}
        }
        socket.send(JSON.stringify(m));
        area.value = "";
    }
    return connectData;
}

function sendNotification(msg: Msg, notification: NotificationConnData){
    if (msg.Msg.GlobalIncrement == "0"){
        globalIncrementMsgNotification(msg, [msg.Msg.SendToUsersId], notification)
    }
}
