package config

const ShortenedURLLen int = 5

type Server struct {
	ServerTimeout int64  `env:"SERVER_TIMEOUT"`
	ServerAddr    string `env:"SERVER_ADDRESS"`
}

type App struct {
	BaseURL string `env:"BASE_URL"`
}
