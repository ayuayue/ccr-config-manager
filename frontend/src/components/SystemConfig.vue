<template>
  <div class="system-config-container">

    <div v-if="loading" class="loading">
      正在加载配置...
    </div>
    
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    
    <div v-else-if="config" class="config-display">
      <div class="config-summary">
        <div class="summary-card">
          <h2>核心配置概览</h2>
          <div class="summary-content">
            <div class="summary-item">
              <span class="label">API Key:</span>
              <span class="value" v-if="config.APIKEY">
                <span :class="{ masked: !showApiKeys['main'] }">{{ getDisplayApiKey(config.APIKEY, !showApiKeys['main']) }}</span>
                <button class="toggle-btn" @click="toggleApiKeyVisibility('main')">
                  {{ showApiKeys['main'] ? '隐藏' : '显示' }}
                </button>
              </span>
              <span class="value" v-else>未设置</span>
            </div>
            
            <div class="summary-item">
              <span class="label">代理设置:</span>
              <span class="value">{{ config.PROXY_URL || '未设置' }}</span>
            </div>
            
            <div class="summary-item">
              <span class="label">主机地址:</span>
              <span class="value">{{ config.HOST || '默认 (127.0.0.1)' }}</span>
            </div>
            
            <div class="summary-item">
              <span class="label">超时时间:</span>
              <span class="value">{{ config.API_TIMEOUT_MS || 600000 }} 毫秒</span>
            </div>
            
            <div class="summary-item">
              <span class="label">日志记录:</span>
              <span class="value">{{ config.LOG ? '已启用' : '已禁用' }}</span>
            </div>
          </div>
        </div>
        
        <div class="summary-card">
          <h2>路由配置</h2>
          <div class="summary-content">
            <div class="summary-item">
              <span class="label">默认路由:</span>
              <span class="value">{{ config.Router?.default || '未设置' }}</span>
            </div>
            
            <div class="summary-item">
              <span class="label">后台任务:</span>
              <span class="value">{{ config.Router?.background || '未设置' }}</span>
            </div>
            
            <div class="summary-item">
              <span class="label">推理任务:</span>
              <span class="value">{{ config.Router?.think || '未设置' }}</span>
            </div>
            
            <div class="summary-item">
              <span class="label">长上下文:</span>
              <span class="value">{{ config.Router?.longContext || '未设置' }}</span>
            </div>
            
            <div class="summary-item">
              <span class="label">网络搜索:</span>
              <span class="value">{{ config.Router?.webSearch || '未设置' }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="providers-section">
        <h2>提供商配置</h2>
        <div class="providers-grid">
          <div v-for="(provider, index) in config.Providers" :key="index" class="provider-card">
            <div class="provider-header">
              <h3>{{ provider.name }}</h3>
            </div>
            <div class="provider-content">
              <div class="provider-summary">
                <div class="summary-item">
                  <span class="label">API URL:</span>
                  <span class="value">{{ provider.api_base_url }}</span>
                </div>
                
                <div class="summary-item">
                  <span class="label">API Key:</span>
                  <span class="value">
                    <span :class="{ masked: !showApiKeys['provider-' + index] }">{{ getDisplayApiKey(provider.api_key, !showApiKeys['provider-' + index]) }}</span>
                    <button class="toggle-btn" @click="toggleApiKeyVisibility('provider-' + index)">
                      {{ showApiKeys['provider-' + index] ? '隐藏' : '显示' }}
                    </button>
                  </span>
                </div>
                
                <div class="summary-item">
                  <span class="label">模型数量:</span>
                  <span class="value">{{ provider.models?.length || 0 }}</span>
                </div>
              </div>
              
              <div class="models-list" v-if="provider.models && provider.models.length > 0">
                <div class="models-header">支持的模型:</div>
                <div class="models-tags">
                  <span v-for="(model, modelIndex) in provider.models" :key="modelIndex" class="model-tag">
                    {{ model }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="full-config-section">
        <h2>完整配置 (JSON)</h2>
        <textarea readonly class="full-config" :value="fullConfigJson"></textarea>
      </div>
    </div>
    
    <div v-else class="no-config">
      <p>未找到配置文件</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { LoadConfig } from '../../wailsjs/go/main/App'

// 响应式数据
const config = ref(null)
const loading = ref(true)
const error = ref('')

// 控制 API Key 显示状态
const showApiKeys = ref({})

// 计算属性：完整配置的 JSON 字符串
const fullConfigJson = computed(() => {
  if (!config.value) return ''
  return JSON.stringify(config.value, null, 2)
})

// 切换 API Key 显示状态
function toggleApiKeyVisibility(key) {
  showApiKeys.value[key] = !showApiKeys.value[key]
}

// 获取显示的 API Key 值
function getDisplayApiKey(key, isMasked = true) {
  if (!key) return ''
  if (!isMasked || showApiKeys.value[key]) {
    return key
  }
  return '•'.repeat(Math.min(key.length, 20))
}

// 加载配置
async function loadSystemConfig() {
  try {
    loading.value = true
    error.value = ''
    
    // 从后端加载配置
    const loadedConfig = await LoadConfig()
    config.value = loadedConfig
    
    // 如果配置为空对象，设置为 null
    if (Object.keys(loadedConfig).length === 0) {
      config.value = null
    }
  } catch (err) {
    error.value = '加载配置时出错: ' + err.message
  } finally {
    loading.value = false
  }
}

// 页面加载完成后初始化
onMounted(() => {
  loadSystemConfig()
})
</script>

<style scoped>
.system-config-container {
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
  border-radius: 6px;
  box-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.header h1 {
  text-align: center;
  margin: 0;
  font-size: 1.5em;
  font-weight: 500;
}

.config-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 15px;
  margin-bottom: 20px;
}

.summary-card {
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.summary-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.12);
}

.summary-card h2 {
  background: linear-gradient(90deg, #3498db 0%, #2c80b9 100%);
  color: white;
  margin: 0;
  padding: 15px 20px;
  font-size: 1.2em;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.summary-content {
  padding: 15px 20px;
}

.summary-item {
  display: flex;
  padding: 12px 0;
  border-bottom: 1px solid #eee;
  transition: background-color 0.2s ease;
}

.summary-item:hover {
  background-color: #f8f9fa;
  border-radius: 4px;
  padding-left: 10px;
  padding-right: 10px;
}

.summary-item:last-child {
  border-bottom: none;
}

.summary-item .label {
  font-weight: 600;
  color: #2c3e50;
  min-width: 120px;
  margin-right: 15px;
  font-size: 0.95em;
  display: flex;
  align-items: center;
}

.summary-item .value {
  flex: 1;
  color: #34495e;
  font-size: 0.95em;
  word-break: break-word;
  display: flex;
  align-items: center;
}

.masked {
  color: #95a5a6;
  font-family: monospace;
}

.toggle-btn {
  margin-left: 8px;
  background-color: #3498db;
  color: white;
  border: none;
  padding: 4px 8px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 0.8em;
  transition: all 0.1s ease;
}

.toggle-btn:hover {
  background-color: #2980b9;
}

.toggle-btn:active {
  transform: scale(0.95);
}

.providers-section h2 {
  background: linear-gradient(90deg, #3498db 0%, #2c80b9 100%);
  color: white;
  margin: 0 0 15px 0;
  padding: 15px 20px;
  font-size: 1.2em;
  font-weight: 600;
  border-radius: 10px 10px 0 0;
  letter-spacing: 0.5px;
}

.providers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.provider-card {
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.provider-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.12);
}

.provider-header {
  background: linear-gradient(90deg, #2c3e50 0%, #1a2530 100%);
  color: white;
  padding: 15px 20px;
}

.provider-header h3 {
  margin: 0;
  font-size: 1.1em;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.provider-content {
  padding: 15px 20px;
}

.provider-summary {
  margin-bottom: 5px;
}

.provider-summary .summary-item {
  padding: 8px 0;
}

.models-list {
  border-top: 1px solid #eee;
  padding-top: 10px;
}

.models-header {
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
  font-size: 0.9em;
}

.models-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.model-tag {
  background: linear-gradient(135deg, #e1f0fa 0%, #d1e7f5 100%);
  color: #2980b9;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.85em;
  font-weight: 500;
  margin: 4px;
  display: inline-block;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
}

.model-tag:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.full-config-section h2 {
  background: linear-gradient(90deg, #3498db 0%, #2c80b9 100%);
  color: white;
  margin: 0 0 15px 0;
  padding: 15px 20px;
  font-size: 1.2em;
  font-weight: 600;
  border-radius: 10px 10px 0 0;
  letter-spacing: 0.5px;
}

.full-config {
  width: 100%;
  height: 550px;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 0 0 10px 10px;
  box-sizing: border-box;
  font-family: 'Courier New', monospace;
  background: linear-gradient(180deg, #fdfdfd 0%, #f8f9fa 100%);
  resize: vertical;
  font-size: 0.9em;
  line-height: 1.4;
  box-shadow: inset 0 2px 5px rgba(0, 0, 0, 0.1);
  border: 1px solid #e1e8ed;
}

.loading, .error, .no-config {
  text-align: center;
  padding: 50px;
  font-size: 18px;
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  margin: 20px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.loading {
  color: #3498db;
}

.loading::before {
  content: "🔄";
  margin-right: 10px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.error {
  color: #e74c3c;
  background: linear-gradient(180deg, #fdf2f2 0%, #fceaea 100%);
  border: 1px solid #fadbd8;
}

.error::before {
  content: "⚠️";
  margin-right: 10px;
}

.no-config {
  color: #7f8c8d;
  background: linear-gradient(180deg, #f8f9fa 0%, #eef2f7 100%);
}

.no-config::before {
  content: "📋";
  margin-right: 10px;
}

@media (max-width: 768px) {
  .system-config-container {
    padding: 10px;
  }
  
  .config-summary {
    grid-template-columns: 1fr;
    gap: 10px;
  }
  
  .summary-card {
    border-radius: 8px;
  }
  
  .summary-card h2 {
    padding: 12px 15px;
    font-size: 1.1em;
  }
  
  .summary-content {
    padding: 12px 15px;
  }
  
  .summary-item {
    flex-direction: column;
    gap: 4px;
    padding: 10px 0;
  }
  
  .summary-item .label {
    min-width: auto;
    margin-right: 0;
    font-weight: 600;
  }
  
  .providers-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .provider-card {
    border-radius: 8px;
  }
  
  .provider-header {
    padding: 12px 15px;
  }
  
  .provider-header h3 {
    font-size: 1em;
  }
  
  .provider-content {
    padding: 12px 15px;
  }
  
  .full-config {
    height: 400px;
  }
  
  .loading, .error, .no-config {
    padding: 30px;
    font-size: 16px;
    margin: 10px;
  }
}

@media (max-width: 480px) {
  .system-config-container {
    padding: 5px;
  }
  
  .summary-card h2 {
    padding: 10px 12px;
    font-size: 1em;
  }
  
  .summary-content {
    padding: 10px 12px;
  }
  
  .summary-item {
    padding: 8px 0;
  }
  
  .summary-item .label,
  .summary-item .value {
    font-size: 0.9em;
  }
  
  .provider-header {
    padding: 10px 12px;
  }
  
  .provider-content {
    padding: 10px 12px;
  }
  
  .model-tag {
    padding: 4px 8px;
    font-size: 0.8em;
    margin: 2px;
  }
  
  .full-config {
    height: 300px;
    padding: 10px;
    font-size: 0.8em;
  }
  
  .loading, .error, .no-config {
    padding: 20px;
    font-size: 14px;
    margin: 5px;
  }
}
</style>