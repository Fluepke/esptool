package common

func Uint16ToBytes(value uint16) []byte {
	return []byte{byte(value & 0xFF),
		byte((value >> 8) & 0xFF)}
}

func Uint32ToBytes(value uint32) []byte {
	return []byte{byte(value & 0xFF),
		byte((value >> 8) & 0xFF),
		byte((value >> 16) & 0xFF),
		byte((value >> 24) & 0xFF),
	}
}

func BytesToUint16(value []byte) uint16 {
	return uint16(value[1])<<8 | uint16(value[0])
}
