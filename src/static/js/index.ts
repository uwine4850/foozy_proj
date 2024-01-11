import '../scss/style.scss';
import {runIfExist} from "./utils";
import {RunWs} from "./chat/chatws/chat_ws";
import {runLazyLoadMsg, runLazyLoadNotReadMsg} from "./chat/lazy_load_msg";
import {RunWsNotification} from "./notification_ws";
import {OnImagesSelect, SendAjaxChatMessage} from "./chat/chat";
import {PopUp} from "./pop_up";
import {searchAjax} from "./search";
import {messageAjaxListen, messageMenu} from "./chat/message_menu";
import {detailLoadImages, resizeChatDetailImages, searchMessages, switchPages} from "./chat/chat_detail";
import {Observer} from "./observer";
import {LazyLoad} from "./lazy_load";

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
    messageMenu();
    let popUp = new PopUp("p1", true);
    popUp.onClick(function (popUp){
        popUp.classList.toggle("pop-up-hide");
    });
    popUp.start();

    let popUpDelete = new PopUp("p2", true);
    popUpDelete.onClick(function (popUp){
        popUp.classList.toggle("pop-up-hide");
    });
    popUpDelete.start();

    let popUpUpdate = new PopUp("p3", true);
    popUpUpdate.onClick(function (popUp){
        popUp.classList.toggle("pop-up-hide");
    });
    popUpUpdate.start();

    messageAjaxListen();

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

const regexChatDetail = /^\/chat-detail\/\d+$/;
if (regexChatDetail.test(window.location.pathname)){
    switchPages();
    resizeChatDetailImages();
    detailLoadImages();
    searchMessages();
}
