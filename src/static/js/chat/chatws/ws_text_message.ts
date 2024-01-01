import {observeMessages} from "../observe_messages";
import {runLazyLoadMsg, runLazyLoadNotReadMsg} from "../lazy_load_msg";
import {IConnectData, IMessage} from "./chat_ws";

export function HandleWsTextMsg(message: IMessage, ws: WebSocket, connectData: IConnectData){
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