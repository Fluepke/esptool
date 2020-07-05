package common

import (
	"bytes"
)

type Command struct {
	Direction Direction
	Opcode    Opcode
	Data      []byte
	Checksum  []byte
}

func (c *Command) ToBytes() []byte {
	b := make([]byte, len(c.Data)+8)
	b[0] = byte(c.Direction)
	b[1] = byte(c.Opcode)
	size := Uint16ToBytes(uint16(len(c.Data)))
	b[2] = size[0]
	b[3] = size[1]
	b[4] = c.Checksum[0]
	b[5] = c.Checksum[1]
	b[6] = c.Checksum[2]
	b[7] = c.Checksum[3]

	for i := 0; i < len(c.Data); i++ {
		b[8+i] = c.Data[i]
	}
	return b
}

func NewCommand(opcode Opcode, data []byte) *Command {
	return &Command{
		Direction: DirectionRequest,
		Opcode:    opcode,
		Data:      data,
		Checksum:  make([]byte, 4),
	}
}

func NewReadRegisterCommand(register uint32) *Command {
	return NewCommand(
		OpcodeReadReg,
		Uint32ToBytes(register),
	)
}

func NewSyncCommand() *Command {
	payload := []byte{0x07, 0x07, 0x12, 0x20}
	payload = append(payload, bytes.Repeat([]byte{0x55}, 32)...)

	return NewCommand(
		OpcodeSync,
		payload,
	)
}

func NewAttachSpiFlashCommand() *Command {
	return NewCommand(OpcodeSpiAttachFlash, make([]byte, 8))
}

func NewReadFlashCommand(offset uint32, size uint32) *Command {
	payload := Uint32ToBytes(offset)
	payload = append(payload, Uint32ToBytes(size)...)

	return NewCommand(
		OpcodeReadFlash,
		payload,
	)
}

func NewChangeBaudrateCommand(newBaudrate uint32, oldBaudrate uint32) *Command {
	payload := Uint32ToBytes(newBaudrate)
	payload = append(payload, Uint32ToBytes(oldBaudrate)...)

	return NewCommand(
		OpcodeChangeBaudrate,
		payload,
	)
}

func NewBeginFlashCommand(eraseSize uint32, numBlocks uint32, blockSize uint32, offset uint32) *Command {
	payload := Uint32ToBytes(eraseSize)
	payload = append(payload, Uint32ToBytes(numBlocks)...)
	payload = append(payload, Uint32ToBytes(blockSize)...)
	payload = append(payload, Uint32ToBytes(offset)...)

	return NewCommand(OpcodeFlashBegin, payload)
}

func NewBeginFlashDeflCommand(eraseSize uint32, numBlocks uint32, blockSize uint32, offset uint32) *Command {
	payload := Uint32ToBytes(eraseSize)
	payload = append(payload, Uint32ToBytes(numBlocks)...)
	payload = append(payload, Uint32ToBytes(blockSize)...)
	payload = append(payload, Uint32ToBytes(offset)...)

	return NewCommand(OpcodeFlashDeflBegin, payload)
}

func calculateChecksum(data []byte) []byte {
	state := uint32(0xEF)

	for _, d := range data {
		state ^= uint32(d)
	}
	return Uint32ToBytes(state)
}

func NewFlashDataCommand(data []byte, sequence uint32) *Command {
	checksum := calculateChecksum(data)
	payload := Uint32ToBytes(uint32(len(data)))
	payload = append(payload, Uint32ToBytes(sequence)...)
	payload = append(payload, Uint32ToBytes(0)...)
	payload = append(payload, Uint32ToBytes(0)...)
	payload = append(payload, data...)

	cmd := NewCommand(OpcodeFlashData, payload)
	cmd.Checksum = checksum

	return cmd
}

func NewFlashDataDeflCommand(data []byte, sequence uint32) *Command {
	checksum := calculateChecksum(data)
	payload := Uint32ToBytes(uint32(len(data)))
	payload = append(payload, Uint32ToBytes(sequence)...)
	payload = append(payload, Uint32ToBytes(0)...)
	payload = append(payload, Uint32ToBytes(0)...)
	payload = append(payload, data...)

	cmd := NewCommand(OpcodeFlashDeflData, payload)
	cmd.Checksum = checksum

	return cmd
}

func NewFlashEndCommand(reboot bool) *Command {
	param := uint32(0)
	if reboot {
		param = 1
	}
	return NewCommand(
		OpcodeFlashEnd,
		Uint32ToBytes(param),
	)
}
