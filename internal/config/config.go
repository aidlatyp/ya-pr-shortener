package config

const ShortenedURLLen int = 5

type Server struct {
	ServerTimeout int64  `env:"SERVER_TIMEOUT"`
	ServerAddr    string `env:"SERVER_ADDRESS"`
	ServerHost    string `env:"SERVER_HOST"`
	ServerPort    string `env:"SERVER_PORT"`
}

type App struct {
	BaseURL  string `env:"BASE_URL"`
	FilePath string `env:"FILE_STORAGE_PATH"`
}
