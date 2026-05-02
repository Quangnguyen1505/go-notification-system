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
		File_name   string `yaml:"file_name" mapstructure:"file_name"`
		Max_size    int    `yaml:"max_size" mapstructure:"max_size"`
		Max_backups int    `yaml:"max_backups" mapstructure:"max_backups"`
		Max_age     int    `yaml:"max_age" mapstructure:"max_age"`
		Compress    bool   `yaml:"compress" mapstructure:"compress"`
		Loglevel    string `yaml:"loglevel" mapstructure:"loglevel"`
	}

	Postgres struct {
		Host            string `yaml:"host" mapstructure:"host"`
		Username        string `yaml:"username" mapstructure:"username"`
		Password        string `yaml:"password" mapstructure:"password"`
		Port            int    `yaml:"port" mapstructure:"port"`
		Dbname          string `yaml:"dbname" mapstructure:"dbname"`
		MaxOpenConns    int    `yaml:"maxOpenConns" mapstructure:"maxOpenConns"`
		MinConns        int    `yaml:"minConns" mapstructure:"minConns"`
		ConnMaxLifeTime int    `yaml:"connMaxLifeTime" mapstructure:"connMaxLifeTime"`
	}
)
