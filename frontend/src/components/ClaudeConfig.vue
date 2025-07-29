<template>
  <div class="claude-config-container">


    <div v-if="status.message" :class="['status', status.type]">
      {{ status.message }}
    </div>

    <!-- 标签页导航 -->
    <div class="tabs">
      <div class="tab-group">
        <button v-for="tab in tabs" :key="tab.id" :class="['tab', { active: activeTab === tab.id }]"
          @click="activeTab = tab.id">
          {{ tab.name }}
        </button>
      </div>
      <div class="tab-actions">
        <button @click="loadConfig" class="btn-secondary">刷新配置</button>
        <button @click="saveConfig" class="btn-primary">保存配置</button>
      </div>
    </div>

    <!-- 标签页内容 -->
    <div class="tab-content">
      <!-- 基础配置和路由配置标签页 -->
      <div v-show="activeTab === 'basic'" class="tab-pane">
        <div class="config-grid">
          <!-- 基础配置在左侧 -->
          <div class="config-group">
            <h3>基础配置</h3>
            <div class="form-grid">
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

          <!-- 路由配置在右侧 -->
          <div class="config-group">
            <div class="form-grid">
              <div class="form-item">
                <label for="longContextThreshold">长上下文阈值</label>
                <input type="number" id="longContextThreshold" v-model.number="config.Router.longContextThreshold"
                  placeholder="例如: 60000">
              </div>
              <div class="form-item">
                <label for="defaultRoute">默认路由</label>
                <select class="route-select" v-model="config.Router.default">
                  <option value="">请选择</option>
                  <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </div>

              <div class="form-item">
                <label for="longContextRoute">长上下文路由</label>
                <select class="route-select" v-model="config.Router.longContext">
                  <option value="">请选择</option>
                  <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </div>

              <div class="form-item">
                <label for="backgroundRoute">后台任务路由</label>
                <select class="route-select" v-model="config.Router.background">
                  <option value="">请选择</option>
                  <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </div>

              <div class="form-item">
                <label for="thinkRoute">推理任务路由</label>
                <select class="route-select" v-model="config.Router.think">
                  <option value="">请选择</option>
                  <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </div>

              <div class="form-item">
                <label for="webSearchRoute">网络搜索路由</label>
                <select class="route-select" v-model="config.Router.webSearch">
                  <option value="">请选择</option>
                  <option v-for="option in getProviderModelOptions()" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </div>
            </div>


          </div>
        </div>
      </div>

      <!-- 提供商配置标签页 -->
      <div v-show="activeTab === 'providers'" class="tab-pane">
        <div class="providers-section">
          <div class="providers-header">
            <h3>提供商配置 ({{ config.Providers?.length || 0 }} 个)</h3>
          </div>
          <div class="providers-list">
            <div v-for="(provider, index) in config.Providers" :key="index" class="provider-item"
              :class="{ expanded: expandedProviders[index] }">
              <div class="provider-summary" @click="toggleProvider(index)">
                <div class="provider-name">提供商 {{ index + 1 }}: {{ provider.name || '未命名' }}</div>
                <div class="expand-icon">{{ expandedProviders[index] ? '▲' : '▼' }}</div>
              </div>
              <div v-show="expandedProviders[index]" class="provider-details">
                <div class="form-grid">
                  <div class="form-item">
                    <label>提供商名称</label>
                    <input type="text" v-model="provider.name" placeholder="例如: openrouter">
                  </div>
                  <div class="form-item">
                    <label>API 基础 URL</label>
                    <input type="text" v-model="provider.api_base_url"
                      placeholder="例如: https://openrouter.ai/api/v1/chat/completions">
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
                  <div class="form-item">
                    <button class="btn-danger" @click.stop="removeProvider(index)">删除提供商</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <button type="button" @click="addProvider" class="btn-secondary">添加提供商</button>
        </div>
      </div>

      <!-- 完整配置标签页 -->
      <div v-show="activeTab === 'full'" class="tab-pane">
        <div class="full-config-section">
          <h3>完整配置 (JSON)</h3>
          <textarea v-model="fullConfigJson" class="full-config"></textarea>
          <div class="actions">
            <button @click="saveConfig" class="btn-primary">保存配置</button>
          </div>
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

