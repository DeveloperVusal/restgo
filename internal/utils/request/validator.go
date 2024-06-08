package request

import (
	"fmt"
	"log"
	"reflect"
	"regexp"

	"apibgo/internal/lang"
	"apibgo/internal/lang/sections"

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

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return false, []string{err.Error()}
		}

		messages := []string{}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.StructField())

			fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Type())
			// // fmt.Println(err.Value())
			// fmt.Println()

			messages = append(messages, v.getMessage(err.StructField(), []string{err.Tag()}, map[string]string{}))
		}

		// from here you can create your own error messages in whatever language you wish
		return false, messages
	}

	return true, []string{}
}

func (v *Validator) getMessage(field string, attributes []string, values map[string]string) string {
	vl := reflect.ValueOf(v.Messages)
	value := vl.FieldByName(field)
	msg := value.String()

	for _, attr := range attributes {
		re, err := regexp.Compile(`:attribute`)

		if err != nil {
			log.Fatal(err)
		}

		msg = string(re.ReplaceAll([]byte(msg), []byte(attr)))
	}

	for key, val := range values {
		re, err := regexp.Compile(`:` + key)

		if err != nil {
			log.Fatal(err)
		}

		msg = string(re.ReplaceAll([]byte(msg), []byte(val)))
	}

	fmt.Println("msg -->", msg)

	return msg
}
