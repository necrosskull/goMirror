package cfg

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	Url  string
	Port string
}

var Settings = &Config{}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables only")
	}

	config := &Config{
		Url:  os.Getenv("URL"),
		Port: os.Getenv("PORT"),
	}

	configMap := make(map[string]string)

	el := reflect.ValueOf(config).Elem()
	for i := 0; i < el.NumField(); i++ {
		field := el.Type().Field(i).Name
		value := el.Field(i).Interface().(string)
		configMap[field] = value
	}

	for key, value := range configMap {
		if value == "" {
			return nil, fmt.Errorf("error loading config: %s is not set", key)
		}
	}

	return config, nil
}

func Init() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	Settings.Port = config.Port
	Settings.Url = config.Url

}
