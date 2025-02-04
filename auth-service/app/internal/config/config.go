package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Logging Logging `yaml:"logging"`
	Server  Server  `yaml:"server"`
}

type Logging struct {
	Level string `yaml:"level"`
}

type Server struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port int    `yaml:"port" env-default:"59999"`
	TLS  TLS    `yaml:"tls"`
}

type TLS struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

func MustLoadConfig() *Config {
	path := getConfigPath()

	if _, err := os.Stat(path); err != nil {
		log.Fatalf("Config file is not exists: config-path=%s error=%s", path, err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("Failed to read config: config-path=%s error=%v", path, err)
	}

	return &cfg
}

func getConfigPath() string {
	configPath := flag.String("config-path", "", "config path")

	flag.Parse()

	if *configPath == "" {
		*configPath = os.Getenv("CONFIG_PATH")
	}

	if *configPath == "" {
		log.Fatal("config path is empty")
	}
	return *configPath
}
