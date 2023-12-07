import {Msg} from "./chat_ws";

export enum NotificationType{
    Connect,
    GlobalIncrementMsg,
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
                let notificationCount = document.getElementById("notification-count");
                let number = parseInt(notificationCount.innerText);
                if (typeof number == "number"){
                    number++;
                    notificationCount.innerText = String(number);
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
