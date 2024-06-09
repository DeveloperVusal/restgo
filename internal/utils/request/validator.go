package request

import (
	"reflect"
	"regexp"
	"strings"

	"apibgo/internal/lang"
	"apibgo/internal/lang/sections"
	"apibgo/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Messages sections.Validation
}

func NewValidator() Validator {
	appLang, _ := lang.Get(lang.Locale())

	return Validator{
		Messages: appLang.Validation,
	}
}

func (v *Validator) Validate(d interface{}) (bool, []string) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(d)

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return false, []string{err.Error()}
		}

		messages := []string{}

		for _, err := range err.(validator.ValidationErrors) {
			values := map[string]string{}

			if len(err.Param()) > 0 {
				values = v.getAppends(err.Tag(), err.Param())
			}

			message := v.getMessage(err.Tag(), []string{err.Field()}, values)
			messages = append(messages, message)
		}

		// from here you can create your own error messages in whatever language you wish
		return false, messages
	}

	return true, []string{}
}

func (v *Validator) getAppends(tag string, param string) map[string]string {
	var keys map[string]string = map[string]string{}
	var field reflect.Value

	vl := reflect.ValueOf(v.Messages)
	tp := reflect.TypeOf(v.Messages)

	for i := 0; i < tp.NumField(); i++ {
		if tp.Field(i).Tag.Get("yaml") == tag {
			field = vl.Field(i)
			break
		}
	}

	if field.IsValid() {
		re, _ := regexp.Compile(`:\w+`)
		matches := re.FindAllString(field.String(), -1)
		params := strings.Split(param, " ")

		if utils.StringInSlice(":other", matches) &&
			utils.StringInSlice(":value", matches) {
			keys["other"] = params[0]
			keys["value"] = params[1]
		} else if utils.StringInSlice(":other", matches) {
			keys["other"] = strings.Join(params, " ")
		} else if utils.StringInSlice(":values", matches) {
			keys["values"] = strings.Join(params, " ")
		} else if utils.StringInSlice(":value", matches) {
			keys["value"] = strings.Join(params, " ")
		}
	}

	return keys
}

func (v *Validator) getMessage(tag string, attributes []string, values map[string]string) string {
	var msg string
	var field reflect.Value

	vl := reflect.ValueOf(v.Messages)
	tp := reflect.TypeOf(v.Messages)

	for i := 0; i < tp.NumField(); i++ {
		if tp.Field(i).Tag.Get("yaml") == tag {
			field = vl.Field(i)
			break
		}
	}

	if field.IsValid() {
		msg = field.String()

		for _, attr := range attributes {
			re, _ := regexp.Compile(`:attribute`)
			msg = string(re.ReplaceAll([]byte(msg), []byte(attr)))
		}

		for key, val := range values {
			re, _ := regexp.Compile(`:` + key)
			msg = string(re.ReplaceAll([]byte(msg), []byte(val)))
		}

		// for key, val := range custom {
		// 	re, _ := regexp.Compile(`:` + key)
		// 	msg = string(re.ReplaceAll([]byte(msg), []byte(val)))
		// }
	}

	return msg
}
