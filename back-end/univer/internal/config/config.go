package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		DSN string `yaml:"dsn"` // DSN (data source name) для БД
	} `yaml:"database"`
	Ids struct {
		Teacher rune `yaml:"teacher"`
		Student rune `yaml:"student"`
		My      rune `yaml:"my"`
	} `yaml:"ids"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

// LoadConfig загружает конфигурацию из YAML-файла
func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
