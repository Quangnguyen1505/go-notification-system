package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	Postgresql struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"postgresql"`
}

func main() {
	configPath := flag.String("config", "cmd/notification/config.yml", "path to YAML config file")
	sslmode := flag.String("sslmode", "disable", "sslmode for postgres DSN")
	flag.Parse()

	cfg := &config{}
	if err := cleanenv.ReadConfig(*configPath, cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.Postgresql.Username,
		cfg.Postgresql.Password,
		cfg.Postgresql.Dbname,
		cfg.Postgresql.Host,
		cfg.Postgresql.Port,
		*sslmode,
	)
}
