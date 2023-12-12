export enum NotificationType{
    WsConnect,
    WsError,
    WsIncrementChatMsgCount,
    WsGlobalIncrementMsg,
    WsGlobalDecrementMsg,
}

export interface INotification{
    Type: number;
    UserIds: string[];
    Msg: Record<string, string>;
}

type MsgActionFunction = (notification: INotification, ws: WebSocket) => void;
type MsgActions = Record<NotificationType, MsgActionFunction>;

const MsgActions: MsgActions = {
    [NotificationType.WsConnect]: handlerConnect,
    [NotificationType.WsError]: handlerError,
    [NotificationType.WsIncrementChatMsgCount]: handlerWsIncrementChatMsgCount,
    [NotificationType.WsGlobalIncrementMsg]: handlerWsGlobalIncrementMsg,
    [NotificationType.WsGlobalDecrementMsg]: handlerWsGlobalDecrementMsg,
}

function handlerWsGlobalDecrementMsg(notification: INotification, ws: WebSocket){
    incrementNotificationCount(false);
}

function handlerWsGlobalIncrementMsg(notification: INotification, ws: WebSocket){
    incrementNotificationCount(true);
}

function handlerConnect(notification: INotification, ws: WebSocket){

}

function handlerError(notification: INotification, ws: WebSocket){
    console.log(notification.Msg.error);
}

function handlerWsIncrementChatMsgCount(notification: INotification, ws: WebSocket){
    if (window.location.pathname === "/chat-list"){
        const chat = document.querySelector(`[data-chatid="${notification.Msg.chatId}"]`);
        let countEl = chat.querySelector("#chat-list-user-msg-count");
        let num = parseInt(countEl.innerHTML);
        if (typeof num == "number"){
            num++;
        }
        countEl.innerHTML = String(num);
    }
}

export function RunWsNotification() {
    const ws = new WebSocket("ws://localhost:8000/notification-ws");
    ws.addEventListener("open", (event) => {
        console.log("Notification connect.");
    });
    ws.addEventListener("close", (event) => {

    });
    ws.onerror = function (ev){
        console.log("Error: ", ev);
    }
    ws.onmessage = function (ev){
        if (ev.data == ""){
            return
        }
        const msg: INotification = JSON.parse(ev.data);
        MsgActions[msg.Type](msg, ws);
    }
}

function incrementNotificationCount(increment: boolean){
    let notificationCount = document.getElementById("notification-count");
    let number = parseInt(notificationCount.innerText);
    if (increment && typeof number == "number"){
        number++;
    }
    if (!increment && typeof number == "number"){
        number--;
    }
    notificationCount.innerText = String(number);
}
