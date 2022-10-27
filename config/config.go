package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
)

type QueryType string

const (
	QueryTypeOne  QueryType = "one"
	QueryTypeMany QueryType = "many"
)

type Config struct {
	DSN     string  `json:"dsn"`
	Queries []Query `json:"queries"`
}
type Query struct {
	Type QueryType `json:"type"`
	Path string    `json:"path"`
	SQL  string    `json:"sql"`
	Argc int       `json:"argc"`
}

func (cfg *Config) Driver() (string, error) {
	u, err := url.Parse(cfg.DSN)
	if err != nil {
		return "", fmt.Errorf("failed to parse DSN: %w", err)
	}
	return u.Scheme, nil
}

var Getenv = os.Getenv

func Parse(data []byte) (Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse: %w", err)
	}
	if cfg.DSN == "" {
		cfg.DSN = Getenv("SQL_PROXY_DSN")
		if cfg.DSN == "" {
			return Config{}, errors.New("missing DSN: JSON config \"dsn\" or environment variable \"SQL_PROXY_DSN\" must be passed")
		}
	}
	return cfg, nil
}
func ParseFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open file: %w", err)
	}
	return Parse(data)
}
