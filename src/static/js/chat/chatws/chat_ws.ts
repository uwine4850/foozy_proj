import {handleError, SendAjaxChatMessage} from "../chat";
import {HandleWsTextMsg} from "./ws_text_message";
import {HandleWsReadMsg} from "./ws_read_message";
import {HandleWsImageMsg} from "./ws_image_message";

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
    HandleWsTextMsg(message, ws, connectData);
}

function handleWsReadMsg(message: IMessage, ws: WebSocket){
    HandleWsReadMsg(message, ws);
}

function handleWsError(message: IMessage, ws: WebSocket){
    handleError(message.Msg.Error);
}

function handleWsImageMsg(message: IMessage, ws: WebSocket){
    HandleWsImageMsg(message, ws, connectData);
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
