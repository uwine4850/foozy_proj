package utils

import (
	"github.com/uwine4850/foozy/pkg/router/form"
	"net/http"
	"net/url"
)

func ConvertApplicationFormFields(fieldsName []string, applicationForm url.Values) (map[string]string, bool) {
	output := map[string]string{}
	for i := 0; i < len(fieldsName); i++ {
		if !applicationForm.Has(fieldsName[i]) {
			return nil, false
		}
		output[fieldsName[i]] = applicationForm.Get(fieldsName[i])
	}
	return output, true
}

func ServerError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err))
}

func ParseForm(r *http.Request) (*form.Form, error) {
	frm := form.NewForm(r)
	err := frm.Parse()
	if err != nil {
		return nil, err
	}
	err = frm.ValidateCsrfToken()
	if err != nil {
		return nil, err
	}
	return frm, nil
}

func RemoveElement[T comparable](slice []T, element T) []T {
	var result []T
	for _, el := range slice {
		if el != element {
			result = append(result, el)
		}
	}
	return result
}
