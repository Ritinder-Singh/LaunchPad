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

type Build struct {
	Command      string `json:"command"`      // e.g. "npm run build"
	StartCommand string `json:"startCommand"` // e.g. "npm start"
	Port         int    `json:"port"`         // port the app listens on for screenshot capture
	OutputDir    string `json:"outputDir"`    // static output dir (for static sites)
}

type Screenshots struct {
	Dir         string `json:"dir"`         // where to look/store screenshots; default ".showcase/screenshots"
	AutoCapture bool   `json:"autoCapture"` // run Puppeteer if no images found; default true
}

type Config struct {
	Name        string      `json:"name"`
	Tagline     string      `json:"tagline"`
	Description string      `json:"description"`
	Features    []string    `json:"features"`
	QuickStart  string      `json:"quickStart"`
	TechStack   []string    `json:"techStack"`
	Media       Media       `json:"media"`
	Links       Links       `json:"links"`
	Template    string      `json:"template"`
	Deploy      Deploy      `json:"deploy"`
	Build       Build       `json:"build"`
	Screenshots Screenshots `json:"screenshots"`
}

// DefaultConfigPath is the preferred config location.
const DefaultConfigPath = ".showcase/config.json"

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// Apply defaults
	if cfg.Screenshots.Dir == "" {
		cfg.Screenshots.Dir = ".showcase/screenshots"
	}
	if !cfg.Screenshots.AutoCapture {
		cfg.Screenshots.AutoCapture = true
	}
	return &cfg, nil
}
