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
	sync.Once
}

// FlagGetter abstracts from flag source
type FlagGetter interface {
	BaseURL() string
	Filename() string
	Addr() string
}

// App do not support hot configuration reload
// so in this case usually makes sense explicitly made configuration happen only once
// Configuration Singleton (?)
var appConfig *AppConfig = nil

func NewAppConfig(appFlags FlagGetter) *AppConfig {
	if appConfig == nil {
		appConfig = &AppConfig{}
		appConfig.Do(func() {
			appConfig.configure(appFlags)
		})
		return appConfig
	}
	return appConfig
}

func (a *AppConfig) IsFilePathSet() bool {
	return a.FilePath != ""
}

func (a *AppConfig) configure(appFlags FlagGetter) {

	// Default configuration - if it will not be overwritten below
	// Low priority
	a.BaseURL = "http://localhost" + ":8080"
	a.FilePath = ""
	a.ServerTimeout = 30
	a.ServerAddr = ":8080"

	// Configure with ENV vars
	// Middle priority
	err := env.Parse(a)
	if err != nil {
		log.Printf("error while parsing application env vars, %v", err)
	}

	// Configure with Flags
	// High priority
	if appFlags.Addr() != "" {
		a.ServerAddr = appFlags.Addr()
	}
	if appFlags.BaseURL() != "" {
		a.BaseURL = appFlags.BaseURL()
	}
	if appFlags.Filename() != "" {
		a.FilePath = appFlags.Filename()
	}

	a.BaseURL += "/"
}
