import {Observer} from "../observer";
import {IConnectData, IMessage, MessageType} from "./chatws/chat_ws";

export function observeMessages(connectData: IConnectData){
    let ob = new Observer("chat-msg-not-read-obs", {root: null, threshold: 0.7});
    ob.run(function (entry) {
        let msgId = entry.target.getAttribute("data-msgid")
        if (msgId == null){
            return
        }
        let m: IMessage = {
            Type: MessageType.WsReadMsg,
            Uid: connectData.Uid,
            ChatId: connectData.ChatId,
            Msg: {"Id": msgId}
        }
        setTimeout(() => {
            connectData.Socket.send(JSON.stringify(m));
        }, 100);
    })
}