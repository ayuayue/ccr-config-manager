<script setup>
import { ref } from 'vue'
import ClaudeConfig from './components/ClaudeConfig.vue'
import { BrowserOpenURL, ClipboardSetText } from '../wailsjs/runtime'
import { ElButton, ElContainer, ElHeader, ElMain, ElMenu, ElMenuItem, ElMessageBox, ElIcon, ElDropdown, ElDropdownMenu, ElDropdownItem } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { ReadREADME, LoadConfig, SaveConfig } from '../wailsjs/go/main/App'

const openGitHub = () => {
  BrowserOpenURL('https://github.com/ayuayue/ccr-config-manager')
}

// 加载配置
const loadConfig = async () => {
  try {
    // 调用ClaudeConfig组件的加载配置方法
    // 这里我们需要通过事件来通知ClaudeConfig组件重新加载配置
    const event = new CustomEvent('reload-config')
    window.dispatchEvent(event)
  } catch (error) {
    ElMessageBox.alert('加载配置时出错: ' + error.message, '错误', {
      confirmButtonText: '确定',
      type: 'error'
    })
  }
}

// 保存配置
const saveConfig = async () => {
  try {
    // 调用ClaudeConfig组件的保存配置方法
    // 这里我们需要通过事件来通知ClaudeConfig组件保存配置
    const event = new CustomEvent('save-config')
    window.dispatchEvent(event)
  } catch (error) {
    ElMessageBox.alert('保存配置时出错: ' + error.message, '错误', {
      confirmButtonText: '确定',
      type: 'error'
    })
  }
}

</script>

<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <div class="header-content">
        <div class="header-left">
          <el-button class="about-button" @click="openGitHub">
            关于
          </el-button>
        </div>
        <h1>CCR 配置管理器</h1>
        <div class="nav-info">
          <el-button @click="loadConfig">刷新配置</el-button>
          <el-button type="primary" @click="saveConfig">保存配置</el-button>
        </div>
      </div>
    </el-header>

    <el-main class="app-main">
      <div class="view-container">
        <ClaudeConfig />
      </div>
    </el-main>
  </el-container>
</template>

<style>
.app-container {
  background-color: #ffffff;
  height: 100vh;
}

.app-header {
  padding: 0 20px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-content h1 {
  margin: 0;
  color: #333333;
  font-size: 24px;
  font-weight: 600;
}

.nav-info {
  display: flex;
  align-items: center;
  gap: 15px;
  color: #333333;
  font-size: 14px;
}

.about-button {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 8px 12px;
}

.about-button:hover {
  background-color: #f5f7fa;
  border-color: #c0c4cc;
}

/* 关于对话框样式 */
.about-dialog {
  width: 800px !important;
  max-width: 90vw !important;
}

.about-dialog .el-message-box__content {
  max-height: 60vh;
  overflow-y: auto;
}

.about-dialog .el-message-box__message {
  padding: 20px;
}

.app-main {
  background-color: #ffffff;
  padding: 0;
}

.view-container {
  height: calc(100vh - 61px);
  background-color: #ffffff;
  padding: 10px;
}
</style>