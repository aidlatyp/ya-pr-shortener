package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v6"
)

// ShortenedURLLen configure "main functionality"
const ShortenedURLLen int = 5

// AppConfig is application specific configuration.
type AppConfig struct {
	BaseURL       string `env:"BASE_URL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	ServerTimeout int64  `env:"SERVER_TIMEOUT"`
	ServerAddr    string `env:"SERVER_ADDRESS"`
}

func (a *AppConfig) IsFilePathSet() bool {
	return a.FilePath != ""
}

// FlagGetter abstracts from flag source
type FlagGetter interface {
	BaseURL() string
	Filename() string
	Addr() string
}

// App do not support hot configuration reload
// so in this case usually makes sense explicitly made configuration happen only once
var once sync.Once

func NewAppConfig(appFlags FlagGetter) AppConfig {
	var appConf AppConfig
	once.Do(func() {

		// Default configuration - if it will not be overwritten below
		// Low priority
		appConf = AppConfig{
			BaseURL:       "http://localhost" + ":8080",
			FilePath:      "",
			ServerTimeout: 30,
			ServerAddr:    ":8080",
		}

		// Configure with ENV vars
		// Middle priority
		err := env.Parse(&appConf)
		if err != nil {
			log.Printf("error while parsing application env vars, %v", err)
		}

		// Configure with Flags
		// High priority
		if appFlags.Addr() != "" {
			appConf.ServerAddr = appFlags.Addr()
		}
		if appFlags.BaseURL() != "" {
			appConf.BaseURL = appFlags.BaseURL()
		}
		if appFlags.Filename() != "" {
			appConf.FilePath = appFlags.Filename()
		}

		appConf.BaseURL += "/"
	})
	return appConf
}
