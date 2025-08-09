package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx     context.Context
	logger  *log.Logger
	logFile *os.File
}

// Config represents the Claude Code Router configuration
type Config struct {
	APIKEY            interface{} `json:"APIKEY,omitempty"`
	PROXY_URL         interface{} `json:"PROXY_URL,omitempty"`
	HOST              interface{} `json:"HOST,omitempty"`
	PORT              interface{} `json:"PORT,omitempty"`
	API_TIMEOUT_MS    interface{} `json:"API_TIMEOUT_MS,omitempty"`
	LOG               interface{} `json:"LOG,omitempty"`
	NPM_GLOBAL_PREFIX interface{} `json:"NPM_GLOBAL_PREFIX,omitempty"`
	Providers         interface{} `json:"Providers,omitempty"`
	Router            interface{} `json:"Router,omitempty"`
}

// Provider represents a model provider configuration
type Provider struct {
	Name        interface{} `json:"name"`
	APIBaseURL  interface{} `json:"api_base_url"`
	APIKey      interface{} `json:"api_key"`
	Models      interface{} `json:"models"`
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
	app := &App{}
	app.initLogger()
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.logger != nil {
		a.logger.Printf("Application shutting down")
	}
	if a.logFile != nil {
		a.logFile.Close()
	}
}

// initLogger initializes the application logger
func (a *App) initLogger() {
	// Get the log file path
	logPath := a.GetAppLogPath()

	// Create the directory if it doesn't exist
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		// If we can't create the directory, we'll log to stderr
		log.SetOutput(os.Stderr)
		log.Printf("Failed to create log directory: %v", err)
		return
	}

	// Check if log file exists and rotate if needed
	a.rotateLogFileIfNeeded(logPath)

	// Open or create the log file (append mode)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// If we can't open the file, we'll log to stderr
		log.SetOutput(os.Stderr)
		log.Printf("Failed to open log file: %v", err)
		return
	}

	// Create a new logger with timestamp prefix
	a.logFile = logFile
	a.logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	// Also log to stderr for development
	multiWriter := io.MultiWriter(logFile, os.Stderr)
	log.SetOutput(multiWriter)

	// Log that the application has started
	a.logger.Printf("Application started")
}

// rotateLogFileIfNeeded rotates the log file if it exceeds 1MB
func (a *App) rotateLogFileIfNeeded(logPath string) {
	// Check if file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// File doesn't exist, no need to rotate
		return
	}

	// Get file info
	fileInfo, err := os.Stat(logPath)
	if err != nil {
		log.Printf("Failed to get log file info: %v", err)
		return
	}

	// Check if file size exceeds 1MB
	const maxSize = 1 * 1024 * 1024 // 1MB
	if fileInfo.Size() > maxSize {
		// Rotate the log file
		backupPath := logPath + ".old"

		// Remove old backup if it exists
		os.Remove(backupPath)

		// Rename current log to backup
		err := os.Rename(logPath, backupPath)
		if err != nil {
			log.Printf("Failed to rotate log file: %v", err)
			return
		}

		log.Printf("Log file rotated: %s", backupPath)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetAppLogPath returns the path to the application log file
func (a *App) GetAppLogPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Application log file should be in .claude-code-router directory
	configDir := filepath.Join(homeDir, ".claude-code-router")
	return filepath.Join(configDir, "config-manager.log")
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
		err := fmt.Errorf("could not determine config path")
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", err)
		}
		return config, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return empty config if file doesn't exist
		if a.logger != nil {
			a.logger.Printf("Config file does not exist at %s, returning empty config", configPath)
		}
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to read config file at %s: %v", configPath, err)
		}
		return config, err
	}

	// Parse JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to parse JSON config at %s: %v", configPath, err)
		}
		return config, err
	}

	// Handle all fields for compatibility
	config = a.processConfigFields(config)

	if a.logger != nil {
		a.logger.Printf("Successfully loaded config from %s", configPath)
	}

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

	// Handle PORT field
	if config.PORT != nil {
		config.PORT = a.convertToInt(config.PORT, 3456)
	} else {
		config.PORT = 3456 // default value
	}

	// Handle API_TIMEOUT_MS field compatibility
	if config.API_TIMEOUT_MS != nil {
		config.API_TIMEOUT_MS = a.convertToInt(config.API_TIMEOUT_MS, 600000)
	} else {
		config.API_TIMEOUT_MS = 600000 // default value
	}

	// Handle NPM_GLOBAL_PREFIX field
	if config.NPM_GLOBAL_PREFIX != nil {
		config.NPM_GLOBAL_PREFIX = a.convertToString(config.NPM_GLOBAL_PREFIX)
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
		err := fmt.Errorf("could not determine config path")
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", err)
		}
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to create directory %s: %v", dir, err)
		}
		return err
	}

	// Convert to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to marshal config to JSON: %v", err)
		}
		return err
	}

	// Write to file
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to write config to %s: %v", configPath, err)
		}
		return err
	}

	if a.logger != nil {
		a.logger.Printf("Successfully saved config to %s", configPath)
	}

	return nil
}

// ReadREADME reads the README.md file content
func (a *App) ReadREADME() (string, error) {
	// Get the executable directory
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Get the directory containing the executable
	execDir := filepath.Dir(execPath)

	// Construct the path to README.md (assuming it's in the same directory as the executable)
	readmePath := filepath.Join(execDir, "README.md")

	// Check if README.md exists in the executable directory
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		// If not found, try the parent directory (for development mode)
		readmePath = filepath.Join(filepath.Dir(execDir), "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			// If still not found, try current working directory
			readmePath = filepath.Join(".", "README.md")
			if _, err := os.Stat(readmePath); os.IsNotExist(err) {
				return "", fmt.Errorf("README.md not found")
			}
		}
	}

	// Read the file content
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ServiceStatus represents the status of the CCR service
type ServiceStatus struct {
	IsRunning bool `json:"isRunning"`
	PID       int  `json:"pid"`
}

// GetServiceStatus checks if the CCR service is running
func (a *App) GetServiceStatus() (ServiceStatus, error) {
	var status ServiceStatus

	// 首先获取配置以确定端口号
	config, err := a.LoadConfig()
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to load config for service status check: %v", err)
		}
		return status, err
	}

	// 获取端口号，如果没有配置则使用默认值3456
	port := 3456
	if config.PORT != nil {
		if portVal, ok := config.PORT.(int); ok {
			port = portVal
		} else if portVal, ok := config.PORT.(float64); ok {
			port = int(portVal)
		} else if portVal, ok := config.PORT.(string); ok {
			if p, err := strconv.Atoi(portVal); err == nil {
				port = p
			}
		}
	}

	if a.logger != nil {
		a.logger.Printf("Checking service status on port %d", port)
	}

	// 根据操作系统查询进程
	pid, isRunning, err := a.findProcessByPort(port)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to find process by port %d: %v", port, err)
		}
		return status, err
	}

	status.IsRunning = isRunning
	status.PID = pid

	if a.logger != nil {
		if isRunning {
			a.logger.Printf("Service is running with PID %d", pid)
		} else {
			a.logger.Printf("Service is not running")
		}
	}

	return status, nil
}

// findProcessByPort finds the process ID listening on the specified port
func (a *App) findProcessByPort(port int) (int, bool, error) {
	// 根据操作系统执行不同的命令
	pid, err := a.getProcessIDByPort(port)
	if err != nil {
		return 0, false, err
	}

	// 如果找到了PID，说明进程正在运行
	if pid > 0 {
		return pid, true, nil
	}

	return 0, false, nil
}

// getProcessIDByPort gets the process ID listening on the specified port
func (a *App) getProcessIDByPort(port int) (int, error) {
	var cmd *exec.Cmd

	// 根据操作系统选择命令
	if isWindows() {
		// Windows: 使用 netstat 命令
		cmd = exec.Command("netstat", "-ano")
	} else {
		// Unix/Linux/macOS: 使用 lsof 命令
		cmd = exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t")
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, nil
	}

	if isWindows() {
		return a.parseWindowsNetstatOutput(string(output), port)
	} else {
		return a.parseUnixLsofOutput(string(output))
	}
}

// parseWindowsNetstatOutput parses the output of Windows netstat command
func (a *App) parseWindowsNetstatOutput(output string, port int) (int, error) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// 查找监听在指定端口的行
		if strings.Contains(line, fmt.Sprintf("0.0.0.0:%d", port)) ||
			strings.Contains(line, fmt.Sprintf("[::]:%d", port)) ||
			strings.Contains(line, fmt.Sprintf("127.0.0.1:%d", port)) ||
			strings.Contains(line, fmt.Sprintf("[::1]:%d", port)) {
			// 查找PID（在最后一列）
			fields := strings.Fields(line)
			if len(fields) > 0 {
				pidStr := fields[len(fields)-1]
				if pid, err := strconv.Atoi(pidStr); err == nil {
					return pid, nil
				}
			}
		}
	}

	return 0, nil
}

// parseUnixLsofOutput parses the output of Unix lsof command
func (a *App) parseUnixLsofOutput(output string) (int, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 && lines[0] != "" {
		pidStr := lines[0]
		if pid, err := strconv.Atoi(pidStr); err == nil {
			return pid, nil
		}
	}

	return 0, nil
}

// isWindows checks if the current OS is Windows
func isWindows() bool {
	return strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") ||
		strings.HasSuffix(os.Getenv("PATH"), ";")
}

