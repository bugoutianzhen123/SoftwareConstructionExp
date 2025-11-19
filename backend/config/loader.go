package config

import (
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

type AppConfig struct {
    Database string `yaml:"database"`
}

func Load() (*AppConfig, error) {
    p := filepath.Join("config", "config.yaml")
    b, err := os.ReadFile(p)
    if err != nil { return nil, err }
    var c AppConfig
    if err := yaml.Unmarshal(b, &c); err != nil { return nil, err }
    return &c, nil
}