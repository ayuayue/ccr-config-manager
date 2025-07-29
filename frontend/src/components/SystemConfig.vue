<template>
  <div class="system-config-container">

    <div v-if="loading" class="loading">
      正在加载配置...
    </div>
    
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    
    <div v-else-if="config" class="config-display">
      <!-- 标签页导航 -->
      <div class="tabs">
        <button 
          v-for="tab in tabs" 
          :key="tab.id" 
          :class="['tab', { active: activeTab === tab.id }]"
          @click="activeTab = tab.id"
        >
          {{ tab.name }}
        </button>
      </div>
      
      <!-- 标签页内容 -->
      <div class="tab-content">
        <!-- 核心配置标签页 -->
        <div v-show="activeTab === 'core'" class="tab-pane">
          <div class="config-grid">
            <div class="config-group">
              <h3>核心配置</h3>
              <div class="config-table">
                <div class="config-row">
                  <span class="label">API Key:</span>
                  <span class="value" v-if="config.APIKEY">
                    <span :class="{ masked: !showApiKeys['main'] }">{{ getDisplayApiKey(config.APIKEY, !showApiKeys['main']) }}</span>
                    <button class="toggle-btn" @click="toggleApiKeyVisibility('main')">
                      {{ showApiKeys['main'] ? '隐藏' : '显示' }}
                    </button>
                  </span>
                  <span class="value" v-else>未设置</span>
                </div>
                
                <div class="config-row">
                  <span class="label">代理设置:</span>
                  <span class="value">{{ config.PROXY_URL || '未设置' }}</span>
                </div>
                
                <div class="config-row">
                  <span class="label">主机地址:</span>
                  <span class="value">{{ config.HOST || '默认 (127.0.0.1)' }}</span>
                </div>
                
                <div class="config-row">
                  <span class="label">超时时间:</span>
                  <span class="value">{{ config.API_TIMEOUT_MS || 600000 }} 毫秒</span>
                </div>
                
                <div class="config-row">
                  <span class="label">日志记录:</span>
                  <span class="value">{{ config.LOG ? '已启用' : '已禁用' }}</span>
                </div>
              </div>
            </div>
            
            <div class="config-group">
              <h3>路由配置</h3>
              <div class="config-table">
                <div class="config-row">
                  <span class="label">默认路由:</span>
                  <span class="value">{{ config.Router?.default || '未设置' }}</span>
                </div>
                
                <div class="config-row">
                  <span class="label">后台任务:</span>
                  <span class="value">{{ config.Router?.background || '未设置' }}</span>
                </div>
                
                <div class="config-row">
                  <span class="label">推理任务:</span>
                  <span class="value">{{ config.Router?.think || '未设置' }}</span>
                </div>
                
                <div class="config-row">
                  <span class="label">长上下文:</span>
                  <span class="value">{{ config.Router?.longContext || '未设置' }}</span>
                </div>
                
                <div class="config-row">
                  <span class="label">网络搜索:</span>
                  <span class="value">{{ config.Router?.webSearch || '未设置' }}</span>
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
              <div 
                v-for="(provider, index) in config.Providers" 
                :key="index" 
                class="provider-item"
                :class="{ expanded: expandedProviders[index] }"
              >
                <div class="provider-summary" @click="toggleProvider(index)">
                  <div class="provider-name">{{ provider.name }}</div>
                  <div class="provider-models-count">{{ provider.models?.length || 0 }} 个模型</div>
                  <div class="expand-icon">{{ expandedProviders[index] ? '▲' : '▼' }}</div>
                </div>
                <div v-show="expandedProviders[index]" class="provider-details">
                  <div class="detail-row">
                    <span class="label">API URL:</span>
                    <span class="value">{{ provider.api_base_url }}</span>
                  </div>
                  
                  <div class="detail-row">
                    <span class="label">API Key:</span>
                    <span class="value">
                      <span :class="{ masked: !showApiKeys['provider-' + index] }">{{ getDisplayApiKey(provider.api_key, !showApiKeys['provider-' + index]) }}</span>
                      <button class="toggle-btn small" @click.stop="toggleApiKeyVisibility('provider-' + index)">
                        {{ showApiKeys['provider-' + index] ? '隐藏' : '显示' }}
                      </button>
                    </span>
                  </div>
                  
                  <div class="models-section" v-if="provider.models && provider.models.length > 0">
                    <div class="models-label">支持的模型:</div>
                    <div class="models-grid">
                      <span v-for="(model, modelIndex) in provider.models" :key="modelIndex" class="model-tag">
                        {{ model }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 完整配置标签页 -->
        <div v-show="activeTab === 'full'" class="tab-pane">
          <div class="full-config-section">
            <h3>完整配置 (JSON)</h3>
            <textarea readonly class="full-config" :value="fullConfigJson"></textarea>
          </div>
        </div>
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

// 标签页相关
const activeTab = ref('core')
const tabs = ref([
  { id: 'core', name: '核心配置' },
  { id: 'providers', name: '提供商配置' },
  { id: 'full', name: '完整配置' }
])

// 控制 API Key 显示状态
const showApiKeys = ref({})

// 控制提供商展开状态
const expandedProviders = ref({})

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

// 切换提供商展开状态
function toggleProvider(index) {
  expandedProviders.value[index] = !expandedProviders.value[index]
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

/* 标签页导航 */
.tabs {
  display: flex;
  gap: 5px;
  margin-bottom: 20px;
  border-bottom: 2px solid #e1e8ed;
}

.tab {
  padding: 10px 20px;
  background: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  cursor: pointer;
  font-size: 0.95em;
  font-weight: 500;
  color: #7f8c8d;
  transition: all 0.3s ease;
  position: relative;
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
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  padding: 20px;
  min-height: 400px;
}

.tab-pane {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 核心配置网格 */
.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 25px;
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

.config-table {
  padding: 0;
}

.config-row {
  display: flex;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
  font-size: 0.9em;
}

.config-row:last-child {
  border-bottom: none;
}

.config-row .label {
  font-weight: 600;
  color: #2c3e50;
  min-width: 100px;
  margin-right: 15px;
  display: flex;
  align-items: center;
}

.config-row .value {
  flex: 1;
  color: #34495e;
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
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8em;
  transition: all 0.2s ease;
}

.toggle-btn.small {
  padding: 3px 6px;
  font-size: 0.75em;
}

.toggle-btn:hover {
  background-color: #2980b9;
  transform: translateY(-1px);
}

.toggle-btn:active {
  transform: scale(0.95);
}

/* 提供商配置 */
.providers-section {
  padding: 10px 0;
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
  gap: 12px;
}

.provider-item {
  border: 1px solid #e1e8ed;
  border-radius: 8px;
  overflow: hidden;
  background: #fafcff;
  transition: all 0.3s ease;
}

.provider-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.provider-summary {
  display: flex;
  align-items: center;
  padding: 15px 20px;
  cursor: pointer;
  background: linear-gradient(90deg, #2c3e50 0%, #1a2530 100%);
  color: white;
}

.provider-name {
  flex: 1;
  font-weight: 600;
  font-size: 1em;
}

.provider-models-count {
  font-size: 0.9em;
  opacity: 0.9;
  margin-right: 15px;
}

.expand-icon {
  font-size: 0.9em;
  transition: transform 0.3s ease;
}

.provider-item.expanded .expand-icon {
  transform: rotate(180deg);
}

.provider-details {
  padding: 15px 20px;
  background: white;
  border-top: 1px solid #e1e8ed;
}

.detail-row {
  display: flex;
  padding: 8px 0;
  border-bottom: 1px solid #f5f5f5;
  font-size: 0.9em;
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-row .label {
  font-weight: 600;
  color: #2c3e50;
  min-width: 100px;
  margin-right: 15px;
}

.detail-row .value {
  flex: 1;
  color: #34495e;
  word-break: break-word;
}

.models-section {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.models-label {
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 10px;
  font-size: 0.9em;
}

.models-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.model-tag {
  background: linear-gradient(135deg, #e1f0fa 0%, #d1e7f5 100%);
  color: #2980b9;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.85em;
  font-weight: 500;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
}

.model-tag:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* 完整配置区域 */
.full-config-section h3 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 1.1em;
  font-weight: 600;
}

.full-config {
  width: 100%;
  height: 800px;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;
  font-family: 'Courier New', monospace;
  background: linear-gradient(180deg, #fdfdfd 0%, #f8f9fa 100%);
  resize: vertical;
  font-size: 0.9em;
  line-height: 1.5;
  box-shadow: inset 0 2px 5px rgba(0, 0, 0, 0.1);
  border: 1px solid #e1e8ed;
}

.loading, .error, .no-config {
  text-align: center;
  padding: 50px;
  font-size: 18px;
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  margin: 20px;
  animation: fadeIn 0.3s ease;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .system-config-container {
    padding: 10px;
  }
  
  .tabs {
    margin-bottom: 15px;
  }
  
  .tab {
    padding: 8px 15px;
    font-size: 0.9em;
  }
  
  .tab-content {
    padding: 15px;
  }
  
  .config-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .config-group {
    padding: 12px;
  }
  
  .config-group h3 {
    font-size: 1em;
    margin-bottom: 12px;
  }
  
  .config-row {
    flex-direction: column;
    gap: 4px;
    padding: 8px 0;
  }
  
  .config-row .label {
    min-width: auto;
    margin-right: 0;
    font-weight: 600;
  }
  
  .provider-summary {
    padding: 12px 15px;
  }
  
  .provider-name {
    font-size: 0.95em;
  }
  
  .provider-details {
    padding: 12px 15px;
  }
  
  .detail-row {
    flex-direction: column;
    gap: 4px;
    padding: 6px 0;
  }
  
  .detail-row .label {
    min-width: auto;
    margin-right: 0;
  }
  
  .models-grid {
    gap: 6px;
  }
  
  .model-tag {
    padding: 4px 10px;
    font-size: 0.8em;
  }
  
  .full-config {
    height: 300px;
  }
  
  .loading, .error, .no-config {
    padding: 30px;
    font-size: 16px;
    margin: 15px;
  }
}

@media (max-width: 480px) {
  .system-config-container {
    padding: 8px;
  }
  
  .tabs {
    flex-wrap: wrap;
  }
  
  .tab {
    padding: 6px 12px;
    font-size: 0.85em;
  }
  
  .tab-content {
    padding: 12px;
  }
  
  .config-group {
    padding: 10px;
  }
  
  .config-group h3 {
    font-size: 0.95em;
    margin-bottom: 10px;
  }
  
  .config-row,
  .detail-row {
    padding: 6px 0;
  }
  
  .config-row .label,
  .config-row .value,
  .detail-row .label,
  .detail-row .value {
    font-size: 0.85em;
  }
  
  .provider-summary {
    padding: 10px 12px;
  }
  
  .provider-name {
    font-size: 0.9em;
  }
  
  .provider-models-count {
    font-size: 0.8em;
    margin-right: 10px;
  }
  
  .provider-details {
    padding: 10px 12px;
  }
  
  .model-tag {
    padding: 3px 8px;
    font-size: 0.75em;
    margin: 1px;
  }
  
  .toggle-btn {
    padding: 2px 4px;
    font-size: 0.7em;
  }
  
  .full-config {
    height: 250px;
    padding: 10px;
    font-size: 0.8em;
  }
  
  .loading, .error, .no-config {
    padding: 20px;
    font-size: 14px;
    margin: 10px;
  }
}
</style>