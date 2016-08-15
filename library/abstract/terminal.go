package abstract

import (
    "syscall"
    "unsafe"
)

type TerminalInfo struct {
	winsize *winsize
}

func NewTerminalInfo() *TerminalInfo {
	return &TerminalInfo{}
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func (terminal *TerminalInfo) Width() int {
	if terminal.winsize == nil {
		terminal.winsize = &winsize{}
		retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(terminal.winsize)))

		if int(retCode) == -1 {
			panic(errno) //TODO
		}
	}

	return int(terminal.winsize.Col)
}
