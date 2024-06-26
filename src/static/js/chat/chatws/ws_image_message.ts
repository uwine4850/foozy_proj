import {observeMessages} from "../observe_messages";
import {runLazyLoadMsg, runLazyLoadNotReadMsg} from "../lazy_load_msg";
import {IConnectData, IMessage} from "./chat_ws";
import {messageMenu} from "../message_menu";

export function HandleWsImageMsg(message: IMessage, ws: WebSocket, connectData: IConnectData){
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
                           <img loading="lazy" src="/${image}">
                       </span>`
    }
    chat_content.innerHTML += `
               <div data-msgid="${message.Msg.Id}" class="chat-content-msg ${classes}">
                   <div class="message-menu message-menu-hide">
                        <button class="message-menu-delete" type="button"><a href="#">
                            <img src="/static/img/del.svg">
                        </a></button>
                        <button class="message-menu-update" type="button"><a href="#">
                            <img src="/static/img/edit.svg">
                       </a></button>
                   </div>
                    ${notReadMy}
                    <div class="chat-content-msg-images">
                        ${chatImages}
                    </div>
                   <div class="chat-content-msg-text">
                       ${message.Msg.Text}
                   </div>
                   <div class="chat-content-msg-date">${message.Msg.Date}</div>
               </div>`;
    messageMenu();
    observeMessages(connectData);
    runLazyLoadMsg(connectData);
    runLazyLoadNotReadMsg(connectData);
}