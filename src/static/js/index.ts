import '../scss/style.scss';
import {runIfExist} from "./utils";
import {Ajax} from "./ajax";
import {RunWs} from "./chat_ws";
import {runLazyLoadMsg, runLazyLoadNotReadMsg} from "./lazy_load_msg";
import {RunWsNotification} from "./notification_ws";
import {OnImagesSelect, SendAjaxChatMessage} from "./chat/chat";
import {PopUp} from "./pop_up";

runIfExist(document.getElementById("pp-del-avatar-label"), function (el){
   el.onclick = function (){
       document.getElementById("pp-del-avatar").classList.toggle("page-panel-checkbox-true");
   }
});

let popUpNotification = new PopUp("notification-pop-up", false);
popUpNotification.onClick(function (popUp){
    document.getElementById("pop-up-activate").classList.toggle("hide-pop-up-activate");
    popUp.classList.toggle("pop-up-hide");
});
popUpNotification.start();

const regex = /^\/chat\/\d+$/;
RunWsNotification();
if (regex.test(window.location.pathname)){
    let popUp = new PopUp("chat-pop-up", true);
    popUp.onClick(function (popUp){
        popUp.classList.toggle("pop-up-hide");
    });
    popUp.start();

    let connectData = RunWs();
    OnImagesSelect();
    runLazyLoadMsg(connectData);
    runLazyLoadNotReadMsg(connectData);
}
