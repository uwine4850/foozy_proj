import {IConnectData, IMessage} from "./chat_ws";

export function HandleWsDeleteMessage(message: IMessage, ws: WebSocket, connectData: IConnectData){
    let msgId = message.Msg.msgId;
    const element = document.querySelector(`[data-msgid="${msgId}"]`);
    element.remove();
}