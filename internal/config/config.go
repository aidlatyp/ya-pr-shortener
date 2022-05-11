package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
)

const ShortenedURLLen int = 5

type ParsedFlags struct {
	addr     *string
	baseURL  *string
	fileName *string
}

func (p *ParsedFlags) Addr() string {
	return *p.addr
}

func (p *ParsedFlags) BaseURL() string {
	return *p.baseURL
}

func (p *ParsedFlags) Filename() string {
	return *p.fileName
}

func NewParsedFlags() ParsedFlags {
	parsed := ParsedFlags{}
	parsed.addr = pflag.StringP("a", "a", "", "Host IP address")
	parsed.baseURL = pflag.StringP("b", "b", "", "Base URL")
	parsed.fileName = pflag.StringP("f", "f", "", "Filename to store URLs")
	pflag.Parse()
	return parsed
}

type ServerConf struct {
	ServerTimeout int64  `env:"SERVER_TIMEOUT"`
	ServerAddr    string `env:"SERVER_ADDRESS"`
}

func NewServerConf(serverFlags ParsedFlags) ServerConf {
	serverConf := ServerConf{
		ServerTimeout: 30,
		ServerAddr:    ":8080",
	}
	if serverFlags.Addr() != "" {
		serverConf.ServerAddr = serverFlags.Addr()
	} else {
		_ = env.Parse(&serverConf)
	}
	return serverConf
}

type App struct {
	BaseURL  string `env:"BASE_URL"`
	FilePath string `env:"FILE_STORAGE_PATH"`
}

func NewAppConf(appFlags ParsedFlags, serverConf ServerConf) App {

	appConf := App{
		BaseURL:  "http://localhost" + serverConf.ServerAddr,
		FilePath: "",
	}

	_ = env.Parse(&appConf)

	if appFlags.BaseURL() != "" {
		appConf.BaseURL = appFlags.BaseURL()
	}

	if appFlags.Filename() != "" {
		appConf.FilePath = appFlags.Filename()
	}

	return appConf
}
