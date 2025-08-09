<script setup>
import ClaudeConfig from './components/ClaudeConfig.vue'
import { BrowserOpenURL } from '../wailsjs/runtime'
import { ElButton, ElContainer, ElHeader, ElMain, ElMenu, ElMenuItem, ElMessageBox, ElIcon, ElDropdown, ElDropdownMenu, ElDropdownItem } from 'element-plus'

const openGitHub = () => {
  BrowserOpenURL('https://github.com/ayuayue/ccr-config-manager/blob/master/README.md')
}

import { GetAppVersion, GetLatestVersionFromGitHub, CompareVersions, DownloadUpdate } from '../wailsjs/go/main/App'
import { ElCard, ElMessage, ElMessageBox as MessageBox } from 'element-plus'

const showAbout = async () => {
  try {
    const version = await GetAppVersion()
    
    // 创建弹窗内容
    const content = `
      <div style="padding: 120px; text-align: center; border-top: 10px">
        <el-card style="border: none; box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);">
          <div style="padding: 20px;">
            <h2>CCR 配置管理器</h2>
            <p style="font-size: 16px; margin: 20px 0;"><strong>版本:</strong> ${version}</p>
            <p style="font-size: 14px; margin: 20px 0; line-height: 1.6;">
              这是一个用于管理和配置 Claude Code Router (CCR) 的图形化工具。
            </p>
            <div style="margin-top: 10px; display: flex; flex-direction: column; gap: 10px;">
              <button id="github-btn" style="
                background-color: #409eff; 
                border: none; 
                color: white; 
                padding: 10px 20px; 
                text-align: center; 
                text-decoration: none; 
                display: inline-block; 
                font-size: 14px; 
                border-radius: 4px; 
                cursor: pointer;
                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
              ">
                <i class="el-icon-link" style="margin-right: 5px;"></i>
                访问 GitHub 项目
              </button>
              <button id="check-update-btn" style="
                background-color: #67c23a; 
                border: none; 
                color: white; 
                padding: 10px 20px; 
                text-align: center; 
                text-decoration: none; 
                display: inline-block; 
                font-size: 14px; 
                border-radius: 4px; 
                cursor: pointer;
                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
              ">
                <i class="el-icon-refresh" style="margin-right: 5px;"></i>
                检查更新
              </button>
            </div>
          </div>
        </el-card>
      </div>
    `
    
    // 创建弹窗
    MessageBox.alert('', '关于', {
      message: content,
      dangerouslyUseHTMLString: true,
      customClass: 'about-dialog',
      callback: () => {
        // 清理事件监听器
        const githubBtn = document.getElementById('github-btn')
        const checkUpdateBtn = document.getElementById('check-update-btn')
        if (githubBtn) {
          githubBtn.removeEventListener('click', handleGithubClick)
        }
        if (checkUpdateBtn) {
          checkUpdateBtn.removeEventListener('click', handleCheckUpdate)
        }
      }
    })
    
    // 添加GitHub按钮点击事件
    const handleGithubClick = () => {
      BrowserOpenURL('https://github.com/ayuayue/ccr-config-manager')
    }
    
    // 添加检查更新按钮点击事件
    const handleCheckUpdate = async () => {
      try {
        // 显示检查更新提示
        ElMessage.info('正在检查更新...')
        
        // 获取最新版本
        const latestVersion = await GetLatestVersionFromGitHub()
        
        // 比较版本
        const hasUpdate = await CompareVersions(version, latestVersion)
        
        if (hasUpdate) {
          // 有新版本，提示用户
          MessageBox.confirm(
            `发现新版本 ${latestVersion}，当前版本 ${version}。是否要下载并安装更新？`,
            '发现新版本',
            {
              confirmButtonText: '更新',
              cancelButtonText: '取消',
              type: 'warning'
            }
          ).then(async () => {
            try {
              // 下载更新
              ElMessage.info('正在下载更新...')
              const filePath = await DownloadUpdate(latestVersion)
              ElMessage.success(`更新已下载到: ${filePath}`)
              
              // 提示用户手动安装
              MessageBox.alert(
                `更新已下载到: ${filePath}。请先卸载旧版本，再手动安装应用程序。`,
                '下载完成',
                {
                  confirmButtonText: '确定',
                  type: 'success'
                }
              )
            } catch (error) {
              ElMessage.error(`下载更新失败: ${error.message || error}`)
            }
          }).catch(() => {
            // 用户取消更新
            ElMessage.info('已取消更新')
          })
        } else {
          // 没有新版本
          ElMessage.success('当前已是最新版本')
        }
      } catch (error) {
        console.error('检查更新失败:', error)
        ElMessage.error(`检查更新失败: ${error.message || error || '未知错误'}`)
      }
    }
    
    // 等待DOM更新后添加事件监听器
    setTimeout(() => {
      const githubBtn = document.getElementById('github-btn')
      const checkUpdateBtn = document.getElementById('check-update-btn')
      if (githubBtn) {
        githubBtn.addEventListener('click', handleGithubClick)
      }
      if (checkUpdateBtn) {
        checkUpdateBtn.addEventListener('click', handleCheckUpdate)
      }
    }, 0)
  } catch (error) {
    // 创建弹窗内容（错误情况）
    const content = `
      <div style="padding: 20px; text-align: center;">
        <el-card style="border: none; box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);">
          <div style="padding: 20px;">
            <h2>CCR 配置管理器</h2>
            <p style="font-size: 16px; margin: 20px 0;"><strong>版本:</strong> 未知</p>
            <p style="font-size: 14px; margin: 20px 0; line-height: 1.6;">
              这是一个用于管理和配置 Claude Code Router (CCR) 的图形化工具。
            </p>
            <div style="margin-top: 30px; display: flex; flex-direction: column; gap: 10px;">
              <button id="github-btn-error" style="
                background-color: #409eff; 
                border: none; 
                color: white; 
                padding: 10px 20px; 
                text-align: center; 
                text-decoration: none; 
                display: inline-block; 
                font-size: 14px; 
                border-radius: 4px; 
                cursor: pointer;
                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
              ">
                <i class="el-icon-link" style="margin-right: 5px;"></i>
                访问 GitHub 项目
              </button>
              <button id="check-update-btn-error" style="
                background-color: #67c23a; 
                border: none; 
                color: white; 
                padding: 10px 20px; 
                text-align: center; 
                text-decoration: none; 
                display: inline-block; 
                font-size: 14px; 
                border-radius: 4px; 
                cursor: pointer;
                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
              ">
                <i class="el-icon-refresh" style="margin-right: 5px;"></i>
                检查更新
              </button>
            </div>
          </div>
        </el-card>
      </div>
    `
    
    // 创建弹窗
    MessageBox.alert('', '关于', {
      message: content,
      dangerouslyUseHTMLString: true,
      customClass: 'about-dialog',
      callback: () => {
        // 清理事件监听器
        const githubBtn = document.getElementById('github-btn-error')
        const checkUpdateBtn = document.getElementById('check-update-btn-error')
        if (githubBtn) {
          githubBtn.removeEventListener('click', handleGithubClick)
        }
        if (checkUpdateBtn) {
          checkUpdateBtn.removeEventListener('click', handleCheckUpdate)
        }
      }
    })
    
    // 添加GitHub按钮点击事件
    const handleGithubClick = () => {
      BrowserOpenURL('https://github.com/ayuayue/ccr-config-manager')
    }
    
    // 添加检查更新按钮点击事件
    const handleCheckUpdate = async () => {
      try {
        // 显示检查更新提示
        ElMessage.info('正在检查更新...')
        
        // 获取最新版本
        const latestVersion = await GetLatestVersionFromGitHub()
        
        // 比较版本
        const hasUpdate = await CompareVersions(version, latestVersion)
        
        if (hasUpdate) {
          // 有新版本，提示用户
          MessageBox.confirm(
            `发现新版本 ${latestVersion}，当前版本 ${version}。是否要下载并安装更新？`,
            '发现新版本',
            {
              confirmButtonText: '更新',
              cancelButtonText: '取消',
              type: 'warning'
            }
          ).then(async () => {
            try {
              // 下载更新
              ElMessage.info('正在下载更新...')
              const filePath = await DownloadUpdate(latestVersion)
              ElMessage.success(`更新已下载到: ${filePath}`)
              
              // 提示用户手动安装
              MessageBox.alert(
                `更新已下载到: ${filePath}。请手动解压并替换当前应用程序。`,
                '下载完成',
                {
                  confirmButtonText: '确定',
                  type: 'success'
                }
              )
            } catch (error) {
              ElMessage.error(`下载更新失败: ${error.message}`)
            }
          }).catch(() => {
            // 用户取消更新
            ElMessage.info('已取消更新')
          })
        } else {
          // 没有新版本
          ElMessage.success('当前已是最新版本')
        }
      } catch (error) {
        console.error('检查更新失败:', error)
        ElMessage.error(`检查更新失败: ${error.message || error || '未知错误'}`)
      }
    }
    
    // 等待DOM更新后添加事件监听器
    setTimeout(() => {
      const githubBtn = document.getElementById('github-btn-error')
      const checkUpdateBtn = document.getElementById('check-update-btn-error')
      if (githubBtn) {
        githubBtn.addEventListener('click', handleGithubClick)
      }
      if (checkUpdateBtn) {
        checkUpdateBtn.addEventListener('click', handleCheckUpdate)
      }
    }, 0)
  }
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
          <el-button class="about-button" @click="showAbout">
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