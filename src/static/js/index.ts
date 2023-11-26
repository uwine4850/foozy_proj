import '../scss/style.scss';
import {runIfExist} from "./utils";
import {Ajax} from "./ajax";
import {RunWs} from "./chat_ws";
import {handleIntersection, LazyLoad} from "./lazy_load";

runIfExist(document.getElementById("header-user"), function (el) {
    el.onclick = function (){
        document.getElementById("hu-menu").classList.toggle("hu-menu-hidden");
    }
});

runIfExist(document.getElementById("pc-menu-filter-value"), function (el){
   el.onclick = function (){
        document.getElementById("filter-items").classList.toggle("filter-items-hidden");
   }
});

runIfExist(document.getElementById("pp-del-avatar-label"), function (el){
   el.onclick = function (){
       document.getElementById("pp-del-avatar").classList.toggle("page-panel-checkbox-true");
   }
});

// Subscribe
let subscribeAjax = new Ajax("/subscribe-post", "subscribe-post");
subscribeAjax.onSuccess(function (response: string) {
    let btn = document.getElementById("profile-button-subscribe");
    let subscribers = document.getElementById("subscribers-value");
    let subscribersValue = parseInt(subscribers.innerHTML);
    btn.classList.toggle("profile-button-unsubscribe");
    if (btn.classList.contains("profile-button-unsubscribe")){
        subscribersValue++;
        btn.innerHTML = '<a><img src="/static/img/subscribe.svg">Unsubscribe</a>';
    } else {
        if (subscribersValue > 0){
            subscribersValue--;
        }
        btn.innerHTML = '<a><img src="/static/img/subscribe.svg">Subscribe</a>';
    }
    subscribers.innerHTML = String(subscribersValue);
});
subscribeAjax.listen()

const regex = /^\/chat\/\d+$/;
if (regex.test(window.location.pathname)){
    RunWs();
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
                    let lastMsgData: string;
                    let myMsg: string;
                    if (i == msg.length-1){
                        lastMsgData = `data-chatid="${response["chatId"]}" data-msgid="${msg[i].Id}" id="last-msg"`;
                    }
                    if (response["uid"] == msg[i].UserId){
                        myMsg = "chat-content-msg-my-msg"
                    }
                    let _msg =`
                        <div ${lastMsgData} class="chat-content-msg ${myMsg}">
                            <div class="chat-content-msg-text">
                                ${ msg[i].Text }
                            </div>
                            <div class="chat-content-msg-date">${ msg[i].Date }</div>
                        </div>`
                    document.getElementById("chat-content").insertAdjacentHTML('afterbegin', _msg);
                }
            }
            const scrollHeightAfter = parentElement.scrollHeight;
            parentElement.scrollTop = scrollTopBefore + (scrollHeightAfter - scrollHeightBefore);
        }, function (error) {
            console.log(error);
        })
    }

}
