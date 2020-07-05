package esp32

type ChipType byte

const (
	ESP32D0WDQ6 ChipType = 0x00
	ESP32D0WDQ5 ChipType = 0x01
	ESP32D2WDQ5 ChipType = 0x02
	ESP32PicoD4 ChipType = 0x05
)

func (c ChipType) String() string {
	str, known := map[ChipType]string{
		ESP32D0WDQ6: "ESP32D0WDQ6",
		ESP32D0WDQ5: "ESP32D0WDQ5",
		ESP32D2WDQ5: "ESP32D2WDQ5",
		ESP32PicoD4: "ESP32-PICO-D4",
	}[c]
	if known {
		return str
	}
	return "unknown ESP32"
}
