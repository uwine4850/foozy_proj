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