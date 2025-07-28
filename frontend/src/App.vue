<script setup>
import { ref } from 'vue'
import SystemConfig from './components/SystemConfig.vue'
import ClaudeConfig from './components/ClaudeConfig.vue'

const currentView = ref('home')
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
  </nav>

  <div class="view-container">
    <SystemConfig v-if="currentView === 'home'" />
    <ClaudeConfig v-if="currentView === 'config'" />
  </div>

  <footer class="app-footer">
    <div class="footer-content">
      <div class="footer-section">
        <h4>Claude Code Router 配置管理器</h4>
        <p>一个用于管理和配置 Claude Code Router 的图形化工具</p>
      </div>
      <div class="footer-section">
        <h4>作者</h4>
        <p>caoayu</p>
      </div>
      <div class="footer-section">
        <h4>项目链接</h4>
        <p><a href="https://github.com/ayuayue/ccr-config-manager" target="_blank" rel="noopener noreferrer">GitHub 仓库</a></p>
        <p><a href="https://github.com/alexrider203/claud-code-router" target="_blank" rel="noopener noreferrer">Claude Code Router</a></p>
      </div>
    </div>
  </footer>
</template>

<style>
.navigation {
  display: flex;
  justify-content: center;
  padding: 15px 0;
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

.view-container {
  height: calc(100vh - 50px);
}

.app-footer {
  background-color: #f8f9fa;
  border-top: 1px solid #e9ecef;
  padding: 20px 0;
  margin-top: auto;
}

.footer-content {
  display: flex;
  justify-content: space-around;
  max-width: 1200px;
  margin: 0 auto;
  flex-wrap: wrap;
}

.footer-section {
  flex: 1;
  min-width: 200px;
  padding: 0 15px;
  margin-bottom: 15px;
}

.footer-section h4 {
  color: #2c3e50;
  margin-bottom: 10px;
  font-size: 16px;
}

.footer-section p {
  color: #6c757d;
  font-size: 14px;
  margin: 5px 0;
}

.footer-section a {
  color: #3498db;
  text-decoration: none;
}

.footer-section a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .footer-content {
    flex-direction: column;
    align-items: center;
  }
  
  .footer-section {
    text-align: center;
    min-width: auto;
  }
}
</style>
