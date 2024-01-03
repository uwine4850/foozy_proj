import {Ajax} from "./ajax";

export function searchAjax(){
    if (window.location.pathname != "/search"){
        return
    }
    let ajax = new Ajax("/search-post", "search-input");
    ajax.onSuccess(function (response){
        const search_error = document.getElementById("search-error");
        search_error.innerHTML = "";
        if (!search_error.classList.contains("search-error-hidden")){
            search_error.classList.add("search-error-hidden");
        }
        if (response["error"]){
            console.log(response["error"]);
            search_error.innerHTML = response["error"];
            search_error.classList.toggle("search-error-hidden");
            return
        }
        const search_items = document.getElementById("search-items");
        let users = response["users"] as Array<any>
        for (let i = 0; i < users.length; i++) {
            let user = users[i];
            let avatar = "/static/img/default.jpeg";
            if (user.Avatar != ""){
                avatar = user.Avatar;
            }
            search_items.innerHTML += `
               <a href="/prof/${user.Id}" class="search-item">
                   <div class="search-item-image">
                       <img src="${avatar}">
                   </div>
                   <div class="search-item-username">
                       ${user.Username}
                   </div>
               </a>
               `
        }
    })
    ajax.onError(function (error) {
        console.log(error);
    })
    ajax.listen();
}