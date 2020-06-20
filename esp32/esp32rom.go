package esp32

import (
	"bytes"
	"fmt"
	"github.com/fluepke/esptool/common"
	"github.com/fluepke/esptool/common/serial"
	"log"
	"time"
)

const (
	efuseRegBase    uint = 0x6001a000
	drRegSysconBase uint = 0x3ff66000
	macEfuseReg     uint = 0x3f41A044 // ESP32-S2 has special block for MAC efuses
)

var (
	flashSizes = map[byte]int{
		0x00: 1048576,  // 1MB
		0x10: 2097152,  // 2MB
		0x20: 4194304,  // 4MB
		0x30: 8388608,  // 8MB
		0x40: 16777216, // 16MB
	}
)

type ESP32ROM struct {
	SerialPort     *serial.Port
	SlipReadWriter *common.SlipReadWriter
	flashAttached  bool
	logger         *log.Logger
	defaultTimeout time.Duration
	defaultRetries int
}

func NewESP32ROM(serialPort *serial.Port, logger *log.Logger) *ESP32ROM {
	return &ESP32ROM{
		SerialPort:     serialPort,
		SlipReadWriter: common.NewSlipReadWriter(serialPort, logger),
		logger:         logger,
		defaultTimeout: 100 * time.Millisecond,
		defaultRetries: 3,
	}
}

func (e *ESP32ROM) Reset() (err error) {
	// set IO0=HIGH
	err = e.SerialPort.SetDTR(false)
	if err != nil {
		return
	}
	// set EN=LOW, chip in reset
	err = e.SerialPort.SetRTS(true)
	if err != nil {
		return
	}

	time.Sleep(100 * time.Millisecond)

	// set IO0=LOW
	err = e.SerialPort.SetDTR(true)
	if err != nil {
		return
	}
	// EN=HIGH, chip out of reset
	err = e.SerialPort.SetRTS(false)

	time.Sleep(5 * time.Millisecond)
	return
}

func (e *ESP32ROM) Connect(maxRetries uint) (err error) {
	err = e.Reset()
	if err != nil {
		return
	}

	err = e.SerialPort.Flush()
	if err != nil {
		return
	}

	for i := uint(0); i < maxRetries; i++ {
		e.logger.Printf("Connecting %d/%d ...\n", i, maxRetries)
		err = e.Sync()
		if err == nil {
			break
		}
	}
	return
}

func (e *ESP32ROM) Sync() (err error) {
	response, err := e.ExecuteCommand(
		common.NewSyncCommand(),
		1000*time.Millisecond,
	)
	if err != nil {
		return err
	}
	if response.Status.Success != true {
		err = fmt.Errorf("Command failed")
	}
	return
}

func (e *ESP32ROM) ReadEfuse(efuseIndex uint) ([4]byte, error) {
	return e.ReadRegister(efuseRegBase + (4 * efuseIndex))
}

func (e *ESP32ROM) ReadRegister(register uint) ([4]byte, error) {
	response, err := e.ExecuteCommand(
		common.NewReadRegisterCommand(uint32(register)),
		e.defaultTimeout,
	)
	if err != nil {
		return [4]byte{}, err
	}
	return response.Value, nil
}

func (e *ESP32ROM) ExecuteCommand(command *common.Command, timeout time.Duration) (*common.Response, error) {
	err := e.SlipReadWriter.Write(command.ToBytes())
	if err != nil {
		return nil, err
	}
	for retryCount := 0; retryCount < 16; retryCount++ {
		responseBuf, err := e.SlipReadWriter.Read(timeout)
		if err != nil {
			return nil, err
		}
		if responseBuf[1] != byte(command.Opcode) {
			e.logger.Printf("Opcode did not match %d/%d\n", retryCount, 16)
			continue
		} else {
			return common.NewResponse(responseBuf)
		}
	}
	return nil, fmt.Errorf("Retrycount exceeded")
}

func (e *ESP32ROM) CheckExecuteCommand(command *common.Command, timeout time.Duration, retries int) (response *common.Response, err error) {
	for retryCount := 0; retryCount < retries; retryCount++ {
		response, err = e.ExecuteCommand(command, timeout)
		if err != nil {
			e.logger.Printf("Executing command %s failed. Retrying %d/%d", command.Opcode.String(), retryCount, retries)
			continue
		}
		if !response.Status.Success {
			err = fmt.Errorf("Device returned for command %s status %s", command.Opcode.String(), response.Status.String())
			e.logger.Printf("Received non success status for command %s. Retrying %d/%d\n", command.Opcode.String(), retryCount, retries)
			continue
		} else {
			break
		}
	}
	return
}

func (e *ESP32ROM) ChangeBaudrate(newBaudrate uint32) error {
	e.logger.Printf("Changing baudrate to %d\n", newBaudrate)
	_, err := e.CheckExecuteCommand(
		common.NewChangeBaudrateCommand(newBaudrate, e.SerialPort.Config.BaudRate),
		e.defaultTimeout,
		e.defaultRetries,
	)
	if err != nil {
		return err
	}

	err = e.SerialPort.SetBaudrate(newBaudrate)
	if err != nil {
		return err
	}

	e.logger.Printf("Changed baudrate to %d", e.SerialPort.Config.BaudRate)
	time.Sleep(10 * time.Millisecond)
	e.SerialPort.Flush() // get rid of crap sent during baud rate change
	return nil
}

func (e *ESP32ROM) ReadPartitionList() (PartitionList, error) {
	e.logger.Print("Reading partiton table from ESP32")

	bindata, err := e.ReadFlash(uint32(partitionTableOffset), uint32(partitionTableMaxSize))

	if err != nil {
		return PartitionList{}, fmt.Errorf("Could not read partition table from chip: %v", err)
	}

	reader := NewPartitionBinaryReader(bytes.NewReader(bindata))

	return reader.ReadAll()
}
