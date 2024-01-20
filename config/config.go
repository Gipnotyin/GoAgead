package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Host string `yaml:"HOST" env-default:"127.0.0.1"`
	Port int32  `yaml:"PORT" env-default:"8081"`
}

func MustConfig() (*Config, error) {
	path := fetchConfigPath()

	if path == "" {
		panic("Path is empty")
	}

	if _, ok := os.Stat(path); os.IsNotExist(ok) {
		return nil, ok
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		res = "config/config.yaml"
	}

	return res
}
