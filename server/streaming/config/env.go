package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func findEnv() string {
	dir, _ := os.Getwd()

	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return filepath.Join(dir, ".env")
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ".env" // fallback
		}
		dir = parent
	}
}

func LoadEnv() error {
	if os.Getenv("ENV") != "CONTAINER" {
		err := godotenv.Load(findEnv())
		if err != nil {
			return err
		}
	}
	return nil
}
