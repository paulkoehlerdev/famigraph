package config

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/samber/do"
	"os"
)

type Config struct {
	Logger struct {
		Level string `json:"level"`
	} `json:"logger"`
	Server struct {
		TLSAddr string `json:"tlsaddr"`
		Addr    string `json:"addr"`
		Domain  string `json:"domain"`
		TLS     struct {
			Crt string `json:"crt"`
			Key string `json:"key"`
		} `json:"tls"`
		ShutdownTimeout string `json:"shutdownTimeout"`
	} `json:"server"`
	Webauthn struct {
		Timeout struct {
			Registration string `json:"registration"`
			Login        string `json:"login"`
		} `json:"timeout"`
		RelyingParty struct {
			ID             string   `json:"id"`
			DisplayName    string   `json:"displayName"`
			AllowedOrigins []string `json:"allowedOrigins"`
		} `json:"relyingParty"`
	} `json:"webauthn"`
	Session struct {
		Expiry       string `json:"expiry"`
		CookiePrefix string `json:"cookiePrefix"`
		JWT          struct {
			Secret string `json:"secret"`
		} `json:"jwt"`
	} `json:"session"`
	Connect struct {
		SigningKey string `json:"signingKey"`
		Expiry     string `json:"expiry"`
	} `json:"connect"`
	Statistics struct {
		UpdateInterval string `json:"updateInterval"`
		UpdateTimeout  string `json:"updateTimeout"`
	} `json:"statistics"`
	Database struct {
		Path string `json:"path"`
	} `json:"database"`
}

func LoadConfig(_ *do.Injector) (Config, error) {
	path := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	if *path == "" {
		return Config{}, fmt.Errorf("path to config file is required")
	}

	file, err := os.ReadFile(*path)
	if err != nil {
		return Config{}, fmt.Errorf("could not read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
