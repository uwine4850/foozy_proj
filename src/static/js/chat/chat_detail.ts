import {LazyLoad} from "../lazy_load";
import {Ajax} from "../ajax";

export function switchPages(){
    let chat_detail_images_btn = document.getElementById("chat-detail-images-btn");
    let chat_detail_messages_btn = document.getElementById("chat-detail-messages-btn");
    let chat_detail_images = document.getElementById("chat-detail-images");
    chat_detail_images_btn.onclick = function (){
        chat_detail_images.classList.remove("chat-detail-images-hide");
        document.getElementById("chat-detail-search").style.display = "none";
    }
    chat_detail_messages_btn.onclick = function (){
        chat_detail_images.classList.add("chat-detail-images-hide");
        document.getElementById("chat-detail-search").style.display = "flex";
    }
}

export function resizeChatDetailImages(){
    let chat_detail_images = document.getElementById("chat-detail-images");
    let w = chat_detail_images.offsetWidth;
    let chat_detail_image = document.getElementsByClassName("chat-detail-image") as HTMLCollectionOf<HTMLElement>;
    let ww = 10;
    for (let i = 0; i < chat_detail_image.length; i++) {
        ww += chat_detail_image[i].offsetWidth + 10;
        if (ww > w){
            chat_detail_images.classList.add("chat-detail-images-grid");
        }
    }
}

interface ChatImage{
    Id: string
    Path: string
}

export function detailLoadImages(){
    let lload = new LazyLoad("last-image", ["imageid"], "/load-images");
    lload.run(function (response) {
        if (response["error"]){
            console.log(response["error"]);
            return
        }
        let chat_detail_images = document.getElementById("chat-detail-images");
        let images = response["images"] as Array<ChatImage>
        if (!images){
            return;
        }
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


export function searchMessages(){
    let ajax = new Ajax("/search-messages", "search-messages");
    ajax.onSuccess(function (response){
        if (response["error"]){
            console.log(response["error"]);
            return;
        }
        let uid = response["UID"];
        let messages = response["messages"];
        let chat_content_detail = document.getElementById("chat-content-detail");
        chat_content_detail.innerHTML = "";
        for (let i = 0; i < messages.length; i++) {
            let images = "";
            if (messages[i]["Images"]){
                for (let j = 0; j < messages[i]["Images"].length; j++) {
                    images += `
                            <span>
                                <img loading="lazy" src="/${messages[i]["Images"][j]["Path"]}">
                            </span>
                    `
                }
            }
            let myMsgClass = "";
            if (uid == messages[i]["UserId"]){
                myMsgClass = "chat-content-msg-my-msg"
            }
            chat_content_detail.innerHTML += `
                    <div class="chat-content-msg chat-content-msg-detail ${myMsgClass}">
                        <div class="chat-content-msg-images">
                            ${images}
                        </div>
                        <div class="chat-content-msg-text">
                            ${messages[i]["Text"]}
                        </div>
                        <div class="chat-content-msg-date">${messages[i]["Date"]}</div>
                    </div>  
            `;
        }
    })
    ajax.listen()
}
