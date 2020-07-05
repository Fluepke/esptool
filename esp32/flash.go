package esp32

import (
	"bytes"
	"compress/zlib"
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

func compressImage(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w, err := zlib.NewWriterLevel(&b, 9)
	_, err = w.Write(data)
	w.Close()
	return b.Bytes(), err
}

func (e *ESP32ROM) WriteFlash(offset uint32, data []byte, useCompression bool) (err error) {
	if !e.flashAttached {
		err = e.AttachSpiFlash()
		if err != nil {
			return err
		}
	}

	var remaining []byte

	numBlocks := (uint32(len(data)) + blockLengthWriteMax - 1) / blockLengthWriteMax
	e.logger.Print("Start Erase procedure")

	if useCompression {
		remaining, err = compressImage(data)
		if err != nil {
			return err
		}
		uncompressedNumBlocks := numBlocks
		numBlocks = (uint32(len(remaining)) + blockLengthWriteMax - 1) / blockLengthWriteMax
		e.logger.Printf("Compressed %d bytes to %d bytes. Ration = %.1f", len(data), len(remaining), float64(len(remaining))/float64(len(data)))
		_, err = e.CheckExecuteCommand(
			common.NewBeginFlashDeflCommand(
				uint32(uncompressedNumBlocks)*blockLengthWriteMax,
				uint32(numBlocks),
				blockLengthWriteMax,
				offset,
			),
			10*time.Second,
			e.defaultRetries)
	} else {
		remaining = make([]byte, len(data))
		copy(remaining, data)
		_, err = e.CheckExecuteCommand(
			common.NewBeginFlashCommand(
				uint32(len(data)),
				uint32(numBlocks),
				blockLengthWriteMax,
				offset,
			),
			10*time.Second,
			e.defaultRetries,
		)
	}

	e.logger.Printf("Block size is %d, block count is %d", blockLengthWriteMax, numBlocks)
	if err != nil {
		return err
	}
	e.logger.Print("Begin Flash success.")

	sequence := uint32(0)

	sent := uint32(0)
	total := uint32(len(remaining))

	time.Sleep(10 * time.Millisecond)

	for {
		if sent >= total {
			break
		}
		fmt.Printf("%d of %d - %.2f \n", sent, total, float64(sent)/float64(total)*100.0)

		blockLength := uint32(total - sent)
		if blockLength > blockLengthWriteMax {
			blockLength = blockLengthWriteMax
		}
		block := remaining[sent : sent+blockLength]

		if !useCompression && blockLength < blockLengthWriteMax {
			block = append(block, bytes.Repeat([]byte{0xFF}, int(blockLengthWriteMax-blockLength))...)
		}

		for retryCount := 0; retryCount < 3; retryCount++ {
			if retryCount > 0 {
				e.logger.Printf("Received error while writing to Flash")
			}
			if useCompression {
				_, err = e.CheckExecuteCommand(
					common.NewFlashDataDeflCommand(
						block,
						sequence,
					),
					e.defaultTimeout*100,
					e.defaultRetries,
				)
				if err == nil {
					break
				}
			} else {
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
		}
		if err != nil {
			return err
		}

		sequence++
		sent += blockLength
	}

	//	_, err = e.CheckExecuteCommand(
	//		common.NewFlashEndCommand(false),
	//		e.defaultTimeout,
	//		e.defaultRetries,
	//	)

	return err
}
