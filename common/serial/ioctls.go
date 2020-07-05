// +build !windows,cgo

package serial

import (
	"golang.org/x/sys/unix"
	"unsafe"
)

func (p *Port) ioctl(req uint64, data uintptr) (err error) {
	_, _, e1 := unix.Syscall(unix.SYS_IOCTL, uintptr(p.file.Fd()), uintptr(req), data)
	if e1 != 0 {
		err = e1
	}
	return
}

func (p *Port) setTermSettings(settings *unix.Termios) error {
	return p.ioctl(unix.TCSETS, uintptr(unsafe.Pointer(settings)))
}

func (p *Port) setModemBitsStatus(status int) error {
	return p.ioctl(unix.TIOCMSET, uintptr(unsafe.Pointer(&status)))
}

func (p *Port) getModemBitsStatus() (int, error) {
	var status int
	err := p.ioctl(unix.TIOCMGET, uintptr(unsafe.Pointer(&status)))
	return status, err
}
