package univenv

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

func Load() {
	paths := getPaths()
	godotenv.Load(paths...)
}

func getPaths() []string {
	dir, _ := os.Getwd()
	projectPath := filepath.Join(filepath.Dir(dir), filepath.Base(dir))
	result := []string{}

	reverseWalkFile(&result, projectPath, "env")

	return result
}

func reverseWalkFile(res *[]string, path string, ext string) {
	files, err := os.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fileExt := filepath.Ext(file.Name())

			if fileExt == ".env" {
				filename := filepath.Join(path, file.Name())

				*res = append(*res, filename)
			}
		}
	}

	for _, file := range files {
		if file.IsDir() && file.Name() == "internal" {
			return
		}
	}

	reverseWalkFile(res, filepath.Join(path, "/../"), ext)
}

func YamlWithEnv(path string) (io.Reader, error) {
	yfile, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	text := string(yfile[:])

	re := regexp.MustCompile(`\${.*}`)
	matches := re.FindAllString(text, -1)

	for _, match := range matches {
		env := match[2:]
		env = env[:len(env)-1]
		text = strings.ReplaceAll(text, match, os.Getenv(env))
	}

	r := bytes.NewReader([]byte(text))

	return r, nil
}