// 标签页相关
const activeTab = ref('basic')
const tabs = ref([
  { id: 'basic', name: '基础配置' },
  { id: 'providers', name: '提供商配置' },
  { id: 'full', name: '完整配置' }
])

// 控制提供商展开状态
const expandedProviders = ref({})

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

// 切换提供商展开状态
function toggleProvider(index) {
  expandedProviders.value[index] = !expandedProviders.value[index]
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
  padding: 12px;
  box-sizing: border-box;
}

/* 标签页导航 */
.tabs {
  display: flex;
  gap: 3px;
  margin-bottom: 15px;
  border-bottom: 2px solid #e1e8ed;
  max-width: 1000px;
  align-items: center;
}

.tab-group {
  display: flex;
  flex: 1;
}

.tab {
  padding: 8px 16px;
  background: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  cursor: pointer;
  font-size: 0.9em;
  font-weight: 500;
  color: #7f8c8d;
  transition: all 0.3s ease;
  position: relative;
  flex-shrink: 0;
}

.tab-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  padding: 0 10px;
  margin-left: auto;
}

.tab:hover {
  color: #3498db;
  background-color: rgba(52, 152, 219, 0.1);
}

.tab.active {
  color: #3498db;
  border-bottom: 3px solid #3498db;
  background-color: rgba(52, 152, 219, 0.05);
}

/* 标签页内容 */
.tab-content {
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.08);
  padding: 15px;
  min-height: 300px;
  max-width: 1200px;
  margin: 0 auto;
}

.tab-pane {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 配置网格 */
.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 20px;
}

.config-group {
  background: #fafcff;
  border-radius: 8px;
  padding: 15px;
  border: 1px solid #e1e8ed;
}

.config-group h3 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 1.1em;
  font-weight: 600;
  padding-bottom: 8px;
  border-bottom: 1px solid #eee;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
}

.form-item {
  margin-bottom: 0;
  padding: 10px;
  border-radius: 6px;
  background-color: #ffffff;
  transition: all 0.2s ease;
  border: 1px solid #f0f0f0;
}

.form-item:hover {
  background-color: #f8f9fa;
  border-color: #e1e8ed;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.route-select {
  margin-bottom: 8px;
}

label {
  display: block;
  margin-bottom: 6px;
  font-weight: 600;
  color: #2c3e50;
  font-size: 0.9em;
}

.checkbox-label {
  display: flex;
  align-items: center;
  font-weight: normal;
  cursor: pointer;
  font-size: 0.9em;
  padding: 8px;
  border-radius: 4px;
  transition: background-color 0.2s ease;
}

.checkbox-label:hover {
  background-color: #f8f9fa;
}

.checkbox-label input[type="checkbox"] {
  margin-right: 8px;
  transform: scale(1.1);
}

input,
textarea,
select {
  max-width: 120px;
  width: 100%;
  padding: 8px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
  font-size: 0.9em;
  font-family: inherit;
  transition: all 0.2s ease;
}

input:focus,
textarea:focus,
select:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.1);
}

textarea {
  min-height: 80px;
  font-family: monospace;
  resize: vertical;
  max-width: 100%;
}

.route-select {
  max-width: 1200px;
  width: 100%;
  padding: 8px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
  font-size: 0.9em;
  font-family: inherit;
  margin-bottom: 8px;
  background-color: white;
  transition: all 0.2s ease;
}

.route-select:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.1);
}

.route-input {
  max-width: 500px;
  width: 100%;
  padding: 8px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
  font-size: 0.9em;
  font-family: inherit;
  transition: all 0.2s ease;
}

.route-input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.1);
}

/* 提供商配置 */
.providers-section {
  padding: 10px 0;
  max-width: 1200px;
}

.providers-header h3 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 1.1em;
  font-weight: 600;
}

