package main

import (
	"fmt"
	"github.com/fluepke/esptool/common/serial"
	"github.com/fluepke/esptool/esp32"
	"log"
)

func bold(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}

func underline(s string) string {
	return fmt.Sprintf("\033[4m%s\033[0m", s)
}

func connectEsp32(portPath string, connectBaudrate uint32, transferBaudrate uint32, retries uint, logger *log.Logger) (*esp32.ESP32ROM, error) {
	serialConfig := serial.NewConfig(portPath, connectBaudrate)
	serialPort, err := serial.OpenPort(serialConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to open serial port: %s", err.Error())
	}
	esp32 := esp32.NewESP32ROM(serialPort, logger)
	err = esp32.Connect(retries)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ESP32: %s", err.Error())
	}
	return esp32, esp32.ChangeBaudrate(transferBaudrate)
}
