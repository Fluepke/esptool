package esp32

import (
    "net"
)

func (e *ESP32ROM) GetChipMAC() (string, error) {
    buf := make([]byte, 6)
    mac0, err := e.ReadEfuse(2)
    if err != nil {
        return "", err
    }
    mac1, err := e.ReadEfuse(1)
    if err != nil {
        return "", err
    }
    buf[0] = mac0[1]
    buf[1] = mac0[0]
    buf[2] = mac1[3]
    buf[3] = mac1[2]
    buf[4] = mac1[1]
    buf[5] = mac1[0]

    mac := net.HardwareAddr(buf[:])
    return mac.String(), nil
}
