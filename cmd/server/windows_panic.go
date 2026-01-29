//go:build windows
// +build windows

package main

import (
	"github.com/pkg/errors"
	"os"
	"syscall"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

// crashLog
func crashLog(f *os.File) error {
	if f == nil {
		return nil
	}

	err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		return errors.Errorf("Failed to redirect stderr to file: %v", err)
	}
	os.Stderr = f
	return nil
}

func setStdHandle(stdHandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdHandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}
