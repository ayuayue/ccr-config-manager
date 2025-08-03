<template>
  <div class="claude-config-container">

    <el-row :gutter="10">
      <el-col :span="3">
        <el-menu :default-active="activeTab" class="el-menu-vertical-demo" @select="activeTab = $event">
          <el-menu-item index="service">服务管理</el-menu-item>
          <el-menu-item index="basic">基础配置</el-menu-item>
          <el-menu-item index="providers">提供商配置</el-menu-item>
          <el-menu-item index="full">完整配置</el-menu-item>
        </el-menu>
      </el-col>

      <el-col :span="20">
        <div v-if="activeTab === 'basic'">
          <el-row :gutter="10">
            <!-- 基础配置在左侧 -->
            <el-col :span="15">
              <el-card class="config-card">
                <div class="card-header">
                  <span>基础配置</span>
                </div>
                <el-form label-width="120px" label-position="left">
                  <el-form-item label="API Key">
                    <el-input type="password" v-model="config.APIKEY" placeholder="用于身份验证的密钥"></el-input>
                    <div class="help-text">设置后，客户端请求必须在 Authorization 请求头或 x-api-key 请求头中提供此密钥</div>
                  </el-form-item>

                  <el-form-item label="代理 URL">
                    <el-input v-model="config.PROXY_URL" placeholder="例如: http://127.0.0.1:7890"></el-input>
                    <div class="help-text">为 API 请求设置代理</div>
                  </el-form-item>

                  <el-form-item label="主机地址">
                    <el-input v-model="config.HOST" placeholder="例如: 0.0.0.0"></el-input>
                    <div class="help-text">设置服务的主机地址。如果未设置 APIKEY，出于安全考虑，主机地址将强制设置为 127.0.0.1</div>
                  </el-form-item>

                  <el-form-item label="端口">
                    <el-input-number v-model="config.PORT" :min="1" :max="65535"
                      placeholder="例如: 3456"></el-input-number>
                    <div class="help-text">设置服务的端口号，默认为 3456</div>
                  </el-form-item>

                  <el-form-item label="API 超时时间">
                    <el-input-number v-model="config.API_TIMEOUT_MS" :min="1"
                      placeholder="例如: 600000"></el-input-number>
                    <span class="unit">毫秒</span>
                  </el-form-item>

                  <el-form-item label="启用日志记录">
                    <el-switch v-model="config.LOG"></el-switch>
                    <div class="help-text">启用后，日志文件将位于 $HOME/.claude-code-router/claude-code-router.log</div>
                  </el-form-item>
                </el-form>
              </el-card>
            </el-col>

            <!-- 路由配置在右侧 -->
            <el-col :span="9">
              <el-card class="config-card">
                <div class="card-header">
                  <span>路由配置</span>
                </div>
                <el-form label-width="120px" label-position="left">
                  <el-form-item label="长上下文阈值">
                    <el-input-number v-model="config.Router.longContextThreshold" :min="1"
                      placeholder="例如: 60000"></el-input-number>
                  </el-form-item>

                  <el-form-item label="默认路由">
                    <el-select v-model="config.Router.default" placeholder="请选择" clearable>
                      <el-option value="" label="请选择"></el-option>
                      <el-option v-for="option in getProviderModelOptions()" :key="option.value" :label="option.label"
                        :value="option.value"></el-option>
                    </el-select>
                  </el-form-item>

                  <el-form-item label="长上下文路由">
                    <el-select v-model="config.Router.longContext" placeholder="请选择" clearable>
                      <el-option value="" label="请选择"></el-option>
                      <el-option v-for="option in getProviderModelOptions()" :key="option.value" :label="option.label"
                        :value="option.value"></el-option>
                    </el-select>
                  </el-form-item>

                  <el-form-item label="后台任务路由">
                    <el-select v-model="config.Router.background" placeholder="请选择" clearable>
                      <el-option value="" label="请选择"></el-option>
                      <el-option v-for="option in getProviderModelOptions()" :key="option.value" :label="option.label"
                        :value="option.value"></el-option>
                    </el-select>
                  </el-form-item>

                  <el-form-item label="推理任务路由">
                    <el-select v-model="config.Router.think" placeholder="请选择" clearable>
                      <el-option value="" label="请选择"></el-option>
                      <el-option v-for="option in getProviderModelOptions()" :key="option.value" :label="option.label"
                        :value="option.value"></el-option>
                    </el-select>
                  </el-form-item>

                  <el-form-item label="网络搜索路由">
                    <el-select v-model="config.Router.webSearch" placeholder="请选择" clearable>
                      <el-option value="" label="请选择"></el-option>
                      <el-option v-for="option in getProviderModelOptions()" :key="option.value" :label="option.label"
                        :value="option.value"></el-option>
                    </el-select>
                  </el-form-item>
                </el-form>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <div v-else-if="activeTab === 'service'">
          <el-row :gutter="10">
            <el-col :span="24">
              <el-card class="config-card">
                <div class="card-header">
                  <h3>服务管理</h3>
                </div>

                <!-- 服务状态部分 -->
                <el-row :gutter="10">
                  <el-col :span="24">
                    <div style="margin-bottom: 20px;">
                      <div class="card-header">
                        <span>服务状态</span>
                      </div>

                      <el-descriptions :column="1">
                        <el-descriptions-item label="运行状态">
                          <el-tag :type="serviceStatus.isRunning ? 'success' : 'danger'">
                            {{ serviceStatus.isRunning ? '运行中' : '已停止' }}
                          </el-tag>
                        </el-descriptions-item>
                        <el-descriptions-item label="进程ID">
                          <span v-if="serviceStatus.pid > 0">{{ serviceStatus.pid }}</span>
                          <span v-else>无</span>
                        </el-descriptions-item>
                        <el-descriptions-item label="版本号">
                          <span v-if="versionLoading">加载中...</span>
                          <span v-else>{{ ccrVersion || '未加载' }}</span>
                        </el-descriptions-item>
                      </el-descriptions>

                      <div style="margin-top: 20px; display: flex; justify-content: flex-end;">
                        <el-button type="primary" @click="refreshServiceStatus" style="margin-right: 10px;">
                          刷新状态
                        </el-button>
                        <el-button v-if="serviceStatus.isRunning" type="danger" @click="stopService"
                          :loading="serviceLoading.stop" style="margin-right: 10px;">
                          停止服务
                        </el-button>
                        <el-button v-if="!serviceStatus.isRunning" type="success" @click="startService"
                          :loading="serviceLoading.start" style="margin-right: 0;">
                          启动服务
                        </el-button>
                      </div>
                    </div>
                  </el-col>
                </el-row>

                <!-- 分割线 -->
                <el-divider></el-divider>

                <!-- 日志部分 -->
                <el-row :gutter="10">
                  <el-col :span="24">
                    <div>
                      <div class="card-header">
                        <span>服务日志</span>

                        <div style="display: flex; justify-content: flex-end;">
                          <el-button type="danger" @click="loadLogs" style="margin-right: 8px;">刷新日志</el-button>

                          <el-button type="danger" @click="clearLogs" style="margin-right: -10px;">清空日志</el-button>
                        </div>
                      </div>

                      <el-input type="textarea" v-model="logs" :rows="15" readonly
                        style="font-family: monospace; font-size: 12px;" placeholder="暂无日志内容"></el-input>
                    </div>
                  </el-col>
                </el-row>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <div v-else-if="activeTab === 'providers'">
          <el-row :gutter="10">
            <el-col :span="24">
              <el-card class="config-card">
                <div class="card-header">
                  <span>提供商配置 ({{ config.Providers?.length || 0 }} 个)</span>
                  <el-button type="primary" @click="addProvider">添加提供商</el-button>
                </div>
                <el-collapse v-model="expandedProviders" accordion>
                  <el-collapse-item v-for="(provider, index) in config.Providers" :key="index" :name="index">
                    <template #title>
                      <div class="provider-title">
                        <span>提供商 {{ index + 1 }}: {{ provider.name || '未命名' }}</span>
                      </div>
                    </template>
                    <el-form label-width="120px" label-position="left">
                      <el-form-item label="提供商名称">
                        <el-input v-model="provider.name" placeholder="例如: openrouter"></el-input>
                      </el-form-item>

                      <el-form-item label="API 基础 URL">
                        <el-input v-model="provider.api_base_url"
                          placeholder="例如: https://openrouter.ai/api/v1/chat/completions"></el-input>
                      </el-form-item>

                      <el-form-item label="API 密钥">
                        <el-input type="password" v-model="provider.api_key" placeholder="提供商的 API 密钥"></el-input>
                      </el-form-item>

                      <el-form-item label="模型列表">
                        <el-input type="textarea" v-model="provider.modelsText" placeholder="例如:
google/gemini-2.5-pro-preview
anthropic/claude-sonnet-4
anthropic/claude-3.5-sonnet" :rows="4"></el-input>
                      </el-form-item>

                      <el-form-item label="转换器 (JSON)">
                        <el-input type="textarea" v-model="provider.transformerText" placeholder='例如:
{
  "use": ["openrouter"]
}' :rows="4"></el-input>
                      </el-form-item>
                      <el-form-item>
                        <el-button type="danger" @click.stop="removeProvider(index)">删除提供商</el-button>
                      </el-form-item>
                    </el-form>
                  </el-collapse-item>
                </el-collapse>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <div v-else-if="activeTab === 'full'">
          <el-row :gutter="10">
            <el-col :span="24">
              <el-card class="config-card">
                <div class="card-header" style="display: flex; justify-content: space-between; align-items: center;">
                  <span>完整配置 (JSON)</span>
                  <el-button type="primary" @click="copyToClipboard(fullConfigJson)" size="small">复制到剪贴板</el-button>
                </div>
                <el-input type="textarea" v-model="fullConfigJson" :rows="35" class="full-config"></el-input>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { LoadConfig, SaveConfig, GetServiceStatus, StartService, StopService, RestartService, ReadLogs, ClearLogs, GetCCRVersion } from '../../wailsjs/go/main/App'
