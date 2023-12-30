import {Ajax} from "../ajax";
import {IConnectData} from "../chat_ws";

export function SendAjaxChatMessage(wsMessageConnectData: IConnectData){
    let aj = new Ajax("/receive-msg", "chat-form");
    aj.onSuccess(function (){
        document.getElementById("chat-form-images").innerHTML = "";
        let input = document.getElementById("images") as HTMLInputElement;
        input.value = "";
        let textarea = document.getElementById("chat-textarea") as HTMLTextAreaElement;
        textarea.value = "";
    });
    aj.onError(function (error){
        console.log("Error when receive chat message: ", error);
    })
    aj.setMultipartFormValue("chatId", wsMessageConnectData.ChatId);
    aj.setMultipartFormValue("uid", wsMessageConnectData.Uid);
    aj.listen();
}

function removeSelectedFile(fileName: string) {
    const fileInput = document.getElementById('images') as HTMLInputElement;
    const newFiles = Array.from(fileInput.files!)
        .filter(file => file.name != fileName) as unknown as File[];

    const dataTransfer = new DataTransfer();
    newFiles.forEach(file => dataTransfer.items.add(file));
    fileInput.files = dataTransfer.files;
}

export function OnImagesSelect(){
    document.getElementById("images-btn")!.onclick = function (){
        document.getElementById('images')!.click();
    }
    document.getElementById('images')!.addEventListener('change', function(this: HTMLInputElement) {
        const fileList = document.getElementById('chat-form-images')!;
        fileList.innerHTML = '';

        const files = this.files!;
        for (let i = 0; i < files.length; i++) {
            let img = URL.createObjectURL(files[i])
            fileList.innerHTML += `
            <div class="chat-form-image">
                <button data-name="${files[i].name}" class="chat-form-remove-image">
                </button>
                <img src="${img}">
            </div>`
        }
        let removeButtons = document.getElementsByClassName("chat-form-remove-image");
        let removeButtonsArray = Array.from(removeButtons);
        removeButtonsArray.forEach(element => {
            element.addEventListener('click', (event) => {
                const el = event.target as HTMLElement;
                const name = el.dataset.name
                removeSelectedFile(name);
                el.parentElement.remove();
            });
        });
    });
}

export function handleError(text: string){
    document.getElementById("chat-pop-up-content").innerHTML = text;
    document.getElementById("pop-up-activate").click();
}
