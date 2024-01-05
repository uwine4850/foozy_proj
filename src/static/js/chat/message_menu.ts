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
        message.classList.toggle("selected-message");
        message.querySelector(".message-menu").classList.toggle("message-menu-hide");
    }, 500);
}

function handleMouseUp() {
    clearTimeout(timeoutId);
}

function messageMenuActions(){
    let deleteButtons = document.getElementsByClassName("message-menu-delete") as HTMLCollectionOf<HTMLElement>;
    let deletePopup = document.getElementById("delete-pop-up-activate");
    for (let i = 0; i < deleteButtons.length; i++) {
        deleteButtons[i].onclick = function (){
            deletePopup.click();
        }
    }
}

export function messageAjaxListen(){
    let ajaxDelete = new Ajax("/message-menu", "delete-pop-up-content");
    ajaxDelete.setUrlFormValue("action", "delete")
    ajaxDelete.onSuccess(function (response){
        handleError(response["error"]);
    })
    ajaxDelete.listen();
    let submit_delete_message = document.getElementById("submit-delete-message") as HTMLButtonElement;
    submit_delete_message.onclick = function (){
        document.getElementById("delete-pop-up-activate").click();
    }
}
