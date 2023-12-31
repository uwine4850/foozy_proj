// <div class="pop-up-panel-bg">
//     <div class="pop-up-panel">
//         <button class="pop-up-activate">Activate</button>
//         <div class="pop-up-content pop-up-hide">
//             Text
//         </div>
//     </div>
//     <div class="pop-up-bg pop-up-bg-hide"></div>
// </div>

export class PopUp {
    private readonly className: string;
    private readonly isBackground: boolean;
    private _onClick: (popUp: HTMLElement, popUpActivate: HTMLElement) => void;
    constructor(className: string, isBackground: boolean) {
        this.className = className;
        this.isBackground = isBackground;
    }

    public onClick(click: (popUp: HTMLElement, popUpActivate: HTMLElement) => void){
        this._onClick = click;
    }

    private processingBg(pop_up_bg: HTMLElement){
        pop_up_bg.onclick = () => {
            let parent = pop_up_bg.parentElement;
            let panel = parent.getElementsByClassName(this.className)[0] as HTMLElement;
            let button = panel.getElementsByClassName("pop-up-activate")[0] as HTMLButtonElement;
            button.click();
        }
    }

    public start(){
        let popUps = Array.from(document.getElementsByClassName(this.className)) as HTMLElement[];
        for (let i = 0; i < popUps.length; i++) {
            let pop_up_bg: HTMLElement;
            if (this.isBackground){
                let pop_up_panel_bg = popUps[i].parentElement;
                pop_up_bg = pop_up_panel_bg.getElementsByClassName("pop-up-bg")[0] as HTMLElement;
            }
            let popUpActivate = popUps[i].getElementsByClassName("pop-up-activate");
            popUpActivate[0].addEventListener("click",  () => {
                let s: HTMLElement = popUps[i].getElementsByClassName("pop-up-content")[0] as HTMLElement
                if (pop_up_bg){
                    this.processingBg(pop_up_bg);
                    pop_up_bg.classList.toggle("pop-up-bg-hide");
                }
                this._onClick(s, popUpActivate[0] as HTMLElement);
            })
        }
    }
}