// findCCRPath finds the CCR command path based on configuration
func (a *App) findCCRPath() (string, error) {
	// 加载配置以获取npm全局安装目录
	config, err := a.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %v", err)
	}

	// 查找ccr命令的绝对路径
	var ccrPath string
	if config.NPM_GLOBAL_PREFIX != nil && config.NPM_GLOBAL_PREFIX != "" {
		// 如果配置了npm全局安装目录，优先在此目录中查找ccr
		npmPrefix := config.NPM_GLOBAL_PREFIX.(string)
		ccrPath = filepath.Join(npmPrefix, "ccr")
		a.logger.Printf("WARNING: CCR command not found in %s", npmPrefix)
		// 检查文件是否存在
		if _, err := os.Stat(ccrPath); os.IsNotExist(err) {
			// 如果文件不存在，尝试添加.exe后缀（Windows）
			a.logger.Printf("WARNING: start find CCR from PATH")

			if isWindows() {
				ccrPath = filepath.Join(npmPrefix, "ccr.exe")
				if _, err := os.Stat(ccrPath); os.IsNotExist(err) {
					// 如果还是找不到，回退到PATH中查找
					ccrPath, err = exec.LookPath("ccr")
					if err != nil {
						ccrPath = "ccr"
					}
				}
			} else {
				// 如果还是找不到，回退到PATH中查找
				ccrPath, err = exec.LookPath("ccr")
				if err != nil {
					ccrPath = "ccr"
				}
				a.logger.Printf("Found CCR at any path")
			}
		}
	} else {
		// 如果没有配置npm全局安装目录，在PATH中查找
		ccrPath, err = exec.LookPath("ccr")
		a.logger.Printf("not config npm global prefix,start find to PATH")
		if err != nil {
			// 如果在PATH中找不到，尝试使用完整路径
			ccrPath = "ccr"
		}
	}

	return ccrPath, nil
}

// StartService starts the CCR service
func (a *App) StartService() error {
	if a.logger != nil {
		a.logger.Printf("Starting CCR service")
	}

	// 查找ccr命令的绝对路径
	ccrPath, err := a.findCCRPath()
	if err != nil {
		errMsg := fmt.Errorf("failed to find CCR path: %v", err)
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", errMsg)
		}
		return errMsg
	}

	if a.logger != nil {
		a.logger.Printf("Found CCR at path: %s", ccrPath)
	}

	// 使用 ccr start 命令启动服务
	cmd := exec.Command(ccrPath, "start")

	// 设置命令属性
	cmd.Dir = "."

	// 设置命令在后台运行，避免创建可见窗口
	cmd.SysProcAttr = getSysProcAttr()

	if a.logger != nil {
		a.logger.Printf("Executing command: %s %s", ccrPath, "start")
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if outputStr == "" {
			outputStr = "No output"
		}

		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to start CCR service. Output: %s", outputStr)
		}

		// 检查是否是HTTP相关错误
		if strings.Contains(outputStr, "HTTP/1.1 400 Bad Request") {
			errMsg := fmt.Errorf("HTTP communication error when starting CCR service. This may be a Wails internal issue. Command output: %s", outputStr)
			return errMsg
		}
		// 检查是否是常见的Windows权限错误
		if strings.Contains(outputStr, "Access is denied") || strings.Contains(outputStr, "拒绝访问") {
			errMsg := fmt.Errorf("permission denied when starting CCR service. Try running this application as administrator. Command output: %s", outputStr)
			return errMsg
		}
		// 检查是否是文件未找到错误
		if strings.Contains(outputStr, "not found") || strings.Contains(outputStr, "not recognized") || strings.Contains(outputStr, "无法将") {
			errMsg := fmt.Errorf("CCR command not found. Please ensure CCR is properly installed and in your PATH. Command output: %s", outputStr)
			return errMsg
		}
		// 增强错误信息，包含更多上下文
		errMsg := fmt.Errorf("failed to start CCR service: %v, output: %s, working dir: %s, ccr path: %s", err, outputStr, cmd.Dir, ccrPath)
		return errMsg
	}

	if a.logger != nil {
		a.logger.Printf("Successfully started CCR service")
	}

	return nil
}

