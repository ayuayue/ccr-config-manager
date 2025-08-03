//go:build windows

package main

import (
	"syscall"
)

// getSysProcAttr returns SysProcAttr to hide the command window on Windows
func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}