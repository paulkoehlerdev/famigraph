package config

import (
	_ "embed"
	"encoding/json"
	"github.com/samber/do"
)

//go:embed config.json
var jsonStr []byte

type Config struct {
	Logger struct {
		Level string `json:"level"`
	} `json:"logger"`
	Server struct {
		Addr   string `json:"addr"`
		Domain string `json:"domain"`
		TLS    struct {
			Enabled bool    `json:"enabled"`
			Crt     *string `json:"crt"`
			Key     *string `json:"key"`
		} `json:"tls"`
		ShutdownTimeout string `json:"shutdownTimeout"`
	} `json:"server"`
	Database struct {
		Path string `json:"path"`
	} `json:"database"`
}

func LoadConfig(_ *do.Injector) (Config, error) {
	var config Config
	if err := json.Unmarshal(jsonStr, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
