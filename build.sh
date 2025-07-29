#!/bin/bash

# Build script for Claude Config Manager

# Check if wails is installed
if ! command -v wails &> /dev/null
then
    echo "Wails is not installed. Installing..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
fi

# Install frontend dependencies
echo "Installing frontend dependencies..."
cd frontend
npm install
cd ..

# Build for current platform
echo "Building for current platform..."
wails build

echo "Build completed successfully!"