(()=>{"use strict";var e={391:(e,t,n)=>{n.r(t)},872:(e,t)=>{Object.defineProperty(t,"__esModule",{value:!0}),t.Ajax=void 0;var n=function(){function e(e,t){this.path=e,this.formId=t}return e.prototype.onSuccess=function(e){this.onSuccessFn=e},e.prototype.onError=function(e){this.onErrorFn=e},e.prototype.listen=function(){var e=this;$(document).ready((function(){$("#"+e.formId).submit((function(t){t.preventDefault();var n=$(e).serialize();$.ajax({type:"POST",url:e.path,data:n,success:function(t){e.onSuccessFn(t)},error:function(t,n,s){e.onErrorFn(s)}})}))}))},e}();t.Ajax=n},862:(e,t)=>{var n,s,o;Object.defineProperty(t,"__esModule",{value:!0}),t.RunWs=void 0,function(e){e[e.Connect=0]="Connect",e[e.TextMsg=1]="TextMsg"}(n||(n={})),t.RunWs=function(){var e=document.getElementById("chat-textarea"),t=new WebSocket("ws://localhost:8000/chat-ws");t.addEventListener("open",(function(e){console.log("Connect.")})),t.addEventListener("message",(function(e){var t=JSON.parse(e.data);switch(t.Type){case n.TextMsg:t.Uid==s?document.getElementById("chat-content").innerHTML+=' <div class="chat-content-msg chat-content-msg-my-msg">\n                        <div class="chat-content-msg-text">\n                            '.concat(t.Msg.Text,'\n                        </div>\n                        <div class="chat-content-msg-date">1/22/3333</div>\n                    </div>'):document.getElementById("chat-content").innerHTML+=' <div class="chat-content-msg">\n                        <div class="chat-content-msg-text">\n                            '.concat(t.Msg.Text,'\n                        </div>\n                        <div class="chat-content-msg-date">1/22/3333</div>\n                    </div>');break;case n.Connect:s=t.Uid,o=t.ChatId}})),t.addEventListener("close",(function(e){console.log("Close.")})),t.addEventListener("error",(function(e){console.error("Error:",e)})),document.getElementById("send").onclick=function(){var c={Type:n.TextMsg,Uid:s,ChatId:o,Msg:{Text:e.value}};t.send(JSON.stringify(c)),e.value=""}}},779:(e,t)=>{Object.defineProperty(t,"__esModule",{value:!0}),t.runIfExist=void 0,t.runIfExist=function(e,t){e&&t(e)}}},t={};function n(s){var o=t[s];if(void 0!==o)return o.exports;var c=t[s]={exports:{}};return e[s](c,c.exports,n),c.exports}n.r=e=>{"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},(()=>{n(391);var e=n(779),t=n(872),s=n(862);(0,e.runIfExist)(document.getElementById("header-user"),(function(e){e.onclick=function(){document.getElementById("hu-menu").classList.toggle("hu-menu-hidden")}})),(0,e.runIfExist)(document.getElementById("pc-menu-filter-value"),(function(e){e.onclick=function(){document.getElementById("filter-items").classList.toggle("filter-items-hidden")}})),(0,e.runIfExist)(document.getElementById("pp-del-avatar-label"),(function(e){e.onclick=function(){document.getElementById("pp-del-avatar").classList.toggle("page-panel-checkbox-true")}}));var o=new t.Ajax("/subscribe-post","subscribe-post");o.onSuccess((function(e){var t=document.getElementById("profile-button-subscribe"),n=document.getElementById("subscribers-value"),s=parseInt(n.innerHTML);t.classList.toggle("profile-button-unsubscribe"),t.classList.contains("profile-button-unsubscribe")?(s++,t.innerHTML='<a><img src="/static/img/subscribe.svg">Unsubscribe</a>'):(s>0&&s--,t.innerHTML='<a><img src="/static/img/subscribe.svg">Subscribe</a>'),n.innerHTML=String(s)})),o.listen(),/^\/chat\/\d+$/.test(window.location.pathname)&&(0,s.RunWs)()})()})();