<template>
  <div class="system-config-container">
    <el-card v-if="loading" shadow="never">
      <div class="loading">
        正在加载配置...
      </div>
    </el-card>
    
    <el-card v-else-if="error" shadow="never">
      <div class="error">
        {{ error }}
      </div>
    </el-card>
    
    <div v-else-if="config">
      <el-row :gutter="10">
        <el-col :span="3">
          <el-menu
            :default-active="activeTab"
            class="el-menu-vertical-demo"
            @select="activeTab = $event"
          >
            <el-menu-item index="core">核心配置</el-menu-item>
            <el-menu-item index="providers">提供商配置</el-menu-item>
            <el-menu-item index="full">完整配置</el-menu-item>
          </el-menu>
        </el-col>
        
        <el-col :span="20">
          <div v-if="activeTab === 'core'">
            <div class="config-grid">
              <el-card class="config-card">
                <template #header>
                  <div class="card-header">
                    <span>核心配置</span>
                  </div>
                </template>
                <el-descriptions :column="1" >
                  <el-descriptions-item label="API Key">
                    <span v-if="config.APIKEY">
                      <span :class="{ masked: !showApiKeys['main'] }">{{ getDisplayApiKey(config.APIKEY, !showApiKeys['main']) }}</span>
                      <el-button size="small" @click="toggleApiKeyVisibility('main')">
                        {{ showApiKeys['main'] ? '隐藏' : '显示' }}
                      </el-button>
                    </span>
                    <span v-else>未设置</span>
                  </el-descriptions-item>
                  
                  <el-descriptions-item label="代理设置">
                    {{ config.PROXY_URL || '未设置' }}
                  </el-descriptions-item>
                  
                  <el-descriptions-item label="主机地址">
                    {{ config.HOST || '默认 (127.0.0.1)' }}
                  </el-descriptions-item>
                  
                  <el-descriptions-item label="超时时间">
                    {{ config.API_TIMEOUT_MS || 600000 }} 毫秒
                  </el-descriptions-item>
                  
                  <el-descriptions-item label="日志记录">
                    <el-tag :type="config.LOG ? 'success' : 'info'">{{ config.LOG ? '已启用' : '已禁用' }}</el-tag>
                  </el-descriptions-item>
                </el-descriptions>
              </el-card>
              
              <el-card class="config-card">
                <template #header>
                  <div class="card-header">
                    <span>路由配置</span>
                  </div>
                </template>
                <el-descriptions :column="1">
                  <el-descriptions-item label="默认路由">
                    {{ config.Router?.default || '未设置' }}
                  </el-descriptions-item>
                  
                  <el-descriptions-item label="后台任务">
                    {{ config.Router?.background || '未设置' }}
                  </el-descriptions-item>
                  
                  <el-descriptions-item label="推理任务">
                    {{ config.Router?.think || '未设置' }}
                  </el-descriptions-item>
                  
                  <el-descriptions-item label="长上下文">
                    {{ config.Router?.longContext || '未设置' }}
                  </el-descriptions-item>
                  
                  <el-descriptions-item label="网络搜索">
                    {{ config.Router?.webSearch || '未设置' }}
                  </el-descriptions-item>
                </el-descriptions>
              </el-card>
            </div>
          </div>
          
          <div v-else-if="activeTab === 'providers'">
            <el-card class="providers-section">
              <template #header>
                <div class="card-header">
                  <span>提供商配置 ({{ config.Providers?.length || 0 }} 个)</span>
                </div>
              </template>
              <el-collapse v-model="expandedProviders" accordion>
                <el-collapse-item 
                  v-for="(provider, index) in config.Providers" 
                  :key="index" 
                  :name="index"
                >
                  <template #title>
                    <div class="provider-title">
                      <span>{{ provider.name }}</span>
                      <el-tag size="small">{{ provider.models?.length || 0 }} 个模型</el-tag>
                    </div>
                  </template>
                  <el-descriptions :column="1" >
                    <el-descriptions-item label="API URL">
                      {{ provider.api_base_url }}
                    </el-descriptions-item>
                    
                    <el-descriptions-item label="API Key">
                      <span>
                        <span :class="{ masked: !showApiKeys['provider-' + index] }">{{ getDisplayApiKey(provider.api_key, !showApiKeys['provider-' + index]) }}</span>
                        <el-button size="small" @click.stop="toggleApiKeyVisibility('provider-' + index)">
                          {{ showApiKeys['provider-' + index] ? '隐藏' : '显示' }}
                        </el-button>
                      </span>
                    </el-descriptions-item>
                  </el-descriptions>
                  
                  <div class="models-section" v-if="provider.models && provider.models.length > 0">
                    <div class="models-label">支持的模型:</div>
                    <div class="models-grid">
                      <el-tag v-for="(model, modelIndex) in provider.models" :key="modelIndex" class="model-tag">
                        {{ model }}
                      </el-tag>
                    </div>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </el-card>
          </div>
          
          <div v-else-if="activeTab === 'full'">
            <el-card class="full-config-section">
              <template #header>
                <div class="card-header">
                  <span>完整配置 (JSON)</span>
                </div>
              </template>
              <el-input
                type="textarea"
                readonly
                :rows="30"
                :value="fullConfigJson"
                class="full-config"
              />
            </el-card>
          </div>
        </el-col>
      </el-row>
    </div>
    
    <el-card v-else shadow="never">
      <div class="no-config">
        <p>未找到配置文件</p>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { LoadConfig } from '../../wailsjs/go/main/App'
import { ElMenu, ElMenuItem, ElRow, ElCol, ElTable, ElTableColumn, ElCard, ElTag, ElCollapse, ElCollapseItem, ElButton } from 'element-plus'

// 响应式数据
const config = ref(null)
const loading = ref(true)
const error = ref('')

// 标签页相关
const activeTab = ref('core')

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
  background-color: #ffffff;
  color: #000000;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  overflow: auto;
  padding: 5px;
  box-sizing: border-box;
}

/* 侧边菜单 */
.el-menu-vertical-demo {
  min-height: 800px;
  border-right: 1px solid #e6e6e6;
  margin-right: 0;
}

.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 150px;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.provider-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  color: #000000;
  font-size: 0.9em;
}

/* 核心配置网格 */
.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 15px;
  margin-left: 0;
}

.config-card {
  margin-bottom: 10px;
  padding: 5px;
  margin-left: 0;
}

.card-header {
  background-color: #f5f7fa;
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
  border-radius: 4px 4px 0 0;
  font-weight: 600;
  color: #303133;
  font-size: 16px;
}

.masked {
  color: #95a5a6;
  font-family: monospace;
}

/* 提供商配置 */
.providers-section {
  padding: 5px 0;
}

.models-section {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #f0f0f0;
}

.models-label {
  font-weight: 600;
  color: #000000;
  margin-bottom: 8px;
  font-size: 0.85em;
  text-align: left;
}

.models-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.model-tag {
  margin: 2px;
  padding: 2px 6px;
  font-size: 0.8em;
}

/* 完整配置区域 */
.full-config {
  font-family: 'Courier New', monospace;
  font-size: 0.9em;
  line-height: 1.5;
}

.loading, .error, .no-config {
  text-align: center;
  padding: 20px;
  font-size: 16px;
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
}

.error::before {
  content: "⚠️";
  margin-right: 10px;
}

.no-config {
  color: #7f8c8d;
}

.no-config::before {
  content: "📋";
  margin-right: 10px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .system-config-container {
    padding: 5px;
  }
  
  .tabs {
    margin-bottom: 10px;
  }
  
  .tab {
    padding: 6px 12px;
    font-size: 0.9em;
  }
  
  .tab-content {
    padding: 10px;
  }
  
  .config-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .config-group {
    padding: 8px;
  }
  
  .config-group h3 {
    font-size: 1em;
    margin-bottom: 8px;
  }
  
  .config-row {
    flex-direction: column;
    gap: 4px;
    padding: 6px 0;
  }
  
  .config-row .label {
    min-width: auto;
    margin-right: 0;
    font-weight: 600;
  }
  
  .provider-summary {
    padding: 10px 12px;
  }
  
  .provider-name {
    font-size: 0.95em;
  }
  
  .provider-details {
    padding: 10px 12px;
  }
  
  .detail-row {
    flex-direction: column;
    gap: 4px;
    padding: 4px 0;
  }
  
  .detail-row .label {
    min-width: auto;
    margin-right: 0;
  }
  
  .models-grid {
    gap: 4px;
  }
  
  .model-tag {
    padding: 2px 6px;
    font-size: 0.8em;
  }
  
  .full-config {
    height: 300px;
  }
  
  .loading, .error, .no-config {
    padding: 20px;
    font-size: 14px;
    margin: 10px;
  }
}

@media (max-width: 480px) {
  .system-config-container {
    padding: 5px;
  }
  
  .tabs {
    flex-wrap: wrap;
  }
  
  .tab {
    padding: 4px 8px;
    font-size: 0.8em;
  }
  
  .tab-content {
    padding: 8px;
  }
  
  .config-group {
    padding: 6px;
  }
  
  .config-group h3 {
    font-size: 0.9em;
    margin-bottom: 6px;
  }
  
  .config-row,
  .detail-row {
    padding: 4px 0;
  }
  
  .config-row .label,
  .config-row .value,
  .detail-row .label,
  .detail-row .value {
    font-size: 0.8em;
  }
  
  .provider-summary {
    padding: 8px 10px;
  }
  
  .provider-name {
    font-size: 0.9em;
  }
  
  .provider-models-count {
    font-size: 0.75em;
    margin-right: 8px;
  }
  
  .provider-details {
    padding: 8px 10px;
  }
  
  .model-tag {
    padding: 2px 6px;
    font-size: 0.7em;
    margin: 1px;
  }
  
  .toggle-btn {
    padding: 1px 3px;
    font-size: 0.65em;
  }
  
  .full-config {
    height: 250px;
    padding: 8px;
    font-size: 0.75em;
  }
  
  .loading, .error, .no-config {
    padding: 15px;
    font-size: 12px;
    margin: 8px;
  }
}
</style>