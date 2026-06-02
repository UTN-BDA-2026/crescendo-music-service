package env

import (
	"os"

	"github.com/joho/godotenv"
)

func Load() error {
	if os.Getenv("ENV") != "CONTAINER" {
		err := godotenv.Load("../.env")
		if err != nil {
			return err
		}
	}
	return nil
}
