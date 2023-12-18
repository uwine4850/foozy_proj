import '../scss/style.scss';
import {runIfExist} from "./utils";
import {Ajax} from "./ajax";
import {RunWs} from "./chat_ws";
import {runLazyLoadMsg, runLazyLoadNotReadMsg} from "./lazy_load_msg";
import {RunWsNotification} from "./notification_ws";

// runIfExist(document.getElementById("header-user"), function (el) {
//     el.onclick = function (){
//         document.getElementById("hu-menu").classList.toggle("hu-menu-hidden");
//     }
// });
//
// runIfExist(document.getElementById("pc-menu-filter-value"), function (el){
//    el.onclick = function (){
//         document.getElementById("filter-items").classList.toggle("filter-items-hidden");
//    }
// });

runIfExist(document.getElementById("pp-del-avatar-label"), function (el){
   el.onclick = function (){
       document.getElementById("pp-del-avatar").classList.toggle("page-panel-checkbox-true");
   }
});

// Subscribe
// let subscribeAjax = new Ajax("/subscribe-post", "subscribe-post");
// subscribeAjax.onSuccess(function (response: string) {
//     let btn = document.getElementById("profile-button-subscribe");
//     let subscribers = document.getElementById("subscribers-value");
//     let subscribersValue = parseInt(subscribers.innerHTML);
//     btn.classList.toggle("profile-button-unsubscribe");
//     if (btn.classList.contains("profile-button-unsubscribe")){
//         subscribersValue++;
//         btn.innerHTML = '<a><img src="/static/img/subscribe.svg">Unsubscribe</a>';
//     } else {
//         if (subscribersValue > 0){
//             subscribersValue--;
//         }
//         btn.innerHTML = '<a><img src="/static/img/subscribe.svg">Subscribe</a>';
//     }
//     subscribers.innerHTML = String(subscribersValue);
// });
// subscribeAjax.listen()

const regex = /^\/chat\/\d+$/;
RunWsNotification();
if (regex.test(window.location.pathname)){
    let connectData = RunWs();
    runLazyLoadMsg(connectData);
    runLazyLoadNotReadMsg(connectData);
}
