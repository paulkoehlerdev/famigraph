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
		Addr            string `json:"addr"`
		ShutdownTimeout string `json:"shutdownTimeout"`
	} `json:"server"`
}

func LoadConfig(_ *do.Injector) (Config, error) {
	var config Config
	if err := json.Unmarshal(jsonStr, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
