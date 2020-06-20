package serial

import "time"

type StopBits   byte

const (
    StopBitsOne     StopBits = iota
    StopBitsOneHalf
    StopBitsTwo
)

type Parity     byte

const (
    ParityNone  Parity = iota
    ParityOdd
    ParityEven
    ParityMark
    ParitySpace
)

type Config struct {
    PortPath    string
    BaudRate    uint32
    ReadTimeout time.Duration
    DataBits    byte
    StopBits    StopBits
    Parity      Parity
}

func NewConfig(portPath string, baudrate uint32) *Config {
    return &Config{
        PortPath: portPath,
        BaudRate: baudrate,
        ReadTimeout: 1 * time.Millisecond,
        DataBits: 8,
        StopBits: StopBitsOne,
        Parity: ParityNone,
    }
}
