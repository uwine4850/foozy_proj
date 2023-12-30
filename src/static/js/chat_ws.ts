import {observeMessages} from "./observe_messages";
import {runLazyLoadMsg, runLazyLoadNotReadMsg} from "./lazy_load_msg";
import {handleError, SendAjaxChatMessage} from "./chat/chat";

export enum MessageType{
    WsConnect,
    WsTextMsg,
    WsReadMsg,
    WsError,
    WsImageNsg
}

export interface IMessage {
    Type: number;
    Uid: string;
    ChatId: string;
    Msg: Record<string, string>;
}

export interface IConnectData {
    Socket: WebSocket;
    Uid: string;
    ChatId: string;
}

type MsgActionFunction = (message: IMessage, ws: WebSocket) => void;
type MsgActions = Record<MessageType, MsgActionFunction>;

const MessageActions: MsgActions = {
    [MessageType.WsConnect]: handleWsConnect,
    [MessageType.WsTextMsg]: handleWsTextMsg,
    [MessageType.WsReadMsg]: handleWsReadMsg,
    [MessageType.WsError]: handleWsError,
    [MessageType.WsImageNsg]: handleWsImageMsg,
}

let connectData: IConnectData = {
    Socket: null,
    Uid: null,
    ChatId: null
}

function handleWsConnect(message: IMessage, ws: WebSocket){
    connectData.Uid = message.Uid;
    connectData.ChatId = message.ChatId;
    SendAjaxChatMessage(connectData);
}

function handleWsTextMsg(message: IMessage, ws: WebSocket){
    const chat_content = document.getElementById("chat-content");
    let classes = ""
    let notReadMy = ""
    if (message.Uid == connectData.Uid){
        classes = "chat-content-msg-my-msg";
        notReadMy = '<div class="chat-msg-not-read-my"></div>'
    } else {
        classes = "chat-msg-not-read chat-msg-not-read-obs";
    }
    chat_content.innerHTML += `
               <div data-msgid="${message.Msg.Id}" class="chat-content-msg ${classes}">
                   ${notReadMy}
                   <div class="chat-content-msg-text">
                       ${message.Msg.Text}
                   </div>
                   <div class="chat-content-msg-date">${message.Msg.Date}</div>
               </div>`;
    observeMessages(connectData);
    runLazyLoadMsg(connectData);
    runLazyLoadNotReadMsg(connectData);
}

function handleWsReadMsg(message: IMessage, ws: WebSocket){
    const element = document.querySelector(`[data-msgid="${message.Msg.Id}"]`);
    if (element.classList.contains("chat-content-msg-my-msg") && element.querySelectorAll(".chat-msg-not-read-my").length > 0){
        element.querySelectorAll(".chat-msg-not-read-my")[0].remove();
    }
    if (element.classList.contains("chat-msg-not-read")){
        element.classList.remove("chat-msg-not-read");
    }
}

function handleWsError(message: IMessage, ws: WebSocket){
    handleError(message.Msg.Error)
}

function handleWsImageMsg(message: IMessage, ws: WebSocket){
    const chat_content = document.getElementById("chat-content");
    let classes = ""
    let notReadMy = ""
    if (message.Uid == connectData.Uid){
        classes = "chat-content-msg-my-msg";
        notReadMy = '<div class="chat-msg-not-read-my"></div>'
    } else {
        classes = "chat-msg-not-read chat-msg-not-read-obs";
    }
    let chatImages = ""
    for (const image of message.Msg.Images.split("\\")) {
        chatImages += `<span>
                           <img src="/${image}">
                       </span>`
    }
    chat_content.innerHTML += `
               <div data-msgid="${message.Msg.Id}" class="chat-content-msg ${classes}">
                   ${notReadMy}
                    <div class="chat-content-msg-images">
                        ${chatImages}
                    </div>
                   <div class="chat-content-msg-text">
                       ${message.Msg.Text}
                   </div>
                   <div class="chat-content-msg-date">${message.Msg.Date}</div>
               </div>`;
    observeMessages(connectData);
    runLazyLoadMsg(connectData);
    runLazyLoadNotReadMsg(connectData);
}

export function RunWs(): IConnectData{
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
        const msg: IMessage = JSON.parse(event.data);
        MessageActions[msg.Type](msg, socket);
    });

    socket.addEventListener("close", (event) => {
        console.log("Close.");
    });

    socket.addEventListener("error", (event) => {
        console.error("Error:", event);
    });

    return connectData;
}