.providers-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

.provider-item {
  border: 1px solid #e1e8ed;
  border-radius: 6px;
  overflow: hidden;
  background: #fafcff;
  transition: all 0.3s ease;
}

.provider-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

.provider-summary {
  display: flex;
  align-items: center;
  padding: 12px 15px;
  cursor: pointer;
  background: linear-gradient(90deg, #2c3e50 0%, #1a2530 100%);
  color: white;
}

.provider-name {
  flex: 1;
  font-weight: 600;
  font-size: 0.9em;
}

.expand-icon {
  font-size: 0.8em;
  transition: transform 0.3s ease;
}

.provider-item.expanded .expand-icon {
  transform: rotate(180deg);
}

.provider-details {
  padding: 12px 15px;
  background: white;
  border-top: 1px solid #e1e8ed;
}

/* 完整配置区域 */
.full-config-section {
  max-width: 1200px;
}

.full-config-section h3 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 1.1em;
  font-weight: 600;
}

.full-config {
  width: 1000;
  height: 800px;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;
  font-family: 'Courier New', monospace;
  background: linear-gradient(180deg, #fdfdfd 0%, #f8f9fa 100%);
  resize: vertical;
  font-size: 0.85em;
  line-height: 1.4;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e1e8ed;
  margin-bottom: 15px;
}

.actions {
  text-align: center;
  padding: 15px 0;
}

button {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9em;
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

.status {
  padding: 15px 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  text-align: center;
  font-size: 0.95em;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  animation: fadeIn 0.3s ease;
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
  font-size: 11px;
  color: #7f8c8d;
  margin-top: 3px;
  font-style: italic;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .claude-config-container {
    padding: 8px;
  }

  .tabs {
    margin-bottom: 12px;
    max-width: 100%;
  }

  .tab {
    padding: 6px 12px;
    font-size: 0.85em;
  }

  .tab-content {
    padding: 12px;
    min-height: 250px;
  }

  .config-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .config-group {
    padding: 10px;
  }

  .config-group h3 {
    font-size: 1em;
    margin-bottom: 10px;
  }

  .form-grid {
    gap: 10px;
  }

  .form-item {
    padding: 8px;
  }

  input,
  textarea,
  select,
  .route-select,
  .route-input {
    max-width: 100%;
  }

  .provider-summary {
    padding: 10px 12px;
  }

  .provider-name {
    font-size: 0.85em;
  }

  .provider-details {
    padding: 10px 12px;
  }

  .providers-section {
    max-width: 100%;
  }

  .full-config-section {
    max-width: 100%;
  }

  .full-config {
    max-width: 100%;
    height: 250px;
    padding: 10px;
  }

  .status {
    padding: 10px 12px;
    font-size: 0.85em;
  }
}

@media (max-width: 480px) {
  .claude-config-container {
    padding: 5px;
  }

  .tabs {
    flex-wrap: wrap;
    margin-bottom: 10px;
  }

  .tab {
    padding: 5px 10px;
    font-size: 0.8em;
  }

  .tab-content {
    padding: 10px;
    min-height: 200px;
  }

  .config-group {
    padding: 8px;
  }

  .config-group h3 {
    font-size: 0.95em;
    margin-bottom: 8px;
  }

  .form-grid {
    gap: 8px;
  }

  .form-item {
    padding: 6px;
  }

  label {
    font-size: 0.8em;
  }

  input,
  textarea,
  select,
  .route-select,
  .route-input {
    padding: 6px 8px;
    font-size: 0.8em;
  }

  .provider-summary {
    padding: 8px 10px;
  }

  .provider-name {
    font-size: 0.8em;
  }

  .provider-details {
    padding: 8px 10px;
  }

  .full-config {
    height: 200px;
    padding: 8px;
    font-size: 0.75em;
  }

  button {
    padding: 5px 8px;
    font-size: 0.75em;
    margin-right: 4px;
    margin-bottom: 4px;
  }

  .status {
    padding: 8px 10px;
    font-size: 0.8em;
  }
}
</style>