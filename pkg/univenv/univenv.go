package univenv

import (
	"log"
	"os"
	"path/filepath"

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
