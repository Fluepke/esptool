package common

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"
)

type SlipReadWriter struct {
	BaseReadWriter io.ReadWriter
	Timeout        time.Duration
	logger         *log.Logger
}

func NewSlipReadWriter(base io.ReadWriter, logger *log.Logger) *SlipReadWriter {
	return &SlipReadWriter{
		BaseReadWriter: base,
		logger:         logger,
	}
}

const (
	SlipHeader     byte = 0xC0
	SlipEscapeChar byte = 0xDB
)

func SlipEncode(b []byte) []byte {
	escapeCharsReplaced := bytes.ReplaceAll(b, []byte{SlipEscapeChar}, []byte{SlipEscapeChar, 0xDD})
	headersReplaced := bytes.ReplaceAll(escapeCharsReplaced, []byte{SlipHeader}, []byte{SlipEscapeChar, 0xDC})

	result := append([]byte{SlipHeader}, headersReplaced...)
	return append(result, SlipHeader)
}

func (s *SlipReadWriter) Write(b []byte) error {
	data := SlipEncode(b)
	n, err := s.BaseReadWriter.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		err := fmt.Errorf("Expected to send %d bytes but transfered only %d bytes.", len(data), n)
		s.logger.Print(err.Error())
		return err
	}
	return nil
}

func (s *SlipReadWriter) Read(timeout time.Duration) ([]byte, error) {
	// read a single byte until we find a slip header
	byteBuf := make([]byte, 1)
	startTime := time.Now()

	// simple state machine
	type slipReadState byte
	const waitingForHeader slipReadState = 0
	const readingContent slipReadState = 1
	const inEscape slipReadState = 2
	var state slipReadState = waitingForHeader

	result := make([]byte, 0)

	for {
		if time.Since(startTime) > timeout {
			err := fmt.Errorf("Read timeout after %v. Received %d bytes", time.Since(startTime), len(result))
			s.logger.Print(err)
			return nil, err
		}
		n, err := s.BaseReadWriter.Read(byteBuf)
		if err != nil {
			if err.Error() == "EOF" {
				continue
			}
			return nil, err
		}
		if n != 1 {
			continue
		}

		switch state {
		case waitingForHeader:
			if byteBuf[0] == SlipHeader {
				state = readingContent
			}
		case readingContent:
			switch byteBuf[0] {
			case SlipHeader:
				return result, nil
			case SlipEscapeChar:
				state = inEscape
			default:
				result = append(result, byteBuf[0])
			}
		case inEscape:
			switch byteBuf[0] {
			case 0xDC:
				result = append(result, SlipHeader)
				state = readingContent
			case 0xDD:
				result = append(result, SlipEscapeChar)
				state = readingContent
			default:
				return nil, fmt.Errorf("Unexpected char %02X after escape character", byteBuf[0])
			}
		}
	}
}
