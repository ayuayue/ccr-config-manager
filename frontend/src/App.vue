<script setup>
import { ref } from 'vue'
import SystemConfig from './components/SystemConfig.vue'
import ClaudeConfig from './components/ClaudeConfig.vue'
import { BrowserOpenURL, ClipboardSetText } from '../wailsjs/runtime'

const currentView = ref('home')

const openGitHub = () => {
  BrowserOpenURL('https://github.com/ayuayue/ccr-config-manager')
}

const copyGitHubUrl = () => {
  ClipboardSetText('https://github.com/ayuayue/ccr-config-manager')
  alert('GitHub 地址已复制到剪贴板')
}
</script>

<template>
  <nav class="navigation">
    <button 
      :class="{ active: currentView === 'home' }" 
      @click="currentView = 'home'"
    >
      系统配置
    </button>
    <button 
      :class="{ active: currentView === 'config' }" 
      @click="currentView = 'config'"
    >
      配置管理
    </button>
    <div class="nav-spacer"></div>
    <div class="nav-info">
      <span class="author">作者: caoayu</span>
      <button class="github-link" @click="openGitHub" @contextmenu.prevent="copyGitHubUrl">
        GitHub 仓库
      </button>
    </div>
  </nav>

  <div class="view-container">
    <SystemConfig v-if="currentView === 'home'" />
    <ClaudeConfig v-if="currentView === 'config'" />
  </div>
</template>

<style>
.navigation {
  display: flex;
  justify-content: center;
  padding: 15px 20px;
  background-color: #2c3e50;
  box-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
}

.navigation button {
  background-color: transparent;
  color: #ecf0f1;
  border: none;
  padding: 8px 16px;
  margin: 0 5px;
  cursor: pointer;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.navigation button:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: white;
}

.navigation button.active {
  background-color: #3498db;
  color: white;
}

.nav-spacer {
  flex-grow: 1;
}

.nav-info {
  display: flex;
  align-items: center;
  gap: 15px;
  color: #ecf0f1;
  font-size: 14px;
}

.author {
  white-space: nowrap;
}

.github-link {
  background-color: #3498db;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s ease;
  white-space: nowrap;
}

.github-link:hover {
  background-color: #2980b9;
}

.view-container {
  height: calc(100vh - 50px);
}
</style>
