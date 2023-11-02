declare const $: any;

export class Ajax{
    public path: string;
    public formId: string;

    private onSuccessFn: (response: string) => void;
    private onErrorFn: (response: string) => void;

    constructor(path: string, formId: string) {
        this.path = path;
        this.formId = formId;
    }

    onSuccess(fn: (error: string)=> void){
        this.onSuccessFn = fn;
    }

    onError(fn: (error: string)=> void){
        this.onErrorFn = fn;
    }

    listen(){
        $(document).ready(() => {
            $('#' + this.formId).submit((e) => {
                e.preventDefault();
                let formData = $(this).serialize();
                $.ajax({
                    type: 'POST',
                    url: this.path,
                    data: formData,
                    success: (response) => {
                        this.onSuccessFn(response)
                    },
                    error: (xhr, status, error) => {
                        this.onErrorFn(error)
                    }
                });
            });
        });
    }
}