package common

import "fmt"

// ErrorCode - ROM Loader Error Codes
type ErrorCode byte

const (
	// ReceivedMessageInvalid parameters or length field is invalid
	ReceivedMessageInvalid ErrorCode = 0x05
	// FailedToActOnReceivedMessage
	FailedToActOnReceivedMessage ErrorCode = 0x06
	// InvalidCRC Invalid CRC in message
	InvalidCRC ErrorCode = 0x07
	// FlashWriteError - after writing a block of data to flash, the ROM loader reads the value back and the 8-bit CRC is compared to the data read from flash. If they don't match, this error is returned.
	FlashWriteError ErrorCode = 0x08
	// FlashReadError SPI read failed
	FlashReadError ErrorCode = 0x09
	// FlashReadLengthError SPI read request length is too long
	FlashReadLengthError ErrorCode = 0x0A
	// DeflateError (ESP32 compressed uploads only)
	DeflateError ErrorCode = 0x0B
)

// String returns a string representation of the ErrorCode
func (e ErrorCode) String() string {
	str, found := map[ErrorCode]string{
		ReceivedMessageInvalid:       "Received message is invalid",
		FailedToActOnReceivedMessage: "Failed to act on received message",
		InvalidCRC:                   "Invalid CRC in message",
		FlashWriteError:              "Flash write error",
		FlashReadError:               "Flash read error",
		FlashReadLengthError:         "Flash read length error",
		DeflateError:                 "Deflate error",
	}[e]
	if found {
		return str
	}
	return fmt.Sprintf("Unknown error %02X", byte(e))
}
