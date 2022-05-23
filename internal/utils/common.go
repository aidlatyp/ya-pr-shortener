package utils

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

func GetConfigs() (*Config, error) {
	config := &Config{}

	config.ServerAddress = os.Getenv("SERVER_ADDRESS")
	if len(config.ServerAddress) == 0 {
		flag.StringVar(&config.ServerAddress, "a", "http://localhost:8080", "server address")
	}
	config.BaseURL = os.Getenv("BASE_URL")
	if len(config.BaseURL) == 0 {
		flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "base url")
	}
	config.FileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	if len(config.FileStoragePath) == 0 {
		flag.StringVar(&config.FileStoragePath, "f", "", "file storage path")
	}

	flag.Parse()
	return config, nil
}
