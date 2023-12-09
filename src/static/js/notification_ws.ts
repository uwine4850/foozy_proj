import {Msg} from "./chat_ws";

export enum NotificationType{
    Connect,
    GlobalIncrementMsg,
    GlobalDecrementMsg,
    IncrementMsg,
}

export interface Notification{
    Type: NotificationType;
    Uid: string[];
    Msg: Record<string, string>;
}

export interface NotificationConnData{
    Socket: WebSocket;
    Uid: string[];
}

let notificationData: NotificationConnData = {
    Socket: null,
    Uid: []
}

export function RunNotificationWS(): NotificationConnData{
    const socket = new WebSocket("ws://localhost:8000/notification-ws")
    socket.addEventListener("open", (event) => {
        console.log("Connect notification.");
        notificationData.Socket = socket
    });
    socket.addEventListener("close", (event) => {
        console.log("Close notification.");
    });

    socket.addEventListener("error", (event) => {
        console.error("Error:", event);
    });
    socket.onmessage = function (ev){
        if (ev.data == ""){
            return
        }
        const notification: Notification = JSON.parse(ev.data);
        switch (notification.Type){
            case NotificationType.Connect:
                break;
            case NotificationType.GlobalIncrementMsg:
                incrementNotificationCount(true);
                break;
            case NotificationType.GlobalDecrementMsg:
                incrementNotificationCount(false);
                break;
            case NotificationType.IncrementMsg:
                if (window.location.pathname === "/chat-list"){
                    const chat = document.querySelector(`[data-chatid="${notification.Msg.ChatId}"]`);
                    let countEl = chat.querySelector("#chat-list-user-msg-count");
                    let num = parseInt(countEl.innerHTML);
                    if (typeof num == "number"){
                        num++;
                    }
                    countEl.innerHTML = String(num);
                }
        }
    }
    return notificationData
}

export function globalIncrementMsgNotification(msg: Msg, sendToUsersId: string[],  notificationData: NotificationConnData){
    let n: Notification = {
        Type: NotificationType.GlobalIncrementMsg,
        Uid: sendToUsersId,
        Msg: {"ChatId": msg.ChatId}
    }
    notificationData.Socket.send(JSON.stringify(n));
}

export function incrementMsgNotification(msg: Msg, sendToUsersId: string[],  notificationData: NotificationConnData){
    let n: Notification = {
        Type: NotificationType.IncrementMsg,
        Uid: sendToUsersId,
        Msg: {"ChatId": msg.ChatId}
    }
    notificationData.Socket.send(JSON.stringify(n));
}

export function globalDecrementMsgNotification(msg: Msg, sendToUsersId: string[],  notificationData: NotificationConnData){
    let n: Notification = {
        Type: NotificationType.GlobalDecrementMsg,
        Uid: sendToUsersId,
        Msg: {"ChatId": msg.ChatId}
    }
    notificationData.Socket.send(JSON.stringify(n));
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
