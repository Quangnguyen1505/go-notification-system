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
		File_name   string `mapstructure:"file_name"`
		Max_size    int    `mapstructure:"max_size"`
		Max_backups int    `mapstructure:"max_backups"`
		Max_age     int    `mapstructure:"max_age"`
		Compress    bool   `mapstructure:"compress"`
		Loglevel    string `mapstructure:"loglevel"`
	}
)
