import {ajaxGET} from "./ajax";

export class LazyLoad{
    private obsElementId: string;
    private dataFields: string[];
    private options: IntersectionObserverInit;
    private url: string;
    private valuesData: Record<string, string> = {};
    constructor(elementId: string, dataFields: string[], url: string) {
        this.obsElementId = elementId;
        this.dataFields = dataFields;
        this.url = url;
    }

    public setOptions(options: IntersectionObserverInit){
        this.options = options;
    }

    private getValues(obsElement: Element){
        // let data: Record<string, string> = {};
        for (const dataField of this.dataFields) {
            let f: string | null = obsElement.getAttribute("data-" + dataField);
            if (f != null){
                this.valuesData[dataField] = f;
            }
        }
        return this.valuesData;
    }

    public setValue(key: string, value: string){
        this.valuesData[key] = value;
    }

    public run(onTrigger: (response) => void, onError: (error: string) => void){
        // let obsElement = document.getElementById(this.obsElementId);
        let obsElement = document.getElementsByClassName(this.obsElementId);
        const observer = new IntersectionObserver((entries: IntersectionObserverEntry[], observer: IntersectionObserver) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    observer.unobserve(entry.target);
                    ajaxGET(this.url, this.getValues(entry.target), function (response){
                        onTrigger(response);
                    }, function (error){
                        onError(error);
                    })
                    entry.target.classList.remove(this.obsElementId)
                    setTimeout(() => {
                        let obsElement = document.getElementsByClassName(this.obsElementId);
                        for (let i = 0; i < obsElement.length; i++) {
                            observer.observe(obsElement.item(i));
                        }
                    }, 1000);
                }
            });
        }, this.options);
        for (let i = 0; i < obsElement.length; i++) {
            observer.observe(obsElement.item(i));
        }
    }
}
