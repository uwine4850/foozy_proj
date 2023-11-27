let oldEntries: Element[] = [];

export class Observer{
    private className: string;
    private options: IntersectionObserverInit;
    constructor(className: string, options: IntersectionObserverInit) {
        this.className = className;
        this.options = options;
    }

    public run(onTrigger: (entry: IntersectionObserverEntry) => void){
        let targets = document.getElementsByClassName(this.className);
        const observer = new IntersectionObserver((entries: IntersectionObserverEntry[], observer: IntersectionObserver) => {
            entries.forEach(entry => {
                if (entry.isIntersecting && !oldEntries.includes(entry.target)) {
                    onTrigger(entry);
                    entry.target.classList.remove(this.className);
                    observer.unobserve(entry.target);
                    oldEntries.push(entry.target);
                }
            });
        }, this.options);
        for (let i = 0; i < targets.length; i++) {
            observer.observe(targets.item(i));
        }
    }
}