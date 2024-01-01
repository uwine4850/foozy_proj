export enum NotificationType{
    WsConnect,
    WsError,
    WsIncrementChatMsgCount,
    WsGlobalIncrementMsg,
    WsGlobalDecrementMsg,
    WsPopUpMessage,
    WsNewChat,
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
    [NotificationType.WsPopUpMessage]: handlerWsPopUpMessage,
    [NotificationType.WsNewChat]: handlerWsNewChat,
}

function handlerWsNewChat(notification: INotification, ws: WebSocket){
    if (window.location.pathname != "/chat-list"){
        return
    }
    let user = JSON.parse(notification.Msg.User);
    let avatar = "/static/img/default.jpeg";
    if (user.Avatar != ""){
        avatar = user.Avatar;
    }
    let chat_list = document.getElementById("chat-list");
    chat_list.insertAdjacentHTML("afterbegin", `
        <a data-chatid="${notification.Msg.ChatId}" href="/chat/${notification.Msg.ChatId}" class="chat-list-item">
        <div class="chat-list-user">
            <div class="chat-list-user-avatar">
                <img src="${avatar}" alt="">
            </div>
            <div class="chat-list-user-data">
                <div class="chat-list-user-username">@${user.Username}</div>
                <div id="chat-list-user-msg" class="chat-list-user-msg">${notification.Msg.Text}</div>
            </div>
        </div>
        <div id="chat-list-user-msg-count" class="chat-list-user-msg-count"></div>
    </a>
    `)
}

function handlerWsPopUpMessage(notification: INotification, ws: WebSocket){
    const regex = /^\/chat\/\d+$/;
    if (!regex.test(window.location.pathname)){
        messagePopUpNotification(notification.Msg);
    }
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
    notificationPopUp(notification.Msg.error);
}

function handlerWsIncrementChatMsgCount(notification: INotification, ws: WebSocket){
    if (window.location.pathname === "/chat-list"){
        const chat = document.querySelector(`[data-chatid="${notification.Msg.chatId}"]`);
        let countEl = chat.querySelector("#chat-list-user-msg-count");
        let num = 0;
        if (countEl.innerHTML.trim() != ""){
            num = parseInt(countEl.innerHTML);
        }
        if (typeof num == "number"){
            num++;
        }
        countEl.innerHTML = String(num);
        // Change last message text
        chat.querySelector("#chat-list-user-msg").innerHTML = notification.Msg.Text;
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
    let notificationCount = document.getElementById("chat-list-msg-count");
    let number = 0;
    if (notificationCount.innerText.trim() != ""){
        number = parseInt(notificationCount.innerText);
    }
    if (increment && typeof number == "number"){
        number++;
    }
    if (!increment && typeof number == "number"){
        number--;
        if (number == 0){
            notificationCount.innerText = "";
            return
        }
    }
    notificationCount.innerText = String(number);
}

function notificationPopUp(text: string){
    document.getElementById("notification-pop-up-content").innerHTML = text;
    document.getElementById("pop-up-activate").click();
}

function messagePopUpNotification(messageData: Record<string, string>){
    document.getElementById("npp-user-avatar").innerHTML = "";
    document.getElementById("npp-user-name").innerHTML = "";
    document.getElementById("npp-msg-images").innerHTML = "";
    document.getElementById("npp-message-text").innerHTML = "";

    const userJsonData = JSON.parse(messageData.User);
    if (userJsonData.Avatar){
        document.getElementById("npp-user-avatar").innerHTML = `<img src="${userJsonData.Avatar}">`;
    } else {
        document.getElementById("npp-user-avatar").innerHTML = `<img src="/static/img/default.jpeg">`;
    }
    document.getElementById("npp-user-name").innerHTML = userJsonData.Name;
    if (messageData.Images){
        let images = messageData.Images.split("\\");
        for (let i = 0; i < images.length; i++) {
            if (i < 3){
                document.getElementById("npp-msg-images").innerHTML += `
            <span>
                <img src="${images[i]}">
            </span>
            `
            } else {
                document.getElementById("npp-message-more-images").classList.remove("npp-message-more-images-hide");
                document.getElementById("npp-message-more-images").innerHTML = `+${images.length-i} images`;
                break;
            }
        }
    }
    document.getElementById("npp-message-text").innerHTML = messageData.Text;
    if (document.getElementById("message-pop-up-activate").classList.contains("hide-pop-up-activate")){
        document.getElementById("message-pop-up-activate").click();
    }
}
