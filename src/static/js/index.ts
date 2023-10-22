import '../scss/style.scss';
import {runIfExist} from "./utils";

runIfExist(document.getElementById("header-user"), function (el) {
    el.onclick = function (){
        document.getElementById("hu-menu").classList.toggle("hu-menu-hidden");
    }
});
