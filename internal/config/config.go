package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env string `json:"env" yaml:"env" mapstructure:"env" env-default:"local"`
	//StoragePath string        `json:"storagePath" yaml:"storagePath" mapstructure:"storagePath" env-required:"true"`
	TokenTTL time.Duration `json:"token_ttl" yaml:"token_ttl" mapstructure:"token_ttl" env-required:"true"`
	Server   ServerConfig  `json:"server" yaml:"server" mapstructure:"server"`
	//Database    Database      `json:"database" yaml:"database" mapstructure:"database"`
}

func NewConfig(env string) *Config {
	return &Config{Env: env}
}

type ServerConfig struct {
	Port    int           `json:"port" yaml:"port" mapstructure:"port"`
	Timeout time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout" env-required:"true"`
}

/*
	type Database struct {
		Host     string `json:"host" yaml:"host" mapstructure:"host" env-required:"true"`
		Port     string `json:"port" yaml:"port" mapstructure:"port" env-required:"true"`
		User     string `json:"user" yaml:"user" mapstructure:"user" env-required:"true"`
		Password string `json:"password" yaml:"password" mapstructure:"password" env-required:"true"`
		Name     string `json:"name" yaml:"name" mapstructure:"name" env-required:"true"`
	}
*/
func MustLoadConfig() Config {
	path := fetchConfigPath(os.Getenv("CONFIG_PATH"))
	if path == "" {
		panic("Config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file does not exist " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Error reading config file " + path + " " + err.Error())
	}

	// TODO: there should be a better way to do this, before was &cfg.
	return cfg

}

// fetchConfigPath fetches the config path from the command line arguments or environment variables.
// Priority: command line args > env vars > default.
// Default value is empty string.
func fetchConfigPath(path string) string {
	var res string
	// --config=<path/to/config,yaml>
	flag.StringVar(&res, "config", "", `Path to config file`)
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
