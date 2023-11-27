import {ajaxGET} from "./ajax";

export class LazyLoad{
    private obsElementId: string;
    private dataFields: string[];
    private options: IntersectionObserverInit;
    private url: string;
    constructor(elementId: string, dataFields: string[], url: string) {
        this.obsElementId = elementId;
        this.dataFields = dataFields;
        this.url = url;
    }

    public setOptions(options: IntersectionObserverInit){
        this.options = options;
    }

    private getValues(obsElement: HTMLElement){
        let data: Record<string, string> = {};
        for (const dataField of this.dataFields) {
            let f: string | null = obsElement.getAttribute("data-" + dataField);
            if (f != null){
                data[dataField] = f;
            }
        }
        return data;
    }

    public run(onTrigger: (response) => void, onError: (error: string) => void){
        let obsElement = document.getElementById(this.obsElementId);
        const observer = new IntersectionObserver((entries: IntersectionObserverEntry[], observer: IntersectionObserver) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    observer.unobserve(entry.target);
                    ajaxGET(this.url, this.getValues(obsElement), function (response){
                        onTrigger(response);
                    }, function (error){
                        onError(error);
                    })
                    obsElement.removeAttribute("id");
                    setTimeout(() => {
                        obsElement = document.getElementById(this.obsElementId);
                        if (obsElement) {
                            observer.observe(obsElement);
                        }
                    }, 100);
                }
            });
        }, this.options);
        observer.observe(obsElement);
    }
}
