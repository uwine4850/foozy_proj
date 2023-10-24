package utils

import (
	"log"
	"net/url"
)

func ConvertApplicationFormFields(fieldsName []string, applicationForm url.Values) (map[string]string, bool) {
	output := map[string]string{}
	for i := 0; i < len(fieldsName); i++ {
		if !applicationForm.Has(fieldsName[i]) {
			log.Println(fieldsName[i])
			return nil, false
		}
		output[fieldsName[i]] = applicationForm.Get(fieldsName[i])
	}
	return output, true
}
