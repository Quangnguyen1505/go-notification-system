package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
	configs "github.com/quangnguyen1505/go-notification-system/pkg/config"
)

type (
	Config struct {
		configs.App  `yaml:"app"`
		configs.HTTP `yaml:"http"`
		configs.Log  `yaml:"logger"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println(dir)

	configCandidates := []string{
		filepath.Join(dir, "config.yml"),
		filepath.Join(dir, "cmd", "notification", "config.yml"),
	}

	var configPath string
	for _, candidate := range configCandidates {
		if _, statErr := os.Stat(candidate); statErr == nil {
			configPath = candidate
			break
		}
	}

	if configPath == "" {
		return nil, fmt.Errorf("config error: could not find config.yml (tried: %v)", configCandidates)
	}

	err = cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	// Optional env overrides (won't error if env vars are missing)
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("env error: %w", err)
	}

	return cfg, nil
}
