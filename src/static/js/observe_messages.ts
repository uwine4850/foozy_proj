import {Observer} from "./observer";
import {ConnectData, Msg, MsgType} from "./chat_ws";

export function observeMessages(connectData: ConnectData){
    let ob = new Observer("chat-msg-not-read-obs", {root: null, threshold: 0.7});
    ob.run(function (entry) {
        let msgId = entry.target.getAttribute("data-msgid")
        if (msgId == null){
            return
        }
        let m: Msg = {
            Type: MsgType.ReadMsg,
            Uid: connectData.Uid,
            ChatId: connectData.ChatId,
            Msg: {"Id": msgId}
        }
        setTimeout(() => {
            connectData.Socket.send(JSON.stringify(m));
        }, 100);
    })
}