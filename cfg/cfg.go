package cfg

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Url string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	config := &Config{
		Url: os.Getenv("URL"),
	}

	return config, nil
}
