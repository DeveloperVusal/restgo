package lang

import (
	"apibgo/internal/lang/sections"
	"apibgo/pkg/univenv"
	"log"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

func Locale() string {
	lang := os.Getenv("APP_LANG")

	return lang
}

func Get(lang string) (Lang, bool) {
	appLang := mustLoad()
	value, ok := appLang[lang]

	return value, ok
}

type Lang struct {
	Mail sections.Mail `yaml:"mail"`
}

func mustLoad() map[string]Lang {
	langPath := os.Getenv("LANG_PATH")

	if langPath == "" {
		log.Fatal("LANG_PATH is not set")
	}

	// check if files exists
	files, err := os.ReadDir(langPath)

	if err != nil {
		log.Fatalf("There is a problem in the folder: %s", langPath)
	}

	langs := make(map[string]Lang)

	for _, file := range files {
		withEnv, err := univenv.YamlWithEnv(langPath + file.Name())

		if err != nil {
			log.Fatalf("cannot set env variables to yaml a file: %s", err)
		}

		var langData Lang

		if err := cleanenv.ParseYAML(withEnv, &langData); err != nil {
			log.Fatalf(file.Name()+" -> cannot read language file: %s", err)
		}

		fileName := strings.Split(file.Name(), ".yaml")
		lang := strings.ToLower(fileName[0])

		langs[lang] = langData
	}

	return langs
}
