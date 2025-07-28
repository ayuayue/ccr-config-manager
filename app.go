package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
)

// App struct
type App struct {
	ctx context.Context
}

// Config represents the Claude Code Router configuration
type Config struct {
	APIKEY              string     `json:"APIKEY,omitempty"`
	PROXY_URL          string     `json:"PROXY_URL,omitempty"`
	HOST               string     `json:"HOST,omitempty"`
	API_TIMEOUT_MS     int        `json:"API_TIMEOUT_MS,omitempty"`
	LOG                bool       `json:"LOG,omitempty"`
	Providers          []Provider `json:"Providers,omitempty"`
	Router             Router     `json:"Router,omitempty"`
}

// Provider represents a model provider configuration
type Provider struct {
	Name       string      `json:"name"`
	APIBaseURL string      `json:"api_base_url"`
	APIKey     string      `json:"api_key"`
	Models     []string    `json:"models"`
	Transformer interface{} `json:"transformer,omitempty"`
}

// Router represents the routing configuration
type Router struct {
	Default              string `json:"default,omitempty"`
	Background           string `json:"background,omitempty"`
	Think                string `json:"think,omitempty"`
	LongContext          string `json:"longContext,omitempty"`
	LongContextThreshold int    `json:"longContextThreshold,omitempty"`
	WebSearch            string `json:"webSearch,omitempty"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetConfigPath returns the path to the Claude Code Router config file
func (a *App) GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".claude-code-router", "config.json")
}

// LoadConfig loads the Claude Code Router configuration
func (a *App) LoadConfig() (Config, error) {
	var config Config
	
	configPath := a.GetConfigPath()
	if configPath == "" {
		return config, fmt.Errorf("could not determine config path")
	}
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return empty config if file doesn't exist
		return config, nil
	}
	
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	
	// Parse JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	
	// Set default values if not present
	if config.API_TIMEOUT_MS == 0 {
		config.API_TIMEOUT_MS = 600000
	}
	
	if config.Router.LongContextThreshold == 0 {
		config.Router.LongContextThreshold = 60000
	}
	
	return config, nil
}

// SaveConfig saves the Claude Code Router configuration
func (a *App) SaveConfig(config Config) error {
	configPath := a.GetConfigPath()
	if configPath == "" {
		return fmt.Errorf("could not determine config path")
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Convert to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	
	// Write to file
	return os.WriteFile(configPath, data, 0644)
}
