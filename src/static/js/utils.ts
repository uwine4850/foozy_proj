export function runIfExist<T>(element: T, fn: (el: T) => void){
    if (element){
        fn(element)
    }
}