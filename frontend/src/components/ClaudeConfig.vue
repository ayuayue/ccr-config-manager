<template>
  <div class="claude-config-container">

    
    <div v-if="status.message" :class="['status', status.type]">
      {{ status.message }}
    </div>
    
    <div class="config-layout">
      <!-- 左侧栏 -->
      <div class="config-sidebar">
        <!-- 基础配置 -->
        <div class="config-section">
          <div class="card-header">
            <h4>基础配置</h4>
          </div>
          <div class="card-content">
            <div class="form-item">
              <label for="apiKey">API Key (可选)</label>
              <input type="password" id="apiKey" v-model="config.APIKEY" placeholder="用于身份验证的密钥">
              <div class="help-text">设置后，客户端请求必须在 Authorization 请求头或 x-api-key 请求头中提供此密钥</div>
            </div>
            
            <div class="form-item">
              <label for="proxyUrl">代理 URL (可选)</label>
              <input type="text" id="proxyUrl" v-model="config.PROXY_URL" placeholder="例如: http://127.0.0.1:7890">
              <div class="help-text">为 API 请求设置代理</div>
            </div>
            
            <div class="form-item">
              <label for="host">主机地址 (可选)</label>
              <input type="text" id="host" v-model="config.HOST" placeholder="例如: 0.0.0.0">
              <div class="help-text">设置服务的主机地址。如果未设置 APIKEY，出于安全考虑，主机地址将强制设置为 127.0.0.1</div>
            </div>
            
            <div class="form-item">
              <label for="apiTimeout">API 超时时间 (毫秒)</label>
              <input type="number" id="apiTimeout" v-model.number="config.API_TIMEOUT_MS" placeholder="例如: 600000">
            </div>
            
            <div class="form-item">
              <label class="checkbox-label">
                <input type="checkbox" v-model="config.LOG"> 启用日志记录
              </label>
              <div class="help-text">启用后，日志文件将位于 $HOME/.claude-code-router.log</div>
            </div>
          </div>
        </div>
        
        <!-- 路由配置 -->
        <div class="config-section">
          <div class="card-header">
            <h4>路由配置</h4>
          </div>
          <div class="card-content">
            <div class="form-item">
              <label for="defaultRoute">默认路由</label>
              <select class="route-select" @change="e => config.Router.default = e.target.value">
                <option value="">请选择或手动输入</option>
                <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value" :selected="config.Router.default === option.value">
                  {{ option.label }}
                </option>
              </select>
              <input type="text" v-model="config.Router.default" placeholder="例如: deepseek,deepseek-chat" class="route-input">
            </div>
            
            <div class="form-item">
              <label for="backgroundRoute">后台任务路由</label>
              <select class="route-select" @change="e => config.Router.background = e.target.value">
                <option value="">请选择或手动输入</option>
                <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value" :selected="config.Router.background === option.value">
                  {{ option.label }}
                </option>
              </select>
              <input type="text" v-model="config.Router.background" placeholder="例如: ollama,qwen2.5-coder:latest" class="route-input">
            </div>
            
            <div class="form-item">
              <label for="thinkRoute">推理任务路由</label>
              <select class="route-select" @change="e => config.Router.think = e.target.value">
                <option value="">请选择或手动输入</option>
                <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value" :selected="config.Router.think === option.value">
                  {{ option.label }}
                </option>
              </select>
              <input type="text" v-model="config.Router.think" placeholder="例如: deepseek,deepseek-reasoner" class="route-input">
            </div>
            
            <div class="form-item">
              <label for="longContextRoute">长上下文路由</label>
              <select class="route-select" @change="e => config.Router.longContext = e.target.value">
                <option value="">请选择或手动输入</option>
                <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value" :selected="config.Router.longContext === option.value">
                  {{ option.label }}
                </option>
              </select>
              <input type="text" v-model="config.Router.longContext" placeholder="例如: openrouter,google/gemini-2.5-pro-preview" class="route-input">
            </div>
            
            <div class="form-item">
              <label for="longContextThreshold">长上下文阈值</label>
              <input type="number" id="longContextThreshold" v-model.number="config.Router.longContextThreshold" placeholder="例如: 60000">
            </div>
            
            <div class="form-item">
              <label for="webSearchRoute">网络搜索路由</label>
              <select class="route-select" @change="e => config.Router.webSearch = e.target.value">
                <option value="">请选择或手动输入</option>
                <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value" :selected="config.Router.webSearch === option.value">
                  {{ option.label }}
                </option>
              </select>
              <input type="text" v-model="config.Router.webSearch" placeholder="例如: gemini,gemini-2.5-flash" class="route-input">
            </div>
          </div>
        </div>
      </div>
      
      <!-- 主内容区 -->
      <div class="config-main">
        <!-- 提供商配置 -->
        <div class="config-section">
          <div class="card-header">
            <h4>提供商配置</h4>
          </div>
          <div class="card-content">
            <div class="providers-container">
              <div class="providers-grid">
                <div v-for="(provider, index) in config.Providers" :key="index" class="provider-card">
                  <div class="provider-header">
                    <h3>提供商 {{ index + 1 }}</h3>
                  </div>
                  <div class="provider-content">
                    <div class="form-item">
                      <label>提供商名称</label>
                      <input type="text" v-model="provider.name" placeholder="例如: openrouter">
                    </div>
                    <div class="form-item">
                      <label>API 基础 URL</label>
                      <input type="text" v-model="provider.api_base_url" placeholder="例如: https://openrouter.ai/api/v1/chat/completions">
                    </div>
                    <div class="form-item">
                      <label>API 密钥</label>
                      <input type="password" v-model="provider.api_key" placeholder="提供商的 API 密钥">
                    </div>
                    <div class="form-item">
                      <label>模型列表 (每行一个)</label>
                      <textarea v-model="provider.modelsText" placeholder="例如:
google/gemini-2.5-pro-preview
anthropic/claude-sonnet-4
anthropic/claude-3.5-sonnet"></textarea>
                    </div>
                    <div class="form-item">
                      <label>转换器 (JSON, 可选)</label>
                      <textarea v-model="provider.transformerText" placeholder='例如:
{
  "use": ["openrouter"]
}'></textarea>
                    </div>
                    <button class="btn-danger" @click="removeProvider(index)">删除提供商</button>
                  </div>
                </div>
              </div>
              <button type="button" @click="addProvider" class="btn-secondary">添加提供商</button>
            </div>
          </div>
        </div>
        
        
        <!-- 操作按钮 -->
        <div class="actions">
          <button @click="saveConfig" class="btn-primary">保存配置</button>
          <button @click="loadConfig" class="btn-secondary">加载配置</button>
          <button @click="clearConfig" class="btn-danger">清空配置</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { LoadConfig, SaveConfig } from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models'

// 状态管理
const status = reactive({
  message: '',
  type: '' // 'success' or 'error'
})

// 配置数据
const config = reactive({
  APIKEY: '',
  PROXY_URL: '',
  HOST: '',
  API_TIMEOUT_MS: 600000,
  LOG: false,
  Providers: [],
  Router: {
    default: '',
    background: '',
    think: '',
    longContext: '',
    longContextThreshold: 60000,
    webSearch: ''
  }
})

// 计算属性：完整配置的 JSON 字符串
const fullConfigJson = computed({
  get() {
    // 构建完整的配置对象
    const fullConfig = {
      ...config,
      Providers: config.Providers.map(provider => {
        const result = {
          name: provider.name,
          api_base_url: provider.api_base_url,
          api_key: provider.api_key,
          models: provider.modelsText ? provider.modelsText.split('\n').filter(model => model.trim() !== '') : []
        }
        
        // 处理转换器
        if (provider.transformerText) {
          try {
            result.transformer = JSON.parse(provider.transformerText)
          } catch (e) {
            console.warn('转换器配置不是有效的 JSON:', e)
          }
        }
        
        return result
      })
    }
    
    // 清理空值
    Object.keys(fullConfig).forEach(key => {
      if (fullConfig[key] === undefined || 
          (Array.isArray(fullConfig[key]) && fullConfig[key].length === 0) ||
          (typeof fullConfig[key] === 'object' && Object.keys(fullConfig[key]).length === 0 && key !== 'Router')) {
        delete fullConfig[key]
      }
    })
    
    return JSON.stringify(fullConfig, null, 2)
  },
  set(value) {
    try {
      const parsed = JSON.parse(value)
      // 更新配置
      Object.keys(parsed).forEach(key => {
        if (key !== 'Providers' && key !== 'Router') {
          config[key] = parsed[key]
        }
      })
      
      // 更新路由配置
      if (parsed.Router) {
        Object.keys(parsed.Router).forEach(key => {
          config.Router[key] = parsed.Router[key]
        })
      }
      
      // 更新提供商配置
      if (parsed.Providers && Array.isArray(parsed.Providers)) {
        config.Providers = parsed.Providers.map(provider => ({
          name: provider.name || '',
          api_base_url: provider.api_base_url || '',
          api_key: provider.api_key || '',
          modelsText: (provider.models || []).join('\n'),
          transformerText: provider.transformer ? JSON.stringify(provider.transformer, null, 2) : ''
        }))
      }
      
      showStatus('配置加载成功', 'success')
    } catch (error) {
      showStatus('JSON 格式错误: ' + error.message, 'error')
    }
  }
})

// 显示状态消息
function showStatus(message, type) {
  status.message = message
  status.type = type
  
  // 3秒后自动隐藏
  setTimeout(() => {
    status.message = ''
    status.type = ''
  }, 3000)
}

// 添加提供商
function addProvider() {
  config.Providers.push({
    name: '',
    api_base_url: '',
    api_key: '',
    modelsText: '',
    transformerText: ''
  })
}

// 获取所有提供商和模型的组合
function getProviderModelOptions() {
  const options = []
  config.Providers.forEach(provider => {
    if (provider.name && provider.modelsText) {
      const models = provider.modelsText.split('\n').filter(model => model.trim() !== '')
      models.forEach(model => {
        options.push({
          label: `${provider.name},${model.trim()}`,
          value: `${provider.name},${model.trim()}`
        })
      })
    }
  })
  return options
}

// 删除提供商
function removeProvider(index) {
  config.Providers.splice(index, 1)
}

// 保存配置
async function saveConfig() {
  try {
    // 构建要保存的配置对象
    const configToSave = {
      APIKEY: config.APIKEY || undefined,
      PROXY_URL: config.PROXY_URL || undefined,
      HOST: config.HOST || undefined,
      API_TIMEOUT_MS: config.API_TIMEOUT_MS || 600000,
      LOG: config.LOG || false,
      Providers: config.Providers.map(provider => {
        const result = {
          name: provider.name,
          api_base_url: provider.api_base_url,
          api_key: provider.api_key,
          models: provider.modelsText ? provider.modelsText.split('\n').filter(model => model.trim() !== '') : []
        }
        
        // 处理转换器
        if (provider.transformerText) {
          try {
            result.transformer = JSON.parse(provider.transformerText)
          } catch (e) {
            console.warn('转换器配置不是有效的 JSON:', e)
          }
        }
        
        return result
      }),
      Router: {
        default: config.Router.default || undefined,
        background: config.Router.background || undefined,
        think: config.Router.think || undefined,
        longContext: config.Router.longContext || undefined,
        longContextThreshold: config.Router.longContextThreshold || 60000,
        webSearch: config.Router.webSearch || undefined
      }
    }
    
    // 调用后端保存配置
    await SaveConfig(configToSave)
    showStatus('配置已保存', 'success')
  } catch (error) {
    showStatus('保存配置时出错: ' + error.message, 'error')
  }
}

// 加载配置
async function loadConfig() {
  try {
    // 从后端加载配置
    const loadedConfig = await LoadConfig()
    
    // 更新本地配置
    config.APIKEY = loadedConfig.APIKEY || ''
    config.PROXY_URL = loadedConfig.PROXY_URL || ''
    config.HOST = loadedConfig.HOST || ''
    config.API_TIMEOUT_MS = loadedConfig.API_TIMEOUT_MS || 600000
    config.LOG = loadedConfig.LOG || false
    
    // 更新路由配置
    if (loadedConfig.Router) {
      config.Router.default = loadedConfig.Router.default || ''
      config.Router.background = loadedConfig.Router.background || ''
      config.Router.think = loadedConfig.Router.think || ''
      config.Router.longContext = loadedConfig.Router.longContext || ''
      config.Router.longContextThreshold = loadedConfig.Router.longContextThreshold || 60000
      config.Router.webSearch = loadedConfig.Router.webSearch || ''
    }
    
    // 更新提供商配置
    if (loadedConfig.Providers && Array.isArray(loadedConfig.Providers)) {
      config.Providers = loadedConfig.Providers.map(provider => ({
        name: provider.name || '',
        api_base_url: provider.api_base_url || '',
        api_key: provider.api_key || '',
        modelsText: (provider.models || []).join('\n'),
        transformerText: provider.transformer ? JSON.stringify(provider.transformer, null, 2) : ''
      }))
    }
    
    showStatus('配置加载成功', 'success')
  } catch (error) {
    showStatus('加载配置时出错: ' + error.message, 'error')
  }
}

// 加载示例配置
function loadExampleConfig() {
  const exampleConfig = {
    "APIKEY": "your-secret-key",
    "PROXY_URL": "http://127.0.0.1:7890",
    "LOG": true,
    "API_TIMEOUT_MS": 600000,
    "Providers": [
      {
        "name": "openrouter",
        "api_base_url": "https://openrouter.ai/api/v1/chat/completions",
        "api_key": "sk-xxx",
        "models": [
          "google/gemini-2.5-pro-preview",
          "anthropic/claude-sonnet-4",
          "anthropic/claude-3.5-sonnet",
          "anthropic/claude-3.7-sonnet:thinking"
        ],
        "transformer": {
          "use": ["openrouter"]
        }
      },
      {
        "name": "deepseek",
        "api_base_url": "https://api.deepseek.com/chat/completions",
        "api_key": "sk-xxx",
        "models": ["deepseek-chat", "deepseek-reasoner"],
        "transformer": {
          "use": ["deepseek"],
          "deepseek-chat": {
            "use": ["tooluse"]
          }
        }
      }
    ],
    "Router": {
      "default": "deepseek,deepseek-chat",
      "background": "ollama,qwen2.5-coder:latest",
      "think": "deepseek,deepseek-reasoner",
      "longContext": "openrouter,google/gemini-2.5-pro-preview",
      "longContextThreshold": 60000,
      "webSearch": "gemini,gemini-2.5-flash"
    }
  }
  
  // 更新配置
  Object.keys(exampleConfig).forEach(key => {
    if (key !== 'Providers' && key !== 'Router') {
      config[key] = exampleConfig[key]
    }
  })
  
  // 更新路由配置
  if (exampleConfig.Router) {
    Object.keys(exampleConfig.Router).forEach(key => {
      config.Router[key] = exampleConfig.Router[key]
    })
  }
  
  // 更新提供商配置
  if (exampleConfig.Providers && Array.isArray(exampleConfig.Providers)) {
    config.Providers = exampleConfig.Providers.map(provider => ({
      name: provider.name || '',
      api_base_url: provider.api_base_url || '',
      api_key: provider.api_key || '',
      modelsText: (provider.models || []).join('\n'),
      transformerText: provider.transformer ? JSON.stringify(provider.transformer, null, 2) : ''
    }))
  }
  
  showStatus('示例配置已加载', 'success')
}

// 清空配置
function clearConfig() {
  if (confirm('确定要清空所有配置吗？')) {
    // 清空基础配置
    config.APIKEY = ''
    config.PROXY_URL = ''
    config.HOST = ''
    config.API_TIMEOUT_MS = 600000
    config.LOG = false
    
    // 清空路由配置
    config.Router = {
      default: '',
      background: '',
      think: '',
      longContext: '',
      longContextThreshold: 60000,
      webSearch: ''
    }
    
    // 清空提供商
    config.Providers = []
    
    showStatus('配置已清空', 'success')
  }
}

// 页面加载完成后初始化
onMounted(() => {
  // 页面加载时自动加载配置
  loadConfig()
})
</script>

<style scoped>
.claude-config-container {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4edf9 100%);
  color: #333;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  overflow: auto;
  padding: 15px;
  box-sizing: border-box;
}

.header {
  background-color: #2c3e50;
  color: white;
  padding: 5px;
  box-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
}

.header h1 {
  text-align: center;
  margin: 0;
  font-size: 1em;
  font-weight: 300;
}

.config-layout {
  display: flex;
  height: calc(100vh - 100px);
  padding: 10px;
  gap: 15px;
}

.config-sidebar {
  flex: 0 0 350px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.config-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.config-section {
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 100%;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.config-section:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
}

.card-header {
  background: linear-gradient(90deg, #3498db 0%, #2c80b9 100%);
  color: white;
  margin: 0;
  padding: 15px 20px;
  font-size: 1.2em;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.card-content {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.form-item {
  margin-bottom: 20px;
  padding: 15px;
  border-radius: 8px;
  background-color: #ffffff;
  transition: all 0.2s ease;
}

.form-item:hover {
  background-color: #f8f9fa;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #2c3e50;
  font-size: 0.95em;
}

.checkbox-label {
  display: flex;
  align-items: center;
  font-weight: normal;
  cursor: pointer;
  font-size: 0.95em;
  padding: 10px;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.checkbox-label:hover {
  background-color: #f8f9fa;
}

.checkbox-label input[type="checkbox"] {
  margin-right: 10px;
  transform: scale(1.2);
}

input, textarea, select {
  width: 100%;
  padding: 12px 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;
  font-size: 0.95em;
  font-family: inherit;
  transition: all 0.2s ease;
}

input:focus, textarea:focus, select:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
}

textarea {
  min-height: 100px;
  font-family: monospace;
  resize: vertical;
}

.route-select {
  width: 100%;
  padding: 12px 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;
  font-size: 0.95em;
  font-family: inherit;
  margin-bottom: 10px;
  background-color: white;
  transition: all 0.2s ease;
}

.route-select:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
}

.route-input {
  width: 100%;
  padding: 12px 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;
  font-size: 0.95em;
  font-family: inherit;
  transition: all 0.2s ease;
}

.route-input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
}

.full-config {
  width: 100%;
  height: 200px;
  padding: 5px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
  font-family: 'Courier New', monospace;
  background-color: #fdfdfd;
  resize: vertical;
  font-size: 0.85em;
  line-height: 1.3;
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.1);
}

button {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.95em;
  font-weight: 500;
  margin-right: 8px;
  margin-bottom: 8px;
  transition: all 0.2s ease;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

button:hover {
  background: linear-gradient(135deg, #45a049 0%, #3d8b40 100%);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

button:active {
  transform: translateY(0);
  box-shadow: 0 2px 3px rgba(0, 0, 0, 0.1);
}

.btn-secondary {
  background: linear-gradient(135deg, #2196F3 0%, #1976D2 100%);
}

.btn-secondary:hover {
  background: linear-gradient(135deg, #1976D2 0%, #1565C0 100%);
}

.btn-danger {
  background: linear-gradient(135deg, #f44336 0%, #d32f2f 100%);
}

.btn-danger:hover {
  background: linear-gradient(135deg, #d32f2f 0%, #c62828 100%);
}

.providers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.provider-card {
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 10px;
  padding: 20px;
  border-left: 4px solid #3498db;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.provider-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
}

.provider-header h3 {
  color: #2c3e50;
  margin: 0 0 15px 0;
  font-size: 1.1em;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.provider-content {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.actions {
  text-align: center;
  padding: 25px;
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
}

.status {
  padding: 15px 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  text-align: center;
  font-size: 0.95em;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.status.success {
  background: linear-gradient(135deg, #d4edda 0%, #c3e6cb 100%);
  color: #155724;
  border: 1px solid #c3e6cb;
}

.status.error {
  background: linear-gradient(135deg, #f8d7da 0%, #f5c6cb 100%);
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.help-text {
  font-size: 12px;
  color: #7f8c8d;
  margin-top: 3px;
  font-style: italic;
}

@media (max-width: 1200px) {
  .config-layout {
    flex-direction: column;
    height: auto;
  }
  
  .config-sidebar {
    flex: none;
  }
  
  .providers-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  }
  
  .config-sidebar,
  .config-main {
    gap: 15px;
  }
}

@media (max-width: 768px) {
  .claude-config-container {
    padding: 10px;
  }
  
  .config-layout {
    padding: 10px;
    gap: 15px;
  }
  
  .config-sidebar,
  .config-main {
    gap: 15px;
  }
  
  .providers-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .provider-card {
    padding: 15px;
  }
  
  .card-header {
    padding: 12px 15px;
    font-size: 1.1em;
  }
  
  .card-content {
    padding: 15px;
  }
  
  .form-item {
    margin-bottom: 15px;
    padding: 12px;
  }
  
  input, textarea, select, .route-select, .route-input {
    padding: 10px 12px;
  }
  
  button {
    padding: 8px 12px;
    font-size: 0.9em;
  }
  
  .header {
    padding: 10px;
  }
  
  .header h1 {
    font-size: 1.2em;
  }
  
  .actions {
    padding: 20px;
  }
}

@media (max-width: 480px) {
  .claude-config-container {
    padding: 5px;
  }
  
  .config-layout {
    padding: 5px;
    gap: 10px;
  }
  
  .config-sidebar,
  .config-main {
    gap: 10px;
  }
  
  .config-section {
    border-radius: 8px;
  }
  
  .card-header {
    padding: 10px 12px;
    font-size: 1em;
  }
  
  .card-content {
    padding: 12px;
  }
  
  .form-item {
    margin-bottom: 12px;
    padding: 10px;
  }
  
  label {
    font-size: 0.9em;
  }
  
  input, textarea, select, .route-select, .route-input {
    padding: 8px 10px;
    font-size: 0.9em;
  }
  
  button {
    padding: 6px 10px;
    font-size: 0.85em;
    margin-right: 5px;
    margin-bottom: 5px;
  }
  
  .providers-grid {
    gap: 10px;
  }
  
  .provider-card {
    padding: 12px;
    border-radius: 8px;
  }
  
  .provider-header h3 {
    font-size: 1em;
    margin: 0 0 10px 0;
  }
  
  .actions {
    padding: 15px;
  }
  
  .status {
    padding: 12px 15px;
    font-size: 0.9em;
  }
}
</style>