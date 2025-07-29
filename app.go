package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
	"strconv"
)

// App struct
type App struct {
	ctx context.Context
}

// Config represents the Claude Code Router configuration
type Config struct {
	APIKEY              interface{} `json:"APIKEY,omitempty"`
	PROXY_URL          interface{} `json:"PROXY_URL,omitempty"`
	HOST               interface{} `json:"HOST,omitempty"`
	API_TIMEOUT_MS     interface{} `json:"API_TIMEOUT_MS,omitempty"`
	LOG                interface{} `json:"LOG,omitempty"`
	Providers          interface{} `json:"Providers,omitempty"`
	Router             interface{} `json:"Router,omitempty"`
}

// Provider represents a model provider configuration
type Provider struct {
	Name       interface{} `json:"name"`
	APIBaseURL interface{} `json:"api_base_url"`
	APIKey     interface{} `json:"api_key"`
	Models     interface{} `json:"models"`
	Transformer interface{} `json:"transformer,omitempty"`
}

// Router represents the routing configuration
type Router struct {
	Default              interface{} `json:"default,omitempty"`
	Background           interface{} `json:"background,omitempty"`
	Think                interface{} `json:"think,omitempty"`
	LongContext          interface{} `json:"longContext,omitempty"`
	LongContextThreshold interface{} `json:"longContextThreshold,omitempty"`
	WebSearch            interface{} `json:"webSearch,omitempty"`
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
	
	// Handle all fields for compatibility
	config = a.processConfigFields(config)
	
	return config, nil
}

// processConfigFields handles type conversion for all config fields
func (a *App) processConfigFields(config Config) Config {
	// Handle LOG field compatibility (string or bool)
	if config.LOG != nil {
		config.LOG = a.convertToBool(config.LOG)
	}
	
	// Handle APIKEY field
	if config.APIKEY != nil {
		config.APIKEY = a.convertToString(config.APIKEY)
	}
	
	// Handle PROXY_URL field
	if config.PROXY_URL != nil {
		config.PROXY_URL = a.convertToString(config.PROXY_URL)
	}
	
	// Handle HOST field
	if config.HOST != nil {
		config.HOST = a.convertToString(config.HOST)
	}
	
	// Handle API_TIMEOUT_MS field compatibility
	if config.API_TIMEOUT_MS != nil {
		config.API_TIMEOUT_MS = a.convertToInt(config.API_TIMEOUT_MS, 600000)
	} else {
		config.API_TIMEOUT_MS = 600000 // default value
	}
	
	// Handle Providers field
	if config.Providers != nil {
		config.Providers = a.processProviders(config.Providers)
	}
	
	// Handle Router field
	if config.Router != nil {
		config.Router = a.processRouter(config.Router)
	}
	
	return config
}

// convertToBool converts various types to bool
func (a *App) convertToBool(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return v == "true"
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}

// convertToString converts various types to string
func (a *App) convertToString(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// convertToInt converts various types to int
func (a *App) convertToInt(value interface{}, defaultValue int) interface{} {
	switch v := value.(type) {
	case string:
		if timeout, err := strconv.Atoi(v); err == nil {
			return timeout
		}
		return defaultValue
	case float64:
		return int(v)
	case int:
		return v
	default:
		return defaultValue
	}
}

// processProviders handles type conversion for providers
func (a *App) processProviders(providers interface{}) interface{} {
	// If it's already a slice of Provider structs, return as is
	if providerSlice, ok := providers.([]Provider); ok {
		return providerSlice
	}
	
	// If it's a slice of interfaces, process each one
	if providerSlice, ok := providers.([]interface{}); ok {
		result := make([]interface{}, len(providerSlice))
		for i, provider := range providerSlice {
			if providerMap, ok := provider.(map[string]interface{}); ok {
				result[i] = a.processProviderMap(providerMap)
			} else {
				result[i] = provider
			}
		}
		return result
	}
	
	// If it's a slice of maps, process each one
	if providerSlice, ok := providers.([]map[string]interface{}); ok {
		result := make([]interface{}, len(providerSlice))
		for i, providerMap := range providerSlice {
			result[i] = a.processProviderMap(providerMap)
		}
		return result
	}
	
	return providers
}

// processProviderMap handles type conversion for a single provider map
func (a *App) processProviderMap(providerMap map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})
	
	for key, value := range providerMap {
		switch key {
		case "name":
			processed[key] = a.convertToString(value)
		case "api_base_url":
			processed[key] = a.convertToString(value)
		case "api_key":
			processed[key] = a.convertToString(value)
		case "models":
			processed[key] = a.processModels(value)
		default:
			processed[key] = value
		}
	}
	
	return processed
}

// processModels handles type conversion for models
func (a *App) processModels(models interface{}) interface{} {
	// If it's already a slice of strings, return as is
	if modelSlice, ok := models.([]string); ok {
		return modelSlice
	}
	
	// If it's a slice of interfaces, convert each one to string
	if modelSlice, ok := models.([]interface{}); ok {
		result := make([]string, len(modelSlice))
		for i, model := range modelSlice {
			result[i] = a.convertToString(model).(string)
		}
		return result
	}
	
	return models
}

// processRouter handles type conversion for router
func (a *App) processRouter(router interface{}) interface{} {
	// If it's already a Router struct, return as is
	if routerStruct, ok := router.(Router); ok {
		return routerStruct
	}
	
	// If it's a map, process each field
	if routerMap, ok := router.(map[string]interface{}); ok {
		processed := make(map[string]interface{})
		
		for key, value := range routerMap {
			switch key {
			case "default", "background", "think", "longContext", "webSearch":
				processed[key] = a.convertToString(value)
			case "longContextThreshold":
				processed[key] = a.convertToInt(value, 60000)
			default:
				processed[key] = value
			}
		}
		
		return processed
	}
	
	return router
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
