import {IMessage} from "./chat_ws";

export function HandleWsReadMsg(message: IMessage, ws: WebSocket){
    const element = document.querySelector(`[data-msgid="${message.Msg.Id}"]`);
    if (element.classList.contains("chat-content-msg-my-msg") && element.querySelectorAll(".chat-msg-not-read-my").length > 0){
        element.querySelectorAll(".chat-msg-not-read-my")[0].remove();
    }
    if (element.classList.contains("chat-msg-not-read")){
        element.classList.remove("chat-msg-not-read");
    }
}