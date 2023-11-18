package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "go-cafe-crawl"

func LoadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	curWorkDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error get path name")
	}

	rootPath := projectName.Find([]byte(curWorkDir))

	err = godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	log.Println("loading .env file")
}

func Environment(key string) string {
	return os.Getenv(key)
}
