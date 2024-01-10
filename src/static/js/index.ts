import '../scss/style.scss';
import {runIfExist} from "./utils";
import {RunWs} from "./chat/chatws/chat_ws";
import {runLazyLoadMsg, runLazyLoadNotReadMsg} from "./chat/lazy_load_msg";
import {RunWsNotification} from "./notification_ws";
import {OnImagesSelect, SendAjaxChatMessage} from "./chat/chat";
import {PopUp} from "./pop_up";
import {searchAjax} from "./search";
import {messageAjaxListen, messageMenu} from "./chat/message_menu";
import {resizeChatDetailImages} from "./chat/chat_detail";
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

interface ChatImage{
    Id: string
    Path: string
}

const regexChatDetail = /^\/chat-detail\/\d+$/;
if (regexChatDetail.test(window.location.pathname)){
    resizeChatDetailImages();
    let lload = new LazyLoad("last-image", ["imageid"], "/load-images")
    lload.run(function (response) {
        if (response["error"]){
            console.log(response["error"]);
            return
        }
        let chat_detail_images = document.getElementById("chat-detail-images");
        let images = response["images"] as Array<ChatImage>
        for (let i = 0; i < images.length; i++) {
            compressAndDisplayImage("http://localhost:8000/" + images[i].Path, 300)
                .then((compressedImageData) => {
                    let d = document.createElement("div")
                    d.classList.add("chat-detail-image");
                    if (i == images.length-1){
                        d.classList.add("last-image");
                    }
                    d.dataset.imageid = images[i].Id;
                    let img = document.createElement("img")
                    img.loading = "lazy"
                    img.src = compressedImageData;
                    d.appendChild(img);
                    chat_detail_images.appendChild(d);
                    resizeChatDetailImages();
                })
                .catch((error) => {
                    console.error('Error:', error);
                });
        }
    }, function (){

    })
}

function compressAndDisplayImage(imagePath: string, maxSize: number): Promise<string> {
    return new Promise((resolve, reject) => {
        const img = new Image();

        img.onload = () => {
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');

            if (!ctx) {
                reject(new Error('Canvas context is not supported'));
                return;
            }

            let width = img.width;
            let height = img.height;

            if (width > height) {
                if (width > maxSize) {
                    height *= maxSize / width;
                    width = maxSize;
                }
            } else {
                if (height > maxSize) {
                    width *= maxSize / height;
                    height = maxSize;
                }
            }

            canvas.width = width;
            canvas.height = height;
            ctx.drawImage(img, 0, 0, width, height);

            const compressedImage = canvas.toDataURL('image/jpeg', 0.7); // изменяем формат и качество изображения
            resolve(compressedImage);
        };

        img.onerror = () => {
            reject(new Error('Failed to load the image'));
        };
        img.src = imagePath;
    });
}
