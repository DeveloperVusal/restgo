package mails

import (
	"apibgo/internal/lang"
	"log"
	"regexp"
)

func Registration(replace map[string]string) (string, string) {
	appLang, _ := lang.Get(lang.Locale())

	subject := appLang.Mail.Registration.Subject
	body := appLang.Mail.Registration.Body

	for key, value := range replace {
		re, err := regexp.Compile(`{{ *` + key + ` *}}`)

		if err != nil {
			log.Fatal(err)
		}

		subject = string(re.ReplaceAll([]byte(subject), []byte(value)))
		body = string(re.ReplaceAll([]byte(body), []byte(value)))
	}

	return subject, body
}
