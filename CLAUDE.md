# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Claude Code Router (CCR) configuration manager application built with Go and Vue.js using the Wails framework. It provides a graphical interface for configuring and managing CCR settings.

## Architecture

- **Backend**: Go with Wails framework
- **Frontend**: Vue 3 + Vite
- **Build System**: Wails handles the integration between Go backend and Vue frontend

### Key Components

1. **Main Application** (`main.go`): Entry point that initializes the Wails application
2. **App Structure** (`app.go`): Contains the main App struct with methods for loading/saving CCR configuration
3. **Frontend** (`frontend/src/`):
   - `App.vue`: Main application component with navigation between SystemConfig and ClaudeConfig views
   - `SystemConfig.vue`: Read-only view of current CCR configuration
   - `ClaudeConfig.vue`: Editable configuration management interface

## Common Development Tasks

### Running the Application

**Development Mode**:
```bash
wails dev
```

This starts a Vite development server with hot reload functionality.

### Building the Application

**Production Build**:
```bash
wails build
```

**Platform-specific Build Scripts**:
- Linux/macOS: `./build.sh`
- Windows: `build.bat`

These scripts automatically install dependencies and build the application.

### Dependency Management

- **Go dependencies**: Managed through `go.mod` and `go.sum`
- **Frontend dependencies**: Managed through `frontend/package.json`

To install all dependencies:
```bash
# Backend dependencies
go mod tidy

# Frontend dependencies
cd frontend && npm install
```

## Configuration Structure

The application manages CCR configuration files located at `~/.claude-code-router/config.json` with the following structure:

- Core settings: APIKEY, PROXY_URL, HOST, API_TIMEOUT_MS, LOG
- Providers: Array of AI service providers with name, API URL, key, and models
- Router: Configuration for different task types (default, background, think, longContext, webSearch)