package config

import (
	"os"
	"encoding/json"
)

type Config struct {
	Server Server `json:"server"`
	Database Database `json:"database"`
	FTP FTP `json:"ftp"`
}

type Server struct {
	Addr string `json:"addr"`
}

type Database struct {
	Path string `json:"path"`
}

type FTP struct {
	RootPath string `json:"root"`
}

var configPath = "config/local.json" // времененное решение, надо env делать

func Load() (*Config, error) {
	var cfg Config
	f, err := os.ReadFile(configPath)
	if err != nil { return nil, err }

	if err := json.Unmarshal(f, &cfg); err != nil { return nil, err }

	return &cfg, nil
}
