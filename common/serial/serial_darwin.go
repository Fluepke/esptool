// +build darwin

package serial

import (
	"fmt"
	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const maxReadTimeout = 255 * 100 * time.Millisecond

type Port struct {
	file   *os.File
	Config *Config
}

func OpenPort(config *Config) (*Port, error) {
	termiosConfig, err := config.toTermios()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(config.PortPath, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0666)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile errored with: %v", err)
	}

	port := &Port{
		file:   file,
		Config: config,
	}

	if err = termios.Tcsetattr(uintptr(unsafe.Pointer(file.Fd())), termios.TCSANOW, termiosConfig); err != nil {
		return nil, fmt.Errorf("termios.Tcsetattr errored with: %v", err)
	}

	return port, nil
}

func (port *Port) Close() {
	port.file.Close()
}

func (p *Port) SetDTR(dtr bool) error {
	var status int
	err := termios.Tiocmget(p.file.Fd(), &status)
	if err != nil {
		return err
	}
	if dtr {
		status |= unix.TIOCM_DTR
	} else {
		status &^= unix.TIOCM_DTR
	}
	return termios.Tiocmset(p.file.Fd(), &status)
}

func (p *Port) SetRTS(rts bool) error {
	var status int
	err := termios.Tiocmget(p.file.Fd(), &status)
	if err != nil {
		return err
	}
	if rts {
		status |= unix.TIOCM_RTS
	} else {
		status &^= unix.TIOCM_RTS
	}
	return termios.Tiocmset(p.file.Fd(), &status)
}

func (p *Port) SetBaudrate(baudrate uint32) error {
	p.Config.BaudRate = baudrate
	termiosConfig, err := p.Config.toTermios()
	if err != nil {
		return err
	}

	return termios.Tcsetattr(uintptr(p.file.Fd()), termios.TCSADRAIN, termiosConfig)
	// return p.setTermSettings(termiosConfig)
}

func (p *Port) Read(b []byte) (int, error) {
	return p.file.Read(b)
}

func (p *Port) Write(b []byte) (int, error) {
	return p.file.Write(b)
}

func (p *Port) Flush() error {
	return termios.Tcflush(p.file.Fd(), termios.TCIOFLUSH)
}

func getBaudrateFlag(baudrate uint32) (uint32, error) {
	mapping := map[uint32]uint32{
		50:     unix.B50,
		75:     unix.B75,
		110:    unix.B110,
		134:    unix.B134,
		150:    unix.B150,
		200:    unix.B200,
		300:    unix.B300,
		600:    unix.B600,
		1200:   unix.B1200,
		1800:   unix.B1800,
		2400:   unix.B2400,
		4800:   unix.B4800,
		9600:   unix.B9600,
		19200:  unix.B19200,
		38400:  unix.B38400,
		57600:  unix.B57600,
		115200: unix.B115200,
		230400: unix.B230400,
		// 460800:  unix.B460800,
		// 500000:  unix.B500000,
		// 576000:  unix.B576000,
		// 921600:  unix.B921600,
		// 1000000: unix.B1000000,
		// 1152000: unix.B1152000,
		// 1500000: unix.B1500000,
		// 2000000: unix.B2000000,
		// 2500000: unix.B2500000,
		// 3000000: unix.B3000000,
		// 3500000: unix.B3500000,
		// 4000000: unix.B4000000,
	}
	return 0, nil
	value, found := mapping[baudrate]
	fmt.Printf("value=%d\n", value)
	if !found {
		return 0, fmt.Errorf("Baudrate %d not supported. Please choose POSIX compliant baudrate.", baudrate)
	}
	return value, nil
}

// getOutputModes returns the termios flags to use for the given configuration
func getOutputModes(config *Config) (cFlag uint32, err error) {
	baudrateFlag, err := getBaudrateFlag(config.BaudRate)
	if err != nil {
		return 0, err
	}
	cFlag = unix.CREAD | unix.CLOCAL | baudrateFlag

	switch config.DataBits {
	case 5:
		cFlag |= unix.CS5
	case 6:
		cFlag |= unix.CS6
	case 7:
		cFlag |= unix.CS7
	case 8:
		cFlag |= unix.CS8
	default:
		return 0, fmt.Errorf("Bad data bits value %d", config.DataBits)
	}

	switch config.StopBits {
	case StopBitsOne:
		// default
	case StopBitsTwo:
		cFlag |= unix.CSTOPB
	default:
		return 0, fmt.Errorf("Bad stop bits value")
	}

	fmt.Printf("cFlag = %X\n", cFlag)
	return cFlag, nil
}

func (c *Config) toTermios() (*syscall.Termios, error) {
	cFlag, err := getOutputModes(c)
	if err != nil {
		return nil, err
	}

	termios := &syscall.Termios{
		//		Iflag: unix.INPCK,
		Cflag: uint64(cFlag),
	}

	termios.Cc[unix.VMIN] = 1

	if c.ReadTimeout > 0 {
		termios.Cc[unix.VMIN] = 0
		if c.ReadTimeout > maxReadTimeout {
			return nil, fmt.Errorf("Bad read timeout")
		}
		termios.Cc[unix.VTIME] = uint8(c.ReadTimeout.Milliseconds() / 100)
	}

	return termios, nil
}
