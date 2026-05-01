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
		GRPC         `yaml:"grpc"`
		configs.Log  `yaml:"logger"`
	}

	GRPC struct {
		NotificationHost string `yaml:"notification_host" env:"GRPC_NOTIFICATION_HOST"`
		NotificationPort int    `yaml:"notification_port" env:"GRPC_NOTIFICATION_PORT"`
		UserHost         string `yaml:"user_host" env:"GRPC_USER_HOST"`
		UserPort         int    `yaml:"user_port" env:"GRPC_USER_PORT"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println("dir:", dir)

	configCandidates := []string{
		filepath.Join(dir, "config.yml"),
		filepath.Join(dir, "cmd", "proxy", "config.yml"),
	}

	var configPath string
	for _, candidate := range configCandidates {
		if _, statErr := os.Stat(candidate); statErr == nil {
			configPath = candidate
			break
		}
	}

	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	// Optional env overrides
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("env error: %w", err)
	}

	if configPath == "" {
		return nil, fmt.Errorf("config error: could not find config.yml (tried: %v)", configCandidates)
	}

	return cfg, nil
}
