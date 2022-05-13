package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

// ShortenedURLLen configure "main functionality"
const ShortenedURLLen int = 5

// AppConf is application specific configuration. As BASE_URL
// is visible and affects user, decided to put it here not to ServerConf
type AppConf struct {
	BaseURL  string `env:"BASE_URL"`
	FilePath string `env:"FILE_STORAGE_PATH"`
}

func (a *AppConf) IsFilePathSet() bool {
	if a.FilePath != "" {
		return true
	}
	return false
}

type AppsFlagsGetter interface {
	BaseURL() string
	Filename() string
}

func NewAppConf(appFlags AppsFlagsGetter, servAddr string) AppConf {
	// Default configuration - if it will not be overwritten below
	// Low priority
	appConf := AppConf{
		BaseURL:  "http://localhost" + servAddr,
		FilePath: "",
	}

	// Configure with ENV vars
	// Middle priority
	// Straight dep from env package, should be used through interface as flags but env seems more common
	// and less user dependent, so no over-engineering here
	err := env.Parse(&appConf)
	if err != nil {
		log.Printf("error while parsing application env vars, %v", err)
	}

	// Configure with flags
	// High priority
	if appFlags.BaseURL() != "" {
		appConf.BaseURL = appFlags.BaseURL()
	}

	if appFlags.Filename() != "" {
		appConf.FilePath = appFlags.Filename()
	}

	appConf.BaseURL += "/"
	return appConf
}

// ServerConf is responsible for sort of "technical" server configuration
// maybe it will make sense to add BASE_URL here
type ServerConf struct {
	ServerTimeout int64  `env:"SERVER_TIMEOUT"`
	ServerAddr    string `env:"SERVER_ADDRESS"`
	// other server conf and restrictions
}

type ServerFlagsGetter interface {
	Addr() string
}

func NewServerConf(serverFlags ServerFlagsGetter) ServerConf {
	serverConf := ServerConf{
		ServerTimeout: 30,
		ServerAddr:    ":8080",
	}

	// todo to flat "ifelse"
	if serverFlags.Addr() != "" {
		serverConf.ServerAddr = serverFlags.Addr()
	} else {
		err := env.Parse(&serverConf)
		if err != nil {
			log.Printf("error while parsing server env vars, %v", err)
		}
	}
	return serverConf
}