import { ClipboardSetText } from '../../wailsjs/runtime'
import {
  ElMenu, ElMenuItem, ElForm, ElFormItem, ElInput, ElSelect, ElOption,
  ElButton, ElCard, ElCollapse, ElCollapseItem, ElSwitch, ElMessage,
  ElRow, ElCol, ElTag, ElDivider, ElDescriptions, ElDescriptionsItem, ElNotification
} from 'element-plus'

// 标签页相关
const activeTab = ref('service')


// 监听标签页切换
watch(activeTab, (newTab) => {
  if (newTab === 'service') {
    // loadLogs()
    // 不再自动加载服务状态和版本号，用户可以手动点击刷新按钮
  }
})

// 控制提供商展开状态
const expandedProviders = ref([])


// 配置数据
const config = reactive({
  APIKEY: '',
  PROXY_URL: '',
  HOST: '',
  PORT: 3456,
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

// 服务管理数据
const serviceStatus = reactive({
  isRunning: false,
  pid: 0
})

// 服务控制按钮加载状态
const serviceLoading = reactive({
  start: false,
  stop: false,
  restart: false
})

// CCR版本号
const ccrVersion = ref('')

// 版本号加载状态
const versionLoading = ref(false)

// 版本号缓存
let versionCache = null
let cacheTimestamp = 0
const CACHE_DURATION = 5 * 60 * 1000 // 5分钟缓存

// 日志数据
const logs = ref('')

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
  ElNotification({
    title: type === 'success' ? '成功' : '错误',
    message: message,
    type: type,
    duration: 3000
  })
}

// 切换提供商展开状态
function toggleProvider(index) {
  const indexInArray = expandedProviders.value.indexOf(index)
  if (indexInArray > -1) {
    // 如果已经展开，则折叠
    expandedProviders.value.splice(indexInArray, 1)
  } else {
    // 如果未展开，则展开
    expandedProviders.value.push(index)
  }
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
      PORT: config.PORT || 3456,
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
    config.PORT = loadedConfig.PORT || 3456
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

// 加载服务状态
async function loadServiceStatus() {
  try {
    const status = await GetServiceStatus()
    serviceStatus.isRunning = status.isRunning
    serviceStatus.pid = status.pid
  } catch (error) {
    showStatus('加载服务状态时出错: ' + error.message, 'error')
  }
}

// 刷新服务状态（包括服务状态和版本号，不包括日志）
async function refreshServiceStatus() {
  try {
    // 显示加载状态
    versionLoading.value = true
    serviceLoading.start = true // 临时使用start的loading状态

    // 加载服务状态
    await loadServiceStatus()

    // 加载版本号
    await loadCCRVersion()

    showStatus('服务状态刷新成功', 'success')
  } catch (error) {
    showStatus('刷新服务状态时出错: ' + error.message, 'error')
  } finally {
    versionLoading.value = false
    serviceLoading.start = false
  }
}

// 启动服务
async function startService() {
  serviceLoading.start = true
  let timeoutId;
  try {
    // 先显示命令已发送提示
    showStatus('服务启动命令已发送', 'success')

    // 设置5秒超时
    const timeoutPromise = new Promise((_, reject) => {
      timeoutId = setTimeout(() => reject(new Error('操作超时，即将刷新页面')), 5000);
    });

    // 执行启动服务操作
    await Promise.race([StartService(), timeoutPromise]);
    clearTimeout(timeoutId);

    // 等待一段时间再检查状态，确保服务完全启动
    await new Promise(resolve => setTimeout(resolve, 1500))
    // 重试几次检查服务状态，确保服务完全启动
    let retries = 0
    const maxRetries = 5
    let success = false
    while (retries < maxRetries) {
      await loadServiceStatus()
      if (serviceStatus.isRunning) {
        showStatus('服务启动成功', 'success')
        success = true
        break
      }
      retries++
      if (retries < maxRetries) {
        await new Promise(resolve => setTimeout(resolve, 1000))
      }
    }
    if (!success && retries >= maxRetries) {
      showStatus('服务启动可能需要更多时间，请稍后手动刷新状态', 'warning')
    }
  } catch (error) {
    console.error('启动服务时出错:', error);
    if (error.message === '操作超时，即将刷新页面') {
      showStatus(error.message, 'warning');
      // 超时后刷新页面
      setTimeout(() => {
        window.location.reload();
      }, 2000);
    } else {
      // 显示更详细的错误信息
      const errorMessage = error.message || '未知错误';
      showStatus('启动服务时出错: ' + errorMessage, 'error');
      // 同时在控制台输出详细错误信息
      console.error('启动服务详细错误信息:', error);
    }
  } finally {
    clearTimeout(timeoutId);
    serviceLoading.start = false
  }
}

// 停止服务
async function stopService() {
  serviceLoading.stop = true
  let timeoutId;
  try {
    // 先显示命令已发送提示
    showStatus('服务停止命令已发送', 'success')

    // 设置5秒超时
    const timeoutPromise = new Promise((_, reject) => {
      timeoutId = setTimeout(() => reject(new Error('操作超时，即将刷新页面')), 5000);
    });

    // 执行停止服务操作
    await Promise.race([StopService(), timeoutPromise]);
    clearTimeout(timeoutId);

    // 等待一段时间再检查状态，确保服务完全停止
    await new Promise(resolve => setTimeout(resolve, 1500))
    // 重试几次检查服务状态，确保服务完全停止
    let retries = 0
    const maxRetries = 5
    let success = false
    while (retries < maxRetries) {
      await loadServiceStatus()
      if (!serviceStatus.isRunning) {
        showStatus('服务停止成功', 'success')
        success = true
        break
      }
      retries++
      if (retries < maxRetries) {
        await new Promise(resolve => setTimeout(resolve, 1000))
      }
    }
    if (!success && retries >= maxRetries) {
      showStatus('服务停止可能需要更多时间，请稍后手动刷新状态', 'warning')
    }
  } catch (error) {
    console.error('停止服务时出错:', error);
    if (error.message === '操作超时，即将刷新页面') {
      showStatus(error.message, 'warning');
      // 超时后刷新页面
      setTimeout(() => {
        window.location.reload();
      }, 2000);
    } else {
      // 显示更详细的错误信息
      const errorMessage = error.message || '未知错误';
      showStatus('停止服务时出错: ' + errorMessage, 'error');
      // 同时在控制台输出详细错误信息
      console.error('停止服务详细错误信息:', error);
    }
  } finally {
    clearTimeout(timeoutId);
    serviceLoading.stop = false
  }
}

// 重启服务
async function restartService() {
  serviceLoading.restart = true
  let timeoutId;
  try {
    // 先显示命令已发送提示
    showStatus('服务重启命令已发送', 'success')

    // 设置5秒超时
    const timeoutPromise = new Promise((_, reject) => {
      timeoutId = setTimeout(() => reject(new Error('操作超时，即将刷新页面')), 5000);
    });

    // 执行重启服务操作
    await Promise.race([RestartService(), timeoutPromise]);
    clearTimeout(timeoutId);

    // 等待一段时间再检查状态，确保服务完全重启
    await new Promise(resolve => setTimeout(resolve, 2000))
    // 重试几次检查服务状态，确保服务完全重启
    let retries = 0
    const maxRetries = 5
    let success = false
    while (retries < maxRetries) {
      await loadServiceStatus()
      if (serviceStatus.isRunning) {
        showStatus('服务重启成功', 'success')
        // 不再自动重新加载日志，用户可以手动点击刷新按钮
        success = true
        break
      }
      retries++
      if (retries < maxRetries) {
        await new Promise(resolve => setTimeout(resolve, 1000))
      }
    }
    if (!success && retries >= maxRetries) {
      showStatus('服务重启可能需要更多时间，请稍后手动刷新状态', 'warning')
    }
  } catch (error) {
    console.error('重启服务时出错:', error);
    if (error.message === '操作超时，即将刷新页面') {
      showStatus(error.message, 'warning');
      // 超时后刷新页面
      setTimeout(() => {
        window.location.reload();
      }, 2000);
    } else {
      // 显示更详细的错误信息
      const errorMessage = error.message || '未知错误';
      showStatus('重启服务时出错: ' + errorMessage, 'error');
      // 同时在控制台输出详细错误信息
      console.error('重启服务详细错误信息:', error);
    }
  } finally {
    clearTimeout(timeoutId);
    serviceLoading.restart = false
  }
}

// 加载日志
async function loadLogs() {
  try {
    const logContent = await ReadLogs()
    logs.value = logContent || ''
  } catch (error) {
    console.error('加载日志时出错:', error)
    logs.value = '读取日志时出错: ' + (error.message || '未知错误')
  }
}

// 加载CCR版本号
async function loadCCRVersion() {
  // 检查缓存
  const now = Date.now()
  if (versionCache && (now - cacheTimestamp) < CACHE_DURATION) {
    ccrVersion.value = versionCache
    return
  }

  versionLoading.value = true
  try {
    const version = await GetCCRVersion()
    ccrVersion.value = version
    // 更新缓存
    versionCache = version
    cacheTimestamp = now
  } catch (error) {
    ccrVersion.value = '无法获取版本号: ' + error.message
    // 清除无效缓存
    versionCache = null
    cacheTimestamp = 0
  } finally {
    versionLoading.value = false
  }
}

// 清空日志
async function clearLogs() {
  try {
    await ClearLogs()
    logs.value = ''
    showStatus('日志已清空', 'success')
    // 不再自动重新加载日志，用户可以手动点击刷新按钮
  } catch (error) {
    console.error('清空日志时出错:', error)
    showStatus('清空日志时出错: ' + (error.message || '未知错误'), 'error')
  }
}

// 复制到剪贴板
async function copyToClipboard(text) {
  try {
    await ClipboardSetText(text)
    showStatus('已复制到剪贴板', 'success')
  } catch (error) {
    showStatus('复制失败: ' + error.message, 'error')
  }
}


// 页面加载完成后初始化
onMounted(() => {
  // 页面加载时自动加载配置
  loadConfig()

  // 如果当前是服务管理页面，1秒后自动刷新服务状态
  if (activeTab.value === 'service') {
    // 1秒后自动刷新服务状态（不包括日志）
    setTimeout(() => {
      refreshServiceStatus()
    }, 1000)
  }

  // 添加事件监听器
  window.addEventListener('reload-config', loadConfig)
  window.addEventListener('save-config', saveConfig)

  // 在组件卸载时移除事件监听器
  onUnmounted(() => {
    window.removeEventListener('reload-config', loadConfig)
    window.removeEventListener('save-config', saveConfig)
    // 清除定时器
    if (versionLoadTimeout) {
      clearTimeout(versionLoadTimeout)
    }
  })
})
</script>

<style scoped>
.claude-config-container {
  width: 100%;
  height: 100%;
  background-color: #ffffff;
  color: #000000;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  overflow: auto;
  padding: 5px;
  box-sizing: border-box;
  position: relative;
}

/* 侧边菜单 */
.el-menu-vertical-demo {
  min-height: 400px;
  border-right: 1px solid #e6e6e6;
  margin-right: 0;
}

.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 150px;
}

/* 标签页内容 */
.tab-content {
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.08);
  padding: 10px;
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
  gap: 15px;
  margin-left: 0;
}

.config-group {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 5px;
  border: 1px solid #e1e8ed;
  margin-left: 0;
}

.config-group h3 {
  margin: 0 0 10px 0;
  color: #000000;
  font-size: 1em;
  font-weight: 600;
  padding-bottom: 5px;
  border-bottom: 1px solid #eee;
}

.card-header {
  /* background-color: #f5f7fa; */
  padding: 15px 20px;
  /* border-bottom: 1px solid #ebeef5; */
  border-radius: 8px 8px 0 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-weight: 600;
  color: #303133;
  font-size: 16px;
}

/* 统一卡片样式 */
.config-card {
  margin-bottom: 20px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  background-color: #fff;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 8px;
}

.form-item {
  margin-bottom: 0;
  padding: 5px;
  border-radius: 6px;
  background-color: #f5f5f5;
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
  margin-bottom: 4px;
  font-weight: 600;
  color: #000000;
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
  padding: 5px 8px;
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
  padding: 5px 8px;
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
  padding: 5px 8px;
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
  padding: 5px 0;
  max-width: 1200px;
}

.providers-header h3 {
  margin: 0 0 10px 0;
  color: #2c3e50;
  font-size: 1em;
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
  background-color: #ffffff;
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
  background-color: #f5f5f5;
  color: #000000;
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
  background-color: #ffffff;
  border-top: 1px solid #e1e8ed;
}

/* 完整配置区域 */
.full-config-section {
  max-width: 1200px;
  padding: 5px;
}

.full-config-section h3 {
  margin: 0 0 10px 0;
  color: #000000;
  font-size: 1em;
  font-weight: 600;
}

.full-config {
  width: 1000;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;
  font-family: 'Courier New', monospace;
  background-color: #ffffff;
  resize: vertical;
  font-size: 0.85em;
  line-height: 1.4;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e1e8ed;
  margin-bottom: 15px;
}

.actions {
  text-align: center;
  padding: 10px 0;
}

.tab-actions {
  text-align: center;
  padding: 10px 0;
}

button {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9em;
  font-weight: 500;
  margin-right: 6px;
  margin-bottom: 6px;
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


.unit {
  margin-left: 10px;
  color: #666;
}

.top-actions {
  position: absolute;
  top: 10px;
  right: 20px;
  z-index: 100;
}

.provider-title {
  font-weight: bold;
  font-size: 0.9em;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .claude-config-container {
    padding: 5px;
  }

  .tabs {
    margin-bottom: 8px;
    max-width: 100%;
  }

  .tab {
    padding: 4px 8px;
    font-size: 0.8em;
  }

  .tab-content {
    padding: 8px;
    min-height: 250px;
  }

  .config-grid {
    grid-template-columns: 1fr;
    gap: 10px;
  }

  .config-group {
    padding: 6px;
  }

  .config-group h3 {
    font-size: 1em;
    margin-bottom: 6px;
  }

  .form-grid {
    gap: 6px;
  }

  .form-item {
    padding: 5px;
  }

  input,
  textarea,
  select,
  .route-select,
  .route-input {
    max-width: 100%;
  }

  .provider-summary {
    padding: 8px 10px;
  }

  .provider-name {
    font-size: 0.85em;
  }

  .provider-details {
    padding: 8px 10px;
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
    padding: 8px;
  }

  .status {
    padding: 8px 10px;
    font-size: 0.8em;
  }
}

@media (max-width: 480px) {
  .claude-config-container {
    padding: 3px;
  }

  .tabs {
    flex-wrap: wrap;
    margin-bottom: 6px;
  }

  .tab {
    padding: 3px 6px;
    font-size: 0.75em;
  }

  .tab-content {
    padding: 6px;
    min-height: 200px;
  }

  .config-group {
    padding: 5px;
  }

  .config-group h3 {
    font-size: 0.9em;
    margin-bottom: 5px;
  }

  .form-grid {
    gap: 5px;
  }

  .form-item {
    padding: 4px;
  }

  label {
    font-size: 0.75em;
  }

  input,
  textarea,
  select,
  .route-select,
  .route-input {
    padding: 4px 6px;
    font-size: 0.75em;
  }

  .provider-summary {
    padding: 6px 8px;
  }

  .provider-name {
    font-size: 0.8em;
  }

  .provider-details {
    padding: 6px 8px;
  }

  .full-config {
    height: 200px;
    padding: 6px;
    font-size: 0.7em;
  }

  button {
    padding: 4px 6px;
    font-size: 0.7em;
    margin-right: 3px;
    margin-bottom: 3px;
  }

  .status {
    padding: 6px 8px;
    font-size: 0.75em;
  }
}
</style>