import {LazyLoad} from "./lazy_load";
import {observeMessages} from "./observe_messages";
import {ConnectData} from "./chat_ws";

export function runLazyLoadMsg(connectData: ConnectData){
    if (document.getElementById('last-msg')) {
        let l = new LazyLoad("last-msg", ["uid", "chatid", "msgid"], "/load-messages")
        l.setOptions({
            root: null,
            threshold: 0.5,
        });
        l.run(function (response){
            const parentElement = document.getElementById('chat-content');
            if (response["error"]){
                console.log(response["error"])
                return
            }
            const scrollTopBefore = parentElement.scrollTop;
            const scrollHeightBefore = parentElement.scrollHeight;
            if (response["messages"]){
                let msg = response["messages"]
                for (let i = 0; i < msg.length; i++) {
                    let lastMsgData: string = `data-msgid="${msg[i].Id}"`;
                    let myMsg: string = "";
                    let isReadMy: string = "";
                    let isRead: string = "";
                    if (i == msg.length-1){
                        lastMsgData = `data-chatid="${response["chatId"]}" data-msgid="${msg[i].Id}" id="last-msg"`;
                    }
                    if (response["uid"] == msg[i].UserId){
                        myMsg = "chat-content-msg-my-msg"
                    }
                    if (response["uid"] != msg[i].UserId && msg[i]["IsRead"] == "0") {
                        isRead = "chat-msg-not-read chat-msg-not-read-obs";
                    }
                    if (msg[i]["IsRead"] == "0"){
                        isReadMy = '<div class="chat-msg-not-read-my"></div>';
                    }
                    let _msg =`
                        <div ${lastMsgData} class="chat-content-msg ${myMsg} ${isRead}">
                            ${isReadMy}
                            <div class="chat-content-msg-text">
                                ${ msg[i].Text }
                            </div>
                            <div class="chat-content-msg-date">${ msg[i].Date }</div>
                        </div>`
                    document.getElementById("chat-content").insertAdjacentHTML('afterbegin', _msg);
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
}