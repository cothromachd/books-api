package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	API API `yaml:"api"`
	DB DB `yaml:"db"`
	RDB RDB `yaml:"rdb"`
}

type API struct {
	Host string `yaml:"host"`
}

type DB struct {
	Conn string `yaml:"conn"`
}

type RDB struct {
	Conn string `yaml:"conn"`
}

func New(path string) (_ *Config, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	return &cfg, yaml.NewDecoder(file).Decode(&cfg)
}