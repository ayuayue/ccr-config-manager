# Claude Code Router 配置管理器

## 关于本项目

这是一个用于管理和配置 Claude Code Router (CCR) 的图形化配置管理工具。该工具提供了直观的用户界面，方便用户配置和管理 CCR 的各种设置。

CCR GitHub 仓库: [https://github.com/alexrider203/claud-code-router](https://github.com/alexrider203/claud-code-router)

本项目 GitHub 仓库: [https://github.com/ayuayue/ccr-config-manager](https://github.com/ayuayue/ccr-config-manager)


![alt text](PixPin_2025-07-29_22-53-40.png)

## 功能特性

- **系统配置展示**: 显示当前 CCR 的配置信息，包括 API 密钥、代理设置等
- **配置管理**: 完整的配置编辑功能，包括基础配置、提供商配置和路由配置
- **提供商管理**: 添加、编辑和删除 AI 服务提供商及其模型
- **路由配置**: 配置不同类型任务的默认路由，包括后台任务、推理任务、长上下文任务和网络搜索任务
- **响应式设计**: 适配不同屏幕尺寸，提供良好的用户体验

## 技术栈

- 前端: Vue 3 + Vite
- 后端: Go + Wails 框架
- 构建工具: Wails

## 快速开始

### 开发模式

要以开发模式运行，请在项目目录中执行以下命令：

```bash
wails dev
```

这将启动一个 Vite 开发服务器，提供快速的热重载功能。如果你想在浏览器中开发并访问 Go 方法，还有一个运行在 http://localhost:34115 的开发服务器。连接到这个地址后，你可以从开发者工具中调用 Go 代码。

### 构建生产版本

要构建可分发的生产版本，请使用以下命令：

```bash
wails build
```

### 自动化构建脚本

项目提供了便捷的构建脚本：

- **Linux/macOS**: 运行 `./build.sh`
- **Windows**: 运行 `build.bat`

## 发布流程

本项目使用 GitHub Actions 进行自动化构建和发布。当推送带有 `v*.*.*` 格式的标签时，会自动触发构建流程，为以下平台生成发行版：

- Windows (AMD64)
- Linux (AMD64)
- macOS (AMD64 - Intel)
- macOS (ARM64 - Apple Silicon)

构建完成的发行版会自动发布到 GitHub Releases 页面。

## 配置说明

配置文件位于用户目录下的 `~/.claude-code-router/config.json` 文件中。

### 基础配置

- **API Key**: 用于身份验证的密钥（可选）
- **代理 URL**: 为 API 请求设置代理（可选）
- **主机地址**: 设置服务的主机地址（可选）
- **API 超时时间**: API 请求的超时时间（毫秒）
- **启用日志记录**: 是否启用日志记录功能

### 提供商配置

可以配置多个 AI 服务提供商，每个提供商包含：
- 提供商名称
- API 基础 URL
- API 密钥
- 模型列表
- 转换器配置（可选）

### 路由配置

- **默认路由**: 默认使用的提供商和模型
- **后台任务路由**: 后台任务使用的提供商和模型
- **推理任务路由**: 推理任务使用的提供商和模型
- **长上下文路由**: 长上下文任务使用的提供商和模型
- **长上下文阈值**: 触发长上下文路由的上下文长度阈值
- **网络搜索路由**: 网络搜索任务使用的提供商和模型

## 使用方法

1. 启动应用后，你会看到两个主要视图：系统配置和配置管理
2. 在系统配置视图中，可以查看当前的配置信息
3. 在配置管理视图中，可以编辑和保存配置
4. 点击"保存配置"按钮将配置保存到 `~/.claude-code-router/config.json` 文件中

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！