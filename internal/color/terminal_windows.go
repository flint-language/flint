//go:build windows
// +build windows

package color

import (
	"syscall"
	"unsafe"
)

func detectTerminal() bool {
	return enableVirtualTerminalProcessing()
}

func enableVirtualTerminalProcessing() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")

	stdoutHandle := syscall.Handle(syscall.Stdout)
	var mode uint32
	ret, _, _ := getConsoleMode.Call(uintptr(stdoutHandle), uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return false
	}

	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	ret, _, _ = setConsoleMode.Call(uintptr(stdoutHandle), uintptr(mode))
	return ret != 0
}
