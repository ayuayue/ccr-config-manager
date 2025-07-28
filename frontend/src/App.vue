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
</style>
