package validators

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

func xssValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		if containsHtml(field.String()) {
			return false
		}
		return true
	}
	return false
}

func containsHtml(value string) bool {
	htmlPattern := "<(\\S*?)[^>]*>.*?|<.*? />"
	pattern := regexp.MustCompile(htmlPattern)
	matches := pattern.FindStringSubmatch(value)
	return len(matches) > 0
}
