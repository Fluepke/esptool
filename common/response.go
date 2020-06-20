package common

import (
	"fmt"
)

const (
	minResponseSize    int = 10
	responseStatusSize int = 2
)

type ResponseStatus struct {
	Success   bool
	ErrorCode ErrorCode
}

type Response struct {
	Direction Direction
	Opcode    Opcode
	Size      uint16
	Value     [4]byte
	Data      []byte
	Status    *ResponseStatus
}

func (r *ResponseStatus) String() string {
	return fmt.Sprintf("success = %t, error_code = %s", r.Success, r.ErrorCode.String())
}

func NewResponseStatus(data []byte) (*ResponseStatus, error) {
	if len(data) != responseStatusSize {
		return nil, fmt.Errorf("Invalid response status length. Received %d bytes, expected exactly %d bytes", len(data), responseStatusSize)
	}
	return &ResponseStatus{
		Success:   data[0] == 0,
		ErrorCode: ErrorCode(data[1]),
	}, nil
}

func NewResponse(data []byte) (*Response, error) {
	if len(data) < minResponseSize {
		return nil, fmt.Errorf("Invalid response length. Received %d bytes, expected at least %d bytes", len(data), minResponseSize)
	}
	response := &Response{
		Direction: Direction(data[0]),
		Opcode:    Opcode(data[1]),
		Size:      BytesToUint16(data[2:4]),
		Data:      data[8 : len(data)-1],
	}
	for i := 0; i < 4; i++ {
		response.Value[i] = data[4+i]
	}
	status, err := NewResponseStatus(data[len(data)-2 : len(data)])
	if err != nil {
		return nil, err
	}
	response.Status = status

	return response, nil
}
