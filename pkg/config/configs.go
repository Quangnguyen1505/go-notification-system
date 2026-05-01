package configs

type (
	App struct {
		Name    string `yaml:"name" env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Host string `yaml:"host" env:"HTTP_HOST"`
		Port int    `yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `yaml:"log_level" env:"LOG_LEVEL"`
	}
)
