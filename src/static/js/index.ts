import '../scss/style.scss';
import {runIfExist} from "./utils";

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
