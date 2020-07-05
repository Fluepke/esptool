package common

type Opcode byte

const (
	// Opcodes supported by ESP8266 ROM bootloader

	OpcodeFlashBegin Opcode = 0x02
	OpcodeFlashData  Opcode = 0x03
	OpcodeFlashEnd   Opcode = 0x04
	OpcodeMemBegin   Opcode = 0x05
	OpcodeMemEnd     Opcode = 0x06
	OpcodeMemData    Opcode = 0x07
	OpcodeSync       Opcode = 0x08
	OpcodeWriteReg   Opcode = 0x09
	OpcodeReadReg    Opcode = 0x0A

	// Opcodes supported by ESP32 ROM bootloader (or ESP8266 with stub bootloader)

	OpcodeSpiSetParams   Opcode = 0x0B
	OpcodeSpiAttachFlash Opcode = 0x0D
	OpcodeReadFlash      Opcode = 0x0E
	OpcodeChangeBaudrate Opcode = 0x0F
	OpcodeFlashDeflBegin Opcode = 0x10
	OpcodeFlashDeflData  Opcode = 0x11
	OpcodeFlashDeflEnd   Opcode = 0x12
	OpcodeSpiFlashMd5    Opcode = 0x13

	// Opcodes supported by software loader only (ESP8266 & ESP32)

	OpcodeEraseFlash    Opcode = 0xD0
	OpcodeEraseRegion   Opcode = 0xD1
	OpcodeReadFlashFast Opcode = 0xD2
	OpcodeRunUserCode   Opcode = 0xD3
)

func (o Opcode) String() string {
	return map[Opcode]string{
		OpcodeFlashBegin:     "Flash Begin",
		OpcodeFlashData:      "Flash Data",
		OpcodeFlashEnd:       "Flash End",
		OpcodeMemBegin:       "Memory Begin",
		OpcodeMemEnd:         "Memory End",
		OpcodeMemData:        "Memory Data",
		OpcodeSync:           "Sync",
		OpcodeWriteReg:       "Write Register",
		OpcodeReadReg:        "Read Register",
		OpcodeSpiSetParams:   "SPI Set Params",
		OpcodeSpiAttachFlash: "SPI Attach Flash",
		OpcodeReadFlash:      "Read Flash",
		OpcodeChangeBaudrate: "Change Baudrate",
		OpcodeFlashDeflBegin: "Flash Deflate Begin",
		OpcodeFlashDeflData:  "Flash Deflate Data",
		OpcodeFlashDeflEnd:   "Flash Deflate End",
		OpcodeSpiFlashMd5:    "SPI Flash MD5",

		OpcodeEraseFlash:    "Erase Flash (stub bootloader)",
		OpcodeEraseRegion:   "Erase Region (stub bootloader)",
		OpcodeReadFlashFast: "Read Flash Fast (stub bootloader)",
		OpcodeRunUserCode:   "Run User Code (stub bootloader)",
	}[o]
}