// StopService stops the CCR service
func (a *App) StopService() error {
	if a.logger != nil {
		a.logger.Printf("Stopping CCR service")
	}

	// 查找ccr命令的绝对路径
	ccrPath, err := a.findCCRPath()
	if err != nil {
		errMsg := fmt.Errorf("failed to find CCR path: %v", err)
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", errMsg)
		}
		return errMsg
	}

	if a.logger != nil {
		a.logger.Printf("Found CCR at path: %s", ccrPath)
	}

	// 使用 ccr stop 命令停止服务
	cmd := exec.Command(ccrPath, "stop")

	// 设置命令属性
	cmd.Dir = "."

	// 设置命令在后台运行，避免创建可见窗口
	cmd.SysProcAttr = getSysProcAttr()

	if a.logger != nil {
		a.logger.Printf("Executing command: %s %s", ccrPath, "stop")
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if outputStr == "" {
			outputStr = "No output"
		}

		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to stop CCR service. Output: %s", outputStr)
		}

		// 检查是否是HTTP相关错误
		if strings.Contains(outputStr, "HTTP/1.1 400 Bad Request") {
			errMsg := fmt.Errorf("HTTP communication error when stopping CCR service. This may be a Wails internal issue. Command output: %s", outputStr)
			return errMsg
		}
		// 检查是否是常见的Windows权限错误
		if strings.Contains(outputStr, "Access is denied") || strings.Contains(outputStr, "拒绝访问") {
			errMsg := fmt.Errorf("permission denied when stopping CCR service. Try running this application as administrator. Command output: %s", outputStr)
			return errMsg
		}
		// 检查是否是文件未找到错误
		if strings.Contains(outputStr, "not found") || strings.Contains(outputStr, "not recognized") || strings.Contains(outputStr, "无法将") {
			errMsg := fmt.Errorf("CCR command not found. Please ensure CCR is properly installed and in your PATH. Command output: %s", outputStr)
			return errMsg
		}
		// 增强错误信息，包含更多上下文
		errMsg := fmt.Errorf("failed to stop CCR service: %v, output: %s, working dir: %s, ccr path: %s", err, outputStr, cmd.Dir, ccrPath)
		return errMsg
	}

	if a.logger != nil {
		a.logger.Printf("Successfully stopped CCR service")
	}

	return nil
}

// RestartService restarts the CCR service
func (a *App) RestartService() error {
	if a.logger != nil {
		a.logger.Printf("Restarting CCR service")
	}

	// 查找ccr命令的绝对路径
	ccrPath, err := a.findCCRPath()
	if err != nil {
		errMsg := fmt.Errorf("failed to find CCR path: %v", err)
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", errMsg)
		}
		return errMsg
	}

	if a.logger != nil {
		a.logger.Printf("Found CCR at path: %s", ccrPath)
	}

	// 使用 ccr restart 命令重启服务
	cmd := exec.Command(ccrPath, "restart")

	// 设置命令属性
	cmd.Dir = "."

	// 设置命令在后台运行，避免创建可见窗口
	cmd.SysProcAttr = getSysProcAttr()

	if a.logger != nil {
		a.logger.Printf("Executing command: %s %s", ccrPath, "restart")
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if outputStr == "" {
			outputStr = "No output"
		}

		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to restart CCR service. Output: %s", outputStr)
		}

		// 检查是否是HTTP相关错误
		if strings.Contains(outputStr, "HTTP/1.1 400 Bad Request") {
			errMsg := fmt.Errorf("HTTP communication error when restarting CCR service. This may be a Wails internal issue. Command output: %s", outputStr)
			return errMsg
		}
		// 检查是否是常见的Windows权限错误
		if strings.Contains(outputStr, "Access is denied") || strings.Contains(outputStr, "拒绝访问") {
			errMsg := fmt.Errorf("permission denied when restarting CCR service. Try running this application as administrator. Command output: %s", outputStr)
			return errMsg
		}
		// 检查是否是文件未找到错误
		if strings.Contains(outputStr, "not found") || strings.Contains(outputStr, "not recognized") || strings.Contains(outputStr, "无法将") {
			errMsg := fmt.Errorf("CCR command not found. Please ensure CCR is properly installed and in your PATH. Command output: %s", outputStr)
			return errMsg
		}
		// 增强错误信息，包含更多上下文
		errMsg := fmt.Errorf("failed to restart CCR service: %v, output: %s, working dir: %s, ccr path: %s", err, outputStr, cmd.Dir, ccrPath)
		return errMsg
	}

	if a.logger != nil {
		a.logger.Printf("Successfully restarted CCR service")
	}

	return nil
}

// GetLogPath returns the path to the CCR log file
func (a *App) GetLogPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// 日志文件应该在 .claude-code-router 目录下
	configDir := filepath.Join(homeDir, ".claude-code-router")
	return filepath.Join(configDir, "claude-code-router.log")
}

// GetCCRVersion returns the version of the CCR service
func (a *App) GetCCRVersion() (string, error) {
	if a.logger != nil {
		a.logger.Printf("Getting CCR version")
	}

	// 尝试从package.json文件读取版本信息
	version, err := a.getCCRViaPackageJSON()
	if err == nil && version != "" {
		if a.logger != nil {
			a.logger.Printf("Found version in package.json: %s", version)
		}
		return version, nil
	}

	if a.logger != nil {
		a.logger.Printf("Failed to get version from package.json, trying command line: %v", err)
	}

	// 查找ccr命令的绝对路径（使用通用方法）
	ccrPath, err := a.findCCRPath()
	if err != nil {
		errMsg := fmt.Errorf("failed to find CCR path: %v", err)
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", errMsg)
		}
		return "", errMsg
	}

	if a.logger != nil {
		a.logger.Printf("Found CCR at path: %s", ccrPath)
	}

	// 使用 ccr -v 命令获取版本号
	cmd := exec.Command(ccrPath, "-v")

	// 设置命令在后台运行，避免创建可见窗口
	cmd.SysProcAttr = getSysProcAttr()

	if a.logger != nil {
		a.logger.Printf("Executing command: %s %s", ccrPath, "-v")
	}

	// 执行命令
	output, err := cmd.Output()
	if err != nil {
		errMsg := fmt.Errorf("failed to get CCR version: %v, output: %s", err, string(output))
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", errMsg)
		}
		return "", errMsg
	}

	// 清理输出，移除换行符
	version = strings.TrimSpace(string(output))

	if a.logger != nil {
		a.logger.Printf("Successfully got CCR version: %s", version)
	}

	return version, nil
}

// getCCRViaPackageJSON tries to read CCR version from package.json
func (a *App) getCCRViaPackageJSON() (string, error) {
	// 查找ccr命令的路径
	ccrPath, err := a.findCCRPath()
	if err != nil {
		return "", fmt.Errorf("failed to find CCR path: %v", err)
	}
	a.logger.Printf("Found CCR at path: %s", ccrPath)
	// 获取ccr命令所在目录
	ccrDir := filepath.Dir(ccrPath)

	// 构建package.json路径
	packageJSONPath := filepath.Join(ccrDir, "node_modules", "@musistudio", "claude-code-router", "package.json")

	// 检查文件是否存在
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return "", fmt.Errorf("package.json not found at: %s", packageJSONPath)
	}

	// 读取文件内容
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return "", fmt.Errorf("failed to read package.json: %v", err)
	}

	// 解析JSON
	var packageInfo map[string]interface{}
	if err := json.Unmarshal(data, &packageInfo); err != nil {
		return "", fmt.Errorf("failed to parse package.json: %v", err)
	}

	// 获取版本号
	if version, ok := packageInfo["version"].(string); ok {
		return version, nil
	}

	return "", fmt.Errorf("version not found in package.json")
}

// ReadLogs reads the CCR log file and limits to last 500 lines
func (a *App) ReadLogs() (string, error) {
	logPath := a.GetLogPath()
	if logPath == "" {
		err := fmt.Errorf("could not determine log path")
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", err)
		}
		return "", err
	}

	if a.logger != nil {
		a.logger.Printf("Reading logs from %s", logPath)
	}

	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// If log file doesn't exist, return empty string instead of error
		if a.logger != nil {
			a.logger.Printf("Log file does not exist at %s, returning empty string", logPath)
		}
		return "", nil
	}

	// 限制读取文件大小，避免读取过大的文件
	fileInfo, err := os.Stat(logPath)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to get log file info for %s: %v", logPath, err)
		}
		return "", fmt.Errorf("failed to get log file info: %v", err)
	}

	// 如果文件大于1MB，只读取最后的部分
	const maxFileSize = 20 * 1024
	if fileInfo.Size() > maxFileSize {
		if a.logger != nil {
			a.logger.Printf("Log file is large (%d bytes), reading last %d bytes", fileInfo.Size(), maxFileSize)
		}

		// 打开文件
		file, err := os.Open(logPath)
		if err != nil {
			if a.logger != nil {
				a.logger.Printf("ERROR: Failed to open log file %s: %v", logPath, err)
			}
			return "", fmt.Errorf("failed to open log file: %v", err)
		}
		defer file.Close()

		// 移动到文件末尾前1MB的位置
		startPos := fileInfo.Size() - maxFileSize
		if startPos < 0 {
			startPos = 0
		}

		// 读取最后1MB的内容
		buffer := make([]byte, maxFileSize)
		n, err := file.ReadAt(buffer, startPos)
		if err != nil && err != io.EOF {
			if a.logger != nil {
				a.logger.Printf("ERROR: Failed to read log file %s: %v", logPath, err)
			}
			return "", fmt.Errorf("failed to read log file: %v", err)
		}

		// 转换为字符串并限制到最后一行
		content := string(buffer[:n])
		lines := strings.Split(content, "\n")

		// 如果有足够多的行，取最后500行
		if len(lines) > 500 {
			lines = lines[len(lines)-500:]
		}

		if a.logger != nil {
			a.logger.Printf("Successfully read %d lines from log file", len(lines))
		}

		return strings.Join(lines, "\n"), nil
	}

	// 文件较小，直接读取
	content, err := os.ReadFile(logPath)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to read log file %s: %v", logPath, err)
		}
		return "", fmt.Errorf("failed to read log file: %v", err)
	}

	// Limit to last 500 lines
	logContent := string(content)
	lines := strings.Split(logContent, "\n")

	// If there are more than 500 lines, keep only the last 500
	if len(lines) > 500 {
		originalCount := len(lines)
		lines = lines[len(lines)-500:]
		if a.logger != nil {
			a.logger.Printf("Limited log output from %d to 500 lines", originalCount)
		}
	}

	if a.logger != nil {
		a.logger.Printf("Successfully read %d lines from log file", len(lines))
	}

	return strings.Join(lines, "\n"), nil
}

// ClearLogs clears the CCR log file
func (a *App) ClearLogs() error {
	logPath := a.GetLogPath()
	if logPath == "" {
		err := fmt.Errorf("could not determine log path")
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", err)
		}
		return err
	}

	if a.logger != nil {
		a.logger.Printf("Clearing logs at %s", logPath)
	}

	// Ensure the directory exists
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to create log directory %s: %v", dir, err)
		}
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Truncate the log file
	err := os.WriteFile(logPath, []byte(""), 0644)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to clear log file %s: %v", logPath, err)
		}
		return err
	}

	if a.logger != nil {
		a.logger.Printf("Successfully cleared log file at %s", logPath)
	}

	return nil
}

// ClearAppLogs clears the application log file
func (a *App) ClearAppLogs() error {
	logPath := a.GetAppLogPath()
	if logPath == "" {
		err := fmt.Errorf("could not determine app log path")
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", err)
		}
		return err
	}

	if a.logger != nil {
		a.logger.Printf("Clearing app logs at %s", logPath)
	}

	// Ensure the directory exists
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to create app log directory %s: %v", dir, err)
		}
		return fmt.Errorf("failed to create app log directory: %v", err)
	}

	// Truncate the log file
	err := os.WriteFile(logPath, []byte(""), 0644)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to clear app log file %s: %v", logPath, err)
		}
		return err
	}

	if a.logger != nil {
		a.logger.Printf("Successfully cleared app log file at %s", logPath)
	}

	return nil
}

// TestLogging is a utility function to test if logging is working
func (a *App) TestLogging() string {
	if a.logger != nil {
		a.logger.Printf("Test log entry at %s", time.Now().Format(time.RFC3339))
		return "Logging test successful"
	}
	return "Logger not initialized"
}

// ReadAppLogs reads the application log file and limits to last 500 lines
func (a *App) ReadAppLogs() (string, error) {
	logPath := a.GetAppLogPath()
	if logPath == "" {
		err := fmt.Errorf("could not determine app log path")
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", err)
		}
		return "", err
	}

	if a.logger != nil {
		a.logger.Printf("Reading app logs from %s", logPath)
	}

	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// If log file doesn't exist, return empty string instead of error
		if a.logger != nil {
			a.logger.Printf("App log file does not exist at %s, returning empty string", logPath)
		}
		return "", nil
	}

	// 限制读取文件大小，避免读取过大的文件
	fileInfo, err := os.Stat(logPath)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to get app log file info for %s: %v", logPath, err)
		}
		return "", fmt.Errorf("failed to get app log file info: %v", err)
	}

	// 如果文件大于1MB，只读取最后的部分
	const maxFileSize = 20 * 1024
	if fileInfo.Size() > maxFileSize {
		if a.logger != nil {
			a.logger.Printf("App log file is large (%d bytes), reading last %d bytes", fileInfo.Size(), maxFileSize)
		}

		// 打开文件
		file, err := os.Open(logPath)
		if err != nil {
			if a.logger != nil {
				a.logger.Printf("ERROR: Failed to open app log file %s: %v", logPath, err)
			}
			return "", fmt.Errorf("failed to open app log file: %v", err)
		}
		defer file.Close()

		// 移动到文件末尾前1MB的位置
		startPos := fileInfo.Size() - maxFileSize
		if startPos < 0 {
			startPos = 0
		}

		// 读取最后1MB的内容
		buffer := make([]byte, maxFileSize)
		n, err := file.ReadAt(buffer, startPos)
		if err != nil && err != io.EOF {
			if a.logger != nil {
				a.logger.Printf("ERROR: Failed to read app log file %s: %v", logPath, err)
			}
			return "", fmt.Errorf("failed to read app log file: %v", err)
		}

		// 转换为字符串并限制到最后一行
		content := string(buffer[:n])
		lines := strings.Split(content, "\n")

		// 如果有足够多的行，取最后500行
		if len(lines) > 500 {
			lines = lines[len(lines)-500:]
		}

		if a.logger != nil {
			a.logger.Printf("Successfully read %d lines from app log file", len(lines))
		}

		return strings.Join(lines, "\n"), nil
	}

	// 文件较小，直接读取
	content, err := os.ReadFile(logPath)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to read app log file %s: %v", logPath, err)
		}
		return "", fmt.Errorf("failed to read app log file: %v", err)
	}

	// Limit to last 500 lines
	logContent := string(content)
	lines := strings.Split(logContent, "\n")

	// If there are more than 500 lines, keep only the last 500
	if len(lines) > 500 {
		originalCount := len(lines)
		lines = lines[len(lines)-500:]
		if a.logger != nil {
			a.logger.Printf("Limited app log output from %d to 500 lines", originalCount)
		}
	}

	if a.logger != nil {
		a.logger.Printf("Successfully read %d lines from app log file", len(lines))
	}

	return strings.Join(lines, "\n"), nil
}

// GetAppVersion returns the version of the application
func (a *App) GetAppVersion() string {
	// For now, return a fixed version
	// In the future, this could read from a version file or be set at build time
	return "v2.0.3"
}

// GetLatestVersionFromGitHub checks GitHub for the latest version tag
// Implements fallback mechanisms to handle rate limits
func (a *App) GetLatestVersionFromGitHub() (string, error) {
	// Try primary method first - GitHub API
	version, err := a.getLatestVersionFromGitHubAPI()
	if err == nil {
		return version, nil
	}
	
	// If primary method fails, try fallback - RSS feed
	if a.logger != nil {
		a.logger.Printf("INFO: Primary GitHub API method failed, trying RSS feed fallback: %v", err)
	}
	version, err = a.getLatestVersionFromRSS()
	if err == nil {
		return version, nil
	}
	
	// All methods failed
	if a.logger != nil {
		a.logger.Printf("ERROR: All methods failed to fetch latest version: %v", err)
	}
	return "", fmt.Errorf("failed to fetch latest version from all sources: %v", err)
}

// getLatestVersionFromGitHubAPI fetches the latest version using GitHub API
func (a *App) getLatestVersionFromGitHubAPI() (string, error) {
	// GitHub API URL for tags
	url := "https://api.github.com/repos/ayuayue/ccr-config-manager/tags"
	
	// Create HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	
	// Add user agent to avoid being blocked
	req.Header.Set("User-Agent", "CCR-Config-Manager")
	
	resp, err := client.Do(req)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to fetch tags from GitHub API: %v", err)
		}
		return "", fmt.Errorf("failed to fetch tags from GitHub API: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		if a.logger != nil {
			a.logger.Printf("ERROR: GitHub API request failed with status: %d", resp.StatusCode)
		}
		return "", fmt.Errorf("GitHub API request failed with status: %d", resp.StatusCode)
	}
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to read response body: %v", err)
		}
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	
	// Parse JSON response
	var tags []map[string]interface{}
	if err := json.Unmarshal(body, &tags); err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to parse JSON response: %v", err)
		}
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}
	
	// Check if we have any tags
	if len(tags) == 0 {
		err := fmt.Errorf("no tags found for repository")
		if a.logger != nil {
			a.logger.Printf("ERROR: %v", err)
		}
		return "", err
	}
	
	// Get the first tag (latest)
	if name, ok := tags[0]["name"].(string); ok {
		if a.logger != nil {
			a.logger.Printf("INFO: Latest version found via API: %s", name)
		}
		return name, nil
	}
	
	err = fmt.Errorf("failed to extract tag name from response")
	if a.logger != nil {
		a.logger.Printf("ERROR: %v", err)
	}
	return "", err
}

// getLatestVersionFromRSS fetches the latest version using GitHub RSS feed
func (a *App) getLatestVersionFromRSS() (string, error) {
	// GitHub RSS feed URL for releases
	url := "https://github.com/ayuayue/ccr-config-manager/releases.atom"
	
	// Create HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create RSS request: %v", err)
	}
	
	// Add user agent to avoid being blocked
	req.Header.Set("User-Agent", "CCR-Config-Manager")
	
	resp, err := client.Do(req)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to fetch RSS feed: %v", err)
		}
		return "", fmt.Errorf("failed to fetch RSS feed: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		if a.logger != nil {
			a.logger.Printf("ERROR: RSS feed request failed with status: %d", resp.StatusCode)
		}
		return "", fmt.Errorf("RSS feed request failed with status: %d", resp.StatusCode)
	}
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to read RSS response body: %v", err)
		}
		return "", fmt.Errorf("failed to read RSS response body: %v", err)
	}
	
	// Parse RSS feed to extract latest version
	// Simple approach: look for first <title> tag that contains a version-like string
	content := string(body)
	
	// Find all title tags
	titleStart := strings.Index(content, "<title>")
	if titleStart == -1 {
		return "", fmt.Errorf("no title tags found in RSS feed")
	}
	
	// Look for version pattern in titles
	// Skip the first title which is usually the feed title
	content = content[titleStart+7:]
	titleStart = strings.Index(content, "<title>")
	if titleStart == -1 {
		return "", fmt.Errorf("no release titles found in RSS feed")
	}
	
	content = content[titleStart+7:]
	titleEnd := strings.Index(content, "</title>")
	if titleEnd == -1 {
		return "", fmt.Errorf("malformed title tag in RSS feed")
	}
	
	title := content[:titleEnd]
	
	// Extract version from title (assuming format "Release vX.X.X" or "vX.X.X")
	// Look for version pattern
	versionRegex := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	matches := versionRegex.FindStringSubmatch(title)
	if len(matches) > 0 {
		if a.logger != nil {
			a.logger.Printf("INFO: Latest version found via RSS: %s", matches[0])
		}
		return matches[0], nil
	}
	
	// Try simpler pattern for any v followed by numbers and dots
	simpleRegex := regexp.MustCompile(`v[\d.]+`)
	matches = simpleRegex.FindStringSubmatch(title)
	if len(matches) > 0 {
		if a.logger != nil {
			a.logger.Printf("INFO: Latest version found via RSS (simple pattern): %s", matches[0])
		}
		return matches[0], nil
	}
	
	return "", fmt.Errorf("failed to extract version from RSS title: %s", title)
}

// CompareVersions compares two version strings
func (a *App) CompareVersions(current, latest string) bool {
	// Simple version comparison assuming semantic versioning (e.g., v1.2.3)
	// Remove 'v' prefix if present
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")
	
	// Split version strings into parts
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")
	
	// Compare each part numerically
	for i := 0; i < len(currentParts) && i < len(latestParts); i++ {
		currentNum, err1 := strconv.Atoi(currentParts[i])
		latestNum, err2 := strconv.Atoi(latestParts[i])
		
		if err1 != nil || err2 != nil {
			// If we can't parse a part as a number, fall back to string comparison
			if latestParts[i] > currentParts[i] {
				return true
			} else if latestParts[i] < currentParts[i] {
				return false
			}
			continue
		}
		
		if latestNum > currentNum {
			return true
		} else if latestNum < currentNum {
			return false
		}
	}
	
	// If all compared parts are equal, check if latest has more parts
	return len(latestParts) > len(currentParts)
}

// DownloadUpdate downloads the latest version from GitHub
func (a *App) DownloadUpdate(version string) (string, error) {
	// Create downloads directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	
	downloadsDir := filepath.Join(homeDir, "Downloads")
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create downloads directory: %v", err)
	}
	
	// Determine the appropriate asset name based on the operating system
	assetName := ""
	fileName := ""
	
	// Get the operating system
	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	
	// Set asset name and file name based on OS
	switch goos {
	case "windows":
		assetName = fmt.Sprintf("claudeConfigManager_%s_windows_amd64_installer.exe", strings.TrimPrefix(version, "v"))
		fileName = fmt.Sprintf("claudeConfigManager_%s_windows_amd64_installer.exe", strings.TrimPrefix(version, "v"))
	case "darwin":
		// For macOS, default to arm64 (Apple Silicon)
		arch := os.Getenv("GOARCH")
		if arch == "" {
			arch = runtime.GOARCH
		}
		if arch == "arm64" {
			assetName = fmt.Sprintf("claudeConfigManager_%s_mac_arm64.zip", strings.TrimPrefix(version, "v"))
			fileName = fmt.Sprintf("claudeConfigManager_%s_mac_arm64.zip", strings.TrimPrefix(version, "v"))
		} else {
			assetName = fmt.Sprintf("claudeConfigManager_%s_mac_intel.zip", strings.TrimPrefix(version, "v"))
			fileName = fmt.Sprintf("claudeConfigManager_%s_mac_intel.zip", strings.TrimPrefix(version, "v"))
		}
	case "linux":
		assetName = fmt.Sprintf("claudeConfigManager_%s_linux_amd64.tar.gz", strings.TrimPrefix(version, "v"))
		fileName = fmt.Sprintf("claudeConfigManager_%s_linux_amd64.tar.gz", strings.TrimPrefix(version, "v"))
	default:
		// Default to Windows installer for unknown platforms
		assetName = fmt.Sprintf("claudeConfigManager_%s_windows_amd64_installer.exe", strings.TrimPrefix(version, "v"))
		fileName = fmt.Sprintf("claudeConfigManager_%s_windows_amd64_installer.exe", strings.TrimPrefix(version, "v"))
	}
	
	// Construct download URL
	url := fmt.Sprintf("https://github.com/ayuayue/ccr-config-manager/releases/download/%s/%s", version, assetName)
	
	// Create file path
	filePath := filepath.Join(downloadsDir, fileName)
	
	if a.logger != nil {
		a.logger.Printf("INFO: Downloading update from %s to %s", url, filePath)
	}
	
	// Create HTTP request with timeout
	client := &http.Client{
		Timeout: 300 * time.Second, // 5 minutes timeout for download
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %v", err)
	}
	
	// Add user agent to avoid being blocked
	req.Header.Set("User-Agent", "CCR-Config-Manager")
	
	resp, err := client.Do(req)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to download update: %v", err)
		}
		return "", fmt.Errorf("failed to download update: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		if a.logger != nil {
			a.logger.Printf("ERROR: Download request failed with status: %d", resp.StatusCode)
		}
		return "", fmt.Errorf("download request failed with status: %d", resp.StatusCode)
	}
	
	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to create file: %v", err)
		}
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	
	// Copy response body to file with progress (optional)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		if a.logger != nil {
			a.logger.Printf("ERROR: Failed to save file: %v", err)
		}
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	
	if a.logger != nil {
		a.logger.Printf("INFO: Update downloaded successfully to %s", filePath)
	}
	
	return filePath, nil
}
