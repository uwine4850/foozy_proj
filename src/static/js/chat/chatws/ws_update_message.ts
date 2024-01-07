import {IConnectData, IMessage} from "./chat_ws";

export function HandleWsUpdateMsg(message: IMessage, ws: WebSocket, connectData: IConnectData){
    let msgId = message.Msg.id;
    const htmlMessage = document.querySelector(`[data-msgid="${msgId}"]`);
    if (message.Msg.text){
        htmlMessage.querySelector(".chat-content-msg-text").innerHTML = message.Msg.text;
    }
    if (message.Msg.delImages){
        let htmlImages = htmlMessage.querySelector(".chat-content-msg-images");
        for (const imagePath of message.Msg.delImages.split("\\")) {
            let s = htmlImages.querySelectorAll("span");
            for (let i = 0; i < s.length; i++) {
                let img = s[i].querySelector("img")
                if (img.src == "http://localhost:8000/" + imagePath){
                    s[i].remove();
                }
            }

        }
    }
}