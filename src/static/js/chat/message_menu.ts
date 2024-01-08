import {Ajax} from "../ajax";
import {handleError} from "./chat";

export function messageMenu(){
    let message = document.getElementsByClassName("chat-content-msg") as HTMLCollectionOf<HTMLElement>;
    for (let i = 0; i < message.length; i++) {
        message[i].onmousedown = handleMouseDown;
        message[i].onmouseup = handleMouseUp;
    }
    messageMenuActions();
}

let timeoutId: ReturnType<typeof setTimeout>;

function handleMouseDown(e: MouseEvent) {
    if (e.button != 0){
        return
    }
    timeoutId = setTimeout(() => {
        let targetElement = e.target as HTMLElement;
        if (targetElement.classList.length == 0){
            return
        }
        let message: HTMLElement;
        if (!targetElement.classList.contains("chat-content-msg")){
            message = targetElement.parentElement;
        } else {
            message = targetElement;
        }
        let msgid = message.dataset.msgid;
        let del_msg_input  = document.getElementById("del-msg-id") as HTMLInputElement;
        del_msg_input.value = msgid;
        setUpdMessageData(message);
        message.classList.toggle("selected-message");
        message.querySelector(".message-menu").classList.toggle("message-menu-hide");
    }, 500);
}

function handleMouseUp() {
    clearTimeout(timeoutId);
}

// Trigger events when menu buttons are pressed.
function messageMenuActions(){
    let deleteButtons = document.getElementsByClassName("message-menu-delete") as HTMLCollectionOf<HTMLElement>;
    let deletePopup = document.getElementById("delete-pop-up-activate");
    for (let i = 0; i < deleteButtons.length; i++) {
        deleteButtons[i].onclick = function (){
            deletePopup.click();
        }
    }
    let updateButtons = document.getElementsByClassName("message-menu-update") as HTMLCollectionOf<HTMLElement>;
    let updatePopup = document.getElementById("update-pop-up-activate");
    for (let i = 0; i < updateButtons.length; i++) {
        updateButtons[i].onclick = function (){
            updatePopup.click();
        }
    }
}

// Starts listening for the message menu form to be sent.
export function messageAjaxListen(){
    // Sending a form delete message.
    let ajaxDelete = new Ajax("/message-menu", "delete-pop-up-content");
    ajaxDelete.setUrlFormValue("action", "delete")
    ajaxDelete.onSuccess(function (response){
        if (response["error"]){
            handleError(response["error"]);
        }
    })
    ajaxDelete.listen();
    let submit_delete_message = document.getElementById("submit-delete-message") as HTMLButtonElement;
    submit_delete_message.onclick = function (){
        document.getElementById("delete-pop-up-activate").click();
    }

    // Sending a message update form.
    let ajaxUpdate = new Ajax("/message-menu", "update-pop-up-content");
    ajaxUpdate.setUrlFormValue("action", "update")
    ajaxUpdate.onSuccess(function (response){
        if (response["error"]){
            handleError(response["error"]);
        }
    });
    ajaxUpdate.listen();
    let submit_update_message = document.getElementById("submit-update-message") as HTMLButtonElement;
    submit_update_message.onclick = function (){
        document.getElementById("update-pop-up-activate").click();
    }
}

// Sets all the required data for the message update form.
function setUpdMessageData(message: HTMLElement){
    let updMsgId = document.getElementById("updMsgId") as HTMLInputElement;
    updMsgId.value = message.dataset.msgid;

    let msgImages = message.querySelector(".chat-content-msg-images");
    let msgImagesData = msgImages.querySelectorAll("span");

    // Setting images in the form.
    // It displays the image itself with a button to delete it and hidden checkboxes to send deletion data to the server.
    // Each checkbox name starts with updRmImage. and then the path to the image.
    if (msgImagesData){
        let chatUpdImages = document.getElementById("chat-upd-images");
        chatUpdImages.innerHTML = "";
        let updImagesCheckboxes = document.getElementById("upd-images-checkboxes");
        updImagesCheckboxes.innerHTML = "";
        for (let i = 0; i < msgImagesData.length; i++) {
            let img = msgImagesData[i].querySelector("img");
            let path = img.src.split("http://localhost:8000/")[1];
            chatUpdImages.innerHTML +=  `
                <div class="chat-upd-image">
                    <button data-updImagePath="/${path}" type="button" class="chat-upd-image-remove"><a>
                        <span>Remove</span>
                    </a></button>
                    <img src="/${path}" alt="">
                </div>
            `;
            updImagesCheckboxes.innerHTML += `
            <input type="checkbox" name="updRmImage./${path}">
            `
        }
        let chatUpdImageRemove = document.getElementsByClassName("chat-upd-image-remove") as HTMLCollectionOf<HTMLElement> ;
        let updImageCheckboxes = document.getElementById("upd-images-checkboxes").querySelectorAll("input");
        for (let i = 0; i < chatUpdImageRemove.length; i++) {
            chatUpdImageRemove[i].onclick = function (){
                chatUpdImageRemove[i].classList.toggle("chat-upd-image-remove-checked");
                updImageCheckboxes[i].click();
            }
        }
    }
    // Set message text.
    let a = document.getElementById("upd-message-text") as HTMLTextAreaElement;
    a.value = message.querySelector(".chat-content-msg-text").innerHTML.trim();
}
