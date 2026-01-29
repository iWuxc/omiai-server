//go:build linux || bsd || darwin
// +build linux bsd darwin

package main

import (
	"os"
	"syscall"
)

// crashLog .
func crashLog(f *os.File) error {
	if f == nil {
		return nil
	}
	return syscall.Dup2(int(f.Fd()), syscall.Stderr)
}
