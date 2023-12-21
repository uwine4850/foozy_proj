const axios = require('axios').default;
declare const $: any;

interface FormDataFormValue {
    name: string,
    value: any,
    fileName?: string
}

interface UrlFormValue{
    name: string,
    value: string,
}

export class Ajax{
    public path: string;
    public formId: string;

    private onSuccessFn: (response: object) => void;
    private onErrorFn: (response: string) => void;
    private formValues: FormDataFormValue[];
    private urlFormValue: UrlFormValue[];

    constructor(path: string, formId: string) {
        this.path = path;
        this.formId = formId;
        this.formValues = [];
    }

    onSuccess(fn: (error: object)=> void){
        this.onSuccessFn = fn;
    }

    onError(fn: (error: string)=> void){
        this.onErrorFn = fn;
    }

    public setUrlFormValue(name: string, value: string){
        let urlFormValue: UrlFormValue = {
            name: name,
            value: value,
        }
        this.urlFormValue.push(urlFormValue);
    }

    public setMultipartFormValue(name: string, value: string){
        let formValue: FormDataFormValue = {
            name: name,
            value: value,
        }
        this.formValues.push(formValue);
    }

    public setMultipartFormFile(name: string, value: Blob, fileName: string){
        let formValue: FormDataFormValue = {
            name: name,
            value: value,
            fileName: fileName
        }
        this.formValues.push(formValue);
    }

    async listen(){
        const form = document.getElementById(this.formId) as HTMLFormElement;
        const formEnctype: string = form.getAttribute('enctype');

        form.addEventListener('submit', (e) =>{
            e.preventDefault();
            let data;
            let formData = new FormData(form);
            for (const formValue of this.formValues) {
                if (!formValue.fileName){
                    formData.append(formValue.name, formValue.value);
                } else {
                    formData.append(formValue.name, formValue.value, formValue.fileName);
                }
            }
            if (formEnctype == "application/x-www-form-urlencoded"){
                let u = new URLSearchParams(formData as any);
                for (const urlFormValue of this.urlFormValue) {
                    u.set(urlFormValue.name, urlFormValue.value);
                }
                data = u.toString()
            } else {
                data = formData;
            }
            this.send(data, formEnctype);
        });
    }

    private async send(formData: any, enctype: string){
        try {
            const response = await axios.post(this.path, formData, {
                headers: {
                    'Content-Type': enctype,
                },
            });
            this.onSuccessFn(response.data);
        } catch (error) {
            this.onErrorFn(error.data);
        }
    }
}

export function ajaxGET(url: string, data: object, onSuccess: (response) => void, onError: (error) => void){
    $.ajax({
        url: url,
        method: 'GET',
        data: data,
        json: true,
        success: function(response) {
            onSuccess(response);
        },
        error: function(xhr, status, error) {
            onError(error);
        }
    });
}
