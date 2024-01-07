import {LazyLoad} from "../lazy_load";
import {observeMessages} from "./observe_messages";
import {IConnectData} from "./chatws/chat_ws";
import {handleError} from "./chat";
import {messageMenu} from "./message_menu";

export function runLazyLoadMsg(connectData: IConnectData){
    let l = new LazyLoad("last-msg", ["first", "msgtype", "uid", "chatid", "msgid"], "/load-messages")
    l.setOptions({
        root: null,
        threshold: 0.5,
    });
    l.setValue("handler", "read");
    l.run(function (response){
        const parentElement = document.getElementById('chat-content');
        if (response["err"]){
            handleError(response["err"])
            return
        }
        const scrollTopBefore = parentElement.scrollTop;
        const scrollHeightBefore = parentElement.scrollHeight;
        let type = response["type"];
        if (response["first"] == 1){
            type = "read";
        }
        if (response["messages"] && type == "read") {
            let msg = response["messages"]
            for (let i = 0; i < msg.length; i++) {
                let msgData: MsgDynamicData = {
                    lastMsgData: `data-first="0" data-msgtype="${type}" data-msgid="${msg[i].Id}"`,
                    isReadMy: "",
                    classes: "",
                }
                if (i == msg.length-1){
                    msgData.lastMsgData += `data-chatid="${response["chatId"]}"`;
                    msgData.classes += "last-msg";
                }
                if (response["uid"] == msg[i].UserId){
                    msgData.classes += " chat-content-msg-my-msg";
                }
                if (response["uid"] != msg[i].UserId && msg[i]["IsRead"] == "0") {
                    continue;
                }
                if (msg[i]["IsRead"] == "0"){
                    msgData.isReadMy = '<div class="chat-msg-not-read-my"></div>';
                }
                document.getElementById("chat-content").insertAdjacentHTML('afterbegin', getMsgText(msg[i], msgData));
            }
            observeMessages(connectData)
            messageMenu();
        } else {
            observeMessages(connectData)
            messageMenu();
        }
        const scrollHeightAfter = parentElement.scrollHeight;
        parentElement.scrollTop = scrollTopBefore + (scrollHeightAfter - scrollHeightBefore);
    }, function (error) {
        console.log(error);
    })
}

export function runLazyLoadNotReadMsg(connectData: IConnectData) {
    let l = new LazyLoad("chat-msg-not-read-obs-down", ["first", "msgtype", "uid", "chatid", "msgid"], "/load-messages")
    l.setOptions({
        root: null,
        threshold: 0.5,
    });
    l.setValue("handler", "notread");
    l.run(function (response){
        const parentElement = document.getElementById('chat-content');
        if (response["err"]){
            handleError(response["err"])
            return
        }
        const scrollTopBefore = parentElement.scrollTop;
        let type = response["type"];
        if (response["first"] == 1){
            type = "notread";
        }
        if (response["messages"] && type == "notread"){
            let msg = response["messages"]
            for (let i = 0; i < msg.length; i++) {
                let msgData: MsgDynamicData = {
                    lastMsgData: `data-first="0" data-msgtype="${type}" data-msgid="${msg[i].Id}"`,
                    isReadMy: "",
                    classes: "",
                }
                if (i == msg.length-1){
                    msgData.lastMsgData += `data-chatid="${response["chatId"]}"`;
                    msgData.classes += "chat-msg-not-read-obs-down";
                }
                if (response["uid"] == msg[i].UserId){
                    msgData.classes += " chat-content-msg-my-msg";
                }
                if (response["uid"] != msg[i].UserId && msg[i]["IsRead"] == "0") {
                    msgData.classes += " chat-msg-not-read chat-msg-not-read-obs"
                }
                if (msg[i]["IsRead"] == "0"){
                    msgData.isReadMy = '<div class="chat-msg-not-read-my"></div>';
                }
                document.getElementById("chat-content").insertAdjacentHTML('beforeend', getMsgText(msg[i], msgData));
            }
            observeMessages(connectData);
            messageMenu();
        } else {
            observeMessages(connectData);
            messageMenu();
        }
        parentElement.scrollTop = scrollTopBefore;
    }, function (error){
        console.log(error);
    })
}

interface MsgDynamicData{
    lastMsgData: string;
    isReadMy: string;
    classes: string;
}

function getMsgText(msg, msgData: MsgDynamicData){
    let chatImages = "";
    if (msg.Images != null){
        for (const image of msg.Images) {
            chatImages += `<span>
                           <img src="/${image.Path}">
                       </span>`;
        }
    }
    return `
        <div ${msgData.lastMsgData} class="chat-content-msg ${msgData.classes}">
            <div class="message-menu message-menu-hide">
                <button class="message-menu-delete" type="button"><a href="#">
                    <img src="/static/img/del.svg">
                </a></button>
                <button class="message-menu-update" type="button"><a href="#">
                    <img src="/static/img/edit.svg">
                </a></button>
            </div>
            ${msgData.isReadMy}
            <div class="chat-content-msg-images">
                ${chatImages}
            </div>
            <div class="chat-content-msg-text">
                ${ msg.Text }
            </div>
            <div class="chat-content-msg-date">${ msg.Date }</div>
        </div>`
}
