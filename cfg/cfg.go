package cfg

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Url  string
	Port string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	config := &Config{
		Url:  os.Getenv("URL"),
		Port: os.Getenv("PORT"),
	}

	return config, nil
}
