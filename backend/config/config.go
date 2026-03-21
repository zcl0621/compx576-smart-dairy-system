package config

import (
	"errors"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "config.yaml"

type PGConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type AppConfig struct {
	Env string `yaml:"env"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type ResendConfig struct {
	APIKey string `yaml:"api_key"`
	From   string `yaml:"from"`
}

type Config struct {
	App    AppConfig    `yaml:"app"`
	PG     PGConfig     `yaml:"pg"`
	Redis  RedisConfig  `yaml:"redis"`
	Resend ResendConfig `yaml:"resend"`
	JWT    JWTConfig    `yaml:"jwt"`
}

var (
	appConfig Config
	once      sync.Once
	initErr   error
)

func InitConfig() error {
	once.Do(func() {
		content, err := os.ReadFile(defaultConfigPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				initErr = errors.New("config.yaml not found")
				return
			}

			initErr = err
			return
		}

		if err := yaml.Unmarshal(content, &appConfig); err != nil {
			initErr = err
			return
		}

	})

	return initErr
}

func Get() Config {
	return appConfig
}
