package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"encoding/json"
	"strconv"
	"strings"
	"runtime"
	"syscall"
	"io"
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
	PORT               interface{} `json:"PORT,omitempty"`
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
	IsRunning bool   `json:"isRunning"`
	PID       int    `json:"pid"`
}

// GetServiceStatus checks if the CCR service is running
func (a *App) GetServiceStatus() (ServiceStatus, error) {
	var status ServiceStatus
	
	// 首先获取配置以确定端口号
	config, err := a.LoadConfig()
	if err != nil {
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
	
	// 根据操作系统查询进程
	pid, isRunning, err := a.findProcessByPort(port)
	if err != nil {
		return status, err
	}
	
	status.IsRunning = isRunning
	status.PID = pid
	
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

// StartService starts the CCR service
func (a *App) StartService() error {
	// 查找ccr命令的绝对路径
	ccrPath, err := exec.LookPath("ccr")
	if err != nil {
		// 如果在PATH中找不到，尝试使用完整路径
		ccrPath = "ccr"
	}
	
	// 使用 ccr start 命令启动服务
	cmd := exec.Command(ccrPath, "start")
	
	// 设置命令属性
	cmd.Dir = "."
	
	// 设置命令在后台运行，避免创建可见窗口
	cmd.SysProcAttr = getSysProcAttr()
	
	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if outputStr == "" {
			outputStr = "No output"
		}
		// 检查是否是HTTP相关错误
		if strings.Contains(outputStr, "HTTP/1.1 400 Bad Request") {
			return fmt.Errorf("HTTP communication error when starting CCR service. This may be a Wails internal issue. Command output: %s", outputStr)
		}
		// 检查是否是常见的Windows权限错误
		if strings.Contains(outputStr, "Access is denied") || strings.Contains(outputStr, "拒绝访问") {
			return fmt.Errorf("permission denied when starting CCR service. Try running this application as administrator. Command output: %s", outputStr)
		}
		// 检查是否是文件未找到错误
		if strings.Contains(outputStr, "not found") || strings.Contains(outputStr, "not recognized") || strings.Contains(outputStr, "无法将") {
			return fmt.Errorf("CCR command not found. Please ensure CCR is properly installed and in your PATH. Command output: %s", outputStr)
		}
		// 增强错误信息，包含更多上下文
		return fmt.Errorf("failed to start CCR service: %v, output: %s, working dir: %s, ccr path: %s", err, outputStr, cmd.Dir, ccrPath)
	}
	
	return nil
}

// StopService stops the CCR service
func (a *App) StopService() error {
	// 查找ccr命令的绝对路径
	ccrPath, err := exec.LookPath("ccr")
	if err != nil {
		// 如果在PATH中找不到，尝试使用完整路径
		ccrPath = "ccr"
	}
	
	// 使用 ccr stop 命令停止服务
	cmd := exec.Command(ccrPath, "stop")
	
	// 设置命令属性
	cmd.Dir = "."
	
	// 设置命令在后台运行，避免创建可见窗口
	cmd.SysProcAttr = getSysProcAttr()
	
	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if outputStr == "" {
			outputStr = "No output"
		}
		// 检查是否是HTTP相关错误
		if strings.Contains(outputStr, "HTTP/1.1 400 Bad Request") {
			return fmt.Errorf("HTTP communication error when stopping CCR service. This may be a Wails internal issue. Command output: %s", outputStr)
		}
		// 检查是否是常见的Windows权限错误
		if strings.Contains(outputStr, "Access is denied") || strings.Contains(outputStr, "拒绝访问") {
			return fmt.Errorf("permission denied when stopping CCR service. Try running this application as administrator. Command output: %s", outputStr)
		}
		// 检查是否是文件未找到错误
		if strings.Contains(outputStr, "not found") || strings.Contains(outputStr, "not recognized") || strings.Contains(outputStr, "无法将") {
			return fmt.Errorf("CCR command not found. Please ensure CCR is properly installed and in your PATH. Command output: %s", outputStr)
		}
		// 增强错误信息，包含更多上下文
		return fmt.Errorf("failed to stop CCR service: %v, output: %s, working dir: %s, ccr path: %s", err, outputStr, cmd.Dir, ccrPath)
	}
	
	return nil
}

