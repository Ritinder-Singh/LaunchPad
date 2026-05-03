package config

import (
	"encoding/json"
	"os"
)

type Links struct {
	GitHub string `json:"github"`
	Live   string `json:"live,omitempty"`
}

type Media struct {
	Images []string `json:"images"`
	Video  string   `json:"video,omitempty"`
}

type Deploy struct {
	Subdomain string `json:"subdomain"`
	Domain    string `json:"domain"`
	Provider  string `json:"provider"`
}

type Config struct {
	Name        string   `json:"name"`
	Tagline     string   `json:"tagline"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
	QuickStart  string   `json:"quickStart"`
	TechStack   []string `json:"techStack"`
	Media       Media    `json:"media"`
	Links       Links    `json:"links"`
	Template    string   `json:"template"`
	Deploy      Deploy   `json:"deploy"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
