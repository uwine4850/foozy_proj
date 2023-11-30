import {LazyLoad} from "./lazy_load";
import {observeMessages} from "./observe_messages";
import {ConnectData} from "./chat_ws";

export function runLazyLoadMsg(connectData: ConnectData){
    let l = new LazyLoad("last-msg", ["first", "msgtype", "uid", "chatid", "msgid"], "/load-messages")
    l.setOptions({
        root: null,
        threshold: 0.5,
    });
    l.setValue("handler", "read");
    l.run(function (response){
        const parentElement = document.getElementById('chat-content');
        if (response["error"]){
            console.log(response["error"])
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
        } else {
            observeMessages(connectData)
        }
        const scrollHeightAfter = parentElement.scrollHeight;
        parentElement.scrollTop = scrollTopBefore + (scrollHeightAfter - scrollHeightBefore);
    }, function (error) {
        console.log(error);
    })
}

export function runLazyLoadNotReadMsg(connectData: ConnectData) {
    let l = new LazyLoad("chat-msg-not-read-obs-down", ["first", "msgtype", "uid", "chatid", "msgid"], "/load-messages")
    l.setOptions({
        root: null,
        threshold: 0.5,
    });
    l.setValue("handler", "notread");
    l.run(function (response){
        const parentElement = document.getElementById('chat-content');
        if (response["error"]){
            console.log(response["error"])
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
                    msgData.lastMsgData = `data-chatid="${response["chatId"]}"`;
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
                document.getElementById("chat-content").insertAdjacentHTML('beforeend', getMsgText(msg[0], msgData));
            }
            observeMessages(connectData)
        } else {
            observeMessages(connectData)
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
    return `
        <div ${msgData.lastMsgData} class="chat-content-msg ${msgData.classes}">
            ${msgData.isReadMy}
            <div class="chat-content-msg-text">
                ${ msg.Text }
            </div>
            <div class="chat-content-msg-date">${ msg.Date }</div>
        </div>`
}
