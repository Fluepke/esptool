package esp32

import (
	"bytes"
	"fmt"
	"github.com/fluepke/esptool/common"
	"time"
)

const blockLengthReadMax uint32 = 64 // TODO check if this value taken from the esptool.py is really true
const blockLengthWriteMax uint32 = 0x400

func (e *ESP32ROM) AttachSpiFlash() (err error) {
	_, err = e.CheckExecuteCommand(
		common.NewAttachSpiFlashCommand(),
		e.defaultTimeout,
		e.defaultRetries,
	)
	if err != nil {
		return err
	}
	e.flashAttached = true
	e.logger.Print("Attach SPI flash success")
	return
}

// func (e *ESP32ROM) SpiSetParams() (err error) {
// 	_, err = e.CheckExecuteCommand(
// 		common.NewSpiSetParamsCommand(
// 			uint32(0),
// 			e.Size(

func (e *ESP32ROM) ReadFlash(offset uint32, size uint32) ([]byte, error) {
	if !e.flashAttached {
		err := e.AttachSpiFlash()
		if err != nil {
			return []byte{}, err
		}
	}

	receivedData := make([]byte, 0)
	for {
		// e.logger.Printf("%d of %d\n", len(receivedData), size)
		if len(receivedData) >= int(size) {
			return receivedData, nil
		}

		blockLength := size - uint32(len(receivedData))
		if blockLength > blockLengthReadMax {
			blockLength = blockLengthReadMax
		}

		response, err := e.CheckExecuteCommand(
			common.NewReadFlashCommand(offset+uint32(len(receivedData)), blockLength),
			e.defaultTimeout,
			e.defaultRetries,
		)
		if err != nil {
			return receivedData, err
		}

		receivedData = append(receivedData, response.Data[:blockLength]...)
	}
}

func (e *ESP32ROM) WriteFlash(offset uint32, data []byte) error {
	if !e.flashAttached {
		err := e.AttachSpiFlash()
		if err != nil {
			return err
		}
	}

	remaining := make([]byte, len(data))
	copy(remaining, data)

	numBlocks := (uint32(len(data)) + blockLengthWriteMax - 1) / blockLengthWriteMax
	e.logger.Print("Start Erase procedure")
	_, err := e.CheckExecuteCommand(
		common.NewBeginFlashCommand(
			uint32(len(data)),
			uint32(numBlocks),
			blockLengthWriteMax,
			offset,
		),
		10*time.Second,
		e.defaultRetries,
	)
	e.logger.Print("Begin Flash success.")
	e.logger.Printf("Block size is %d, block count is %d", blockLengthWriteMax, numBlocks)
	if err != nil {
		return err
	}

	sequence := uint32(0)

	sent := uint32(0)
	total := uint32(len(data))

	time.Sleep(100 * time.Millisecond)

	for {
		if sent >= total {
			break
		}
		fmt.Printf("%d of %d - %.2f \n", sent, total, float64(sent)/float64(total)*100.0)

		blockLength := uint32(total - sent)
		if blockLength > blockLengthWriteMax {
			blockLength = blockLengthWriteMax
		}
		block := remaining[sent : sent+blockLength] // TODO we might need to pad the last block

		if blockLength < blockLengthWriteMax {
			block = append(block, bytes.Repeat([]byte{0xFF}, int(blockLengthWriteMax-blockLength))...)
		}

		for retryCount := 0; retryCount < 3; retryCount++ {
			if retryCount > 0 {
				e.logger.Printf("Received error while writing to Flash: %s", err.Error())
			}
			_, err = e.CheckExecuteCommand(
				common.NewFlashDataCommand(
					block,
					sequence,
				),
				e.defaultTimeout,
				e.defaultRetries,
			)

			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}

		sequence++
		sent += blockLength
	}

	_, err = e.CheckExecuteCommand(
		common.NewFlashEndCommand(true),
		e.defaultTimeout,
		e.defaultRetries,
	)

	return err
}