// RestartService restarts the CCR service
func (a *App) RestartService() error {
	// 查找ccr命令的绝对路径
	ccrPath, err := exec.LookPath("ccr")
	if err != nil {
		// 如果在PATH中找不到，尝试使用完整路径
		ccrPath = "ccr"
	}
	
	// 使用 ccr restart 命令重启服务
	cmd := exec.Command(ccrPath, "restart")
	
	// 设置命令属性
	cmd.Dir = "."
	
	// 设置命令在后台运行，避免创建可见窗口
	cmd.SysProcAttr = getSysProcAttr()
	
	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if outputStr == "" {
			outputStr = "No output"
		}
		// 检查是否是HTTP相关错误
		if strings.Contains(outputStr, "HTTP/1.1 400 Bad Request") {
			return fmt.Errorf("HTTP communication error when restarting CCR service. This may be a Wails internal issue. Command output: %s", outputStr)
		}
		// 检查是否是常见的Windows权限错误
		if strings.Contains(outputStr, "Access is denied") || strings.Contains(outputStr, "拒绝访问") {
			return fmt.Errorf("permission denied when restarting CCR service. Try running this application as administrator. Command output: %s", outputStr)
		}
		// 检查是否是文件未找到错误
		if strings.Contains(outputStr, "not found") || strings.Contains(outputStr, "not recognized") || strings.Contains(outputStr, "无法将") {
			return fmt.Errorf("CCR command not found. Please ensure CCR is properly installed and in your PATH. Command output: %s", outputStr)
		}
		// 增强错误信息，包含更多上下文
		return fmt.Errorf("failed to restart CCR service: %v, output: %s, working dir: %s, ccr path: %s", err, outputStr, cmd.Dir, ccrPath)
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
	// 尝试从package.json文件读取版本信息
	version, err := a.getCCRViaPackageJSON()
	if err == nil && version != "" {
		return version, nil
	}
	
	// 如果无法从package.json读取，则使用命令行方式
	ccrPath, err := exec.LookPath("ccr")
	if err != nil {
		// 如果在PATH中找不到，尝试使用完整路径
		ccrPath = "ccr"
	}
	
	// 使用 ccr -v 命令获取版本号
	cmd := exec.Command(ccrPath, "-v")
	
	// 设置命令在后台运行，避免创建可见窗口
	cmd.SysProcAttr = getSysProcAttr()
	
	// 执行命令
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get CCR version: %v, output: %s", err, string(output))
	}
	
	// 清理输出，移除换行符
	version = strings.TrimSpace(string(output))
	
	return version, nil
}

// getCCRViaPackageJSON tries to read CCR version from package.json
func (a *App) getCCRViaPackageJSON() (string, error) {
	// 查找ccr命令的路径
	ccrPath, err := exec.LookPath("ccr")
	if err != nil {
		return "", fmt.Errorf("ccr command not found: %v", err)
	}
	
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
		return "", fmt.Errorf("could not determine log path")
	}
	
	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// If log file doesn't exist, return empty string instead of error
		return "", nil
	}
	
	// 限制读取文件大小，避免读取过大的文件
	fileInfo, err := os.Stat(logPath)
	if err != nil {
		return "", fmt.Errorf("failed to get log file info: %v", err)
	}
	
	// 如果文件大于1MB，只读取最后的部分
	const maxFileSize = 20 * 1024 
	if fileInfo.Size() > maxFileSize {
		// 打开文件
		file, err := os.Open(logPath)
		if err != nil {
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
			return "", fmt.Errorf("failed to read log file: %v", err)
		}
		
		// 转换为字符串并限制到最后一行
		content := string(buffer[:n])
		lines := strings.Split(content, "\n")
		
		// 如果有足够多的行，取最后500行
		if len(lines) > 500 {
			lines = lines[len(lines)-500:]
		}
		
		return strings.Join(lines, "\n"), nil
	}
	
	// 文件较小，直接读取
	content, err := os.ReadFile(logPath)
	if err != nil {
		return "", fmt.Errorf("failed to read log file: %v", err)
	}
	
	// Limit to last 500 lines
	logContent := string(content)
	lines := strings.Split(logContent, "\n")
	
	// If there are more than 500 lines, keep only the last 500
	if len(lines) > 500 {
		lines = lines[len(lines)-500:]
	}
	
	return strings.Join(lines, "\n"), nil
}

// ClearLogs clears the CCR log file
func (a *App) ClearLogs() error {
	logPath := a.GetLogPath()
	if logPath == "" {
		return fmt.Errorf("could not determine log path")
	}
	
	// Ensure the directory exists
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}
	
	// Truncate the log file
	return os.WriteFile(logPath, []byte(""), 0644)
}
