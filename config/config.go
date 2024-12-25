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
		CookiePrefix string `json:"cookiePrefix"`
		JWT          struct {
			Secret string `json:"secret"`
		} `json:"jwt"`
	} `json:"session"`
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
