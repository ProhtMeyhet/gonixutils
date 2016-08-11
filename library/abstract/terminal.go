package abstract

import (
    "syscall"
    "unsafe"
)

type Terminal struct {

}

func NewTerminal() *Terminal {
	return &Terminal{}
}

type winsize struct {
    Row    uint16
    Col    uint16
    Xpixel uint16
    Ypixel uint16
}

func (terminal *Terminal) GetWidth() int {
    ws := &winsize{}
    retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
        uintptr(syscall.Stdin),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(ws)))

    if int(retCode) == -1 {
        panic(errno) //TODO
    }
    return int(ws.Col)
}
