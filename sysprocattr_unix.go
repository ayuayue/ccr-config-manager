//go:build !windows

package main

import (
	"syscall"
)

// getSysProcAttr returns SysProcAttr for Unix-like systems (Linux, macOS)
func getSysProcAttr() *syscall.SysProcAttr {
	// Unix-like systems don't need HideWindow attribute
	return nil
}