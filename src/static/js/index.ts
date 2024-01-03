import '../scss/style.scss';
import {runIfExist} from "./utils";
import {RunWs} from "./chat/chatws/chat_ws";
import {runLazyLoadMsg, runLazyLoadNotReadMsg} from "./chat/lazy_load_msg";
import {RunWsNotification} from "./notification_ws";
import {OnImagesSelect, SendAjaxChatMessage} from "./chat/chat";
import {PopUp} from "./pop_up";
import {searchAjax} from "./search";

runIfExist(document.getElementById("pp-del-avatar-label"), function (el){
   el.onclick = function (){
       document.getElementById("pp-del-avatar").classList.toggle("page-panel-checkbox-true");
   }
});

let popUpNotification = new PopUp("notification-pop-up", false);
popUpNotification.onClick(function (popUp, popUpActivate){
    popUpActivate.classList.toggle("hide-pop-up-activate");
    popUp.classList.toggle("pop-up-hide");
});
popUpNotification.start();

RunWsNotification();
const regex = /^\/chat\/\d+$/;
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

const regexProf = /^\/prof\/\d+$/;
if (regexProf.test(window.location.pathname)){
    let popUpNotification = new PopUp("pop-up-panel-new-chat", true);
    popUpNotification.onClick(function (popUp, popUpActivate){
        popUpActivate.classList.toggle("hide-pop-up-activate-new-chat");
        popUp.classList.toggle("pop-up-hide");
    });
    popUpNotification.start();
}

searchAjax();
