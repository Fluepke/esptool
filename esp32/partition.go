package esp32

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type PartitionType int

const (
	// PartitionTypeApp Application partition type
	PartitionTypeApp PartitionType = iota
	// PartitionTypeData Data partition type
	PartitionTypeData
)

var partitionTypeToString = map[PartitionType]string{
	PartitionTypeApp:  "app",
	PartitionTypeData: "data",
}

var partitionTypeToUint8 = map[PartitionType]uint8{
	PartitionTypeApp:  0x00,
	PartitionTypeData: 0x01,
}

// String returns a string representation of the given PartitionType
func (p PartitionType) String() string {
	name, found := partitionTypeToString[p]
	if found {
		return name
	}
	return strconv.Itoa(int(p))
}

func (p PartitionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// ToUint8 returns a uint8 representation as defined in the ESP-IDF
func (p PartitionType) ToUint8() uint8 {
	value, found := partitionTypeToUint8[p]
	if found {
		return value
	}
	return uint8(p)
}

func PartitionTypeFromUint8(value uint8) PartitionType {
	for partitiontype, rawValue := range partitionTypeToUint8 {
		if rawValue == value {
			return partitiontype
		}
	}
	return PartitionType(value)
}

// ParsePartitionType parses a string into a PartitionType
func ParsePartitionType(value string) (PartitionType, error) {
	if numericValue, err := strconv.ParseInt(value, 0, 32); err == nil {
		return PartitionType(numericValue), nil
	}
	for partitionType, name := range partitionTypeToString {
		if strings.ToLower(name) == value {
			return partitionType, nil
		}
	}
	return PartitionType(0), fmt.Errorf("Illegal partition type '%s'", value)
}

func (p PartitionType) UnmarshalJSON(data []byte) error {
	partitionType, err := ParsePartitionType(string(data))
	p = partitionType
	return err
}

// PartitionSubType as defined in esp_partition_subtype_t
type PartitionSubType int

const (
	// PartitionSubTypeFactory Factory application partition
	PartitionSubTypeFactory PartitionSubType = iota
	// PartitionSubTypeOTA0 OTA partition 0
	PartitionSubTypeOTA0
	// PartitionSubTypeOTA1 OTA partition 1
	PartitionSubTypeOTA1
	// PartitionSubTypeOTA2 OTA partition 2
	PartitionSubTypeOTA2
	// PartitionSubTypeOTA3 OTA partition 3
	PartitionSubTypeOTA3
	// PartitionSubTypeOTA4 OTA partition 4
	PartitionSubTypeOTA4
	// PartitionSubTypeOTA5 OTA partition 5
	PartitionSubTypeOTA5
	// PartitionSubTypeOTA6 OTA partition 6
	PartitionSubTypeOTA6
	// PartitionSubTypeOTA7 OTA partition 7
	PartitionSubTypeOTA7
	// PartitionSubTypeOTA8 OTA partition 8
	PartitionSubTypeOTA8
	// PartitionSubTypeOTA9 OTA partition 9
	PartitionSubTypeOTA9
	// PartitionSubTypeOTA10 OTA partition 10
	PartitionSubTypeOTA10
	// PartitionSubTypeOTA11 OTA partition 11
	PartitionSubTypeOTA11
	// PartitionSubTypeOTA12 OTA partition 12
	PartitionSubTypeOTA12
	// PartitionSubTypeOTA13 OTA partition 13
	PartitionSubTypeOTA13
	// PartitionSubTypeOTA14 OTA partition 14
	PartitionSubTypeOTA14
	// PartitionSubTypeOTA15 OTA partition 15
	PartitionSubTypeOTA15
	// PartitionSubTypeTest Test application partition
	PartitionSubTypeTest
	// PartitionSubTypePHY PHY init data partition
	PartitionSubTypePHY
	// PartitionSubTypeNVS NVS partition
	PartitionSubTypeNVS
	// PartitionSubTypeCoredump COREDUMP partition
	PartitionSubTypeCoredump
	// PartitionSubTypeNvsKeys Partition for NVS keys
	PartitionSubTypeNvsKeys
	// PartitionSubTypeEfuse Partition for emulate eFuse bits
	PartitionSubTypeEfuse
	// PartitionSubTypeEspHttpd ESPHTTPD partition
	PartitionSubTypeEspHttpd
	// PartitionSubTypeFAT FAT partition
	PartitionSubTypeFAT
	// PartitionSubTypeSpiffs SPIFFS partition
	PartitionSubTypeSpiffs
)

var PartitionSubTypeToString = map[PartitionSubType]string{
	PartitionSubTypeFactory:  "factory",
	PartitionSubTypeOTA0:     "ota0",
	PartitionSubTypeOTA1:     "ota1",
	PartitionSubTypeOTA2:     "ota2",
	PartitionSubTypeOTA3:     "ota3",
	PartitionSubTypeOTA4:     "ota4",
	PartitionSubTypeOTA5:     "ota5",
	PartitionSubTypeOTA6:     "ota6",
	PartitionSubTypeOTA7:     "ota7",
	PartitionSubTypeOTA8:     "ota8",
	PartitionSubTypeOTA9:     "ota9",
	PartitionSubTypeOTA10:    "ota10",
	PartitionSubTypeOTA11:    "ota11",
	PartitionSubTypeOTA12:    "ota12",
	PartitionSubTypeOTA13:    "ota13",
	PartitionSubTypeOTA14:    "ota14",
	PartitionSubTypeOTA15:    "ota15",
	PartitionSubTypeTest:     "test",
	PartitionSubTypePHY:      "phy",
	PartitionSubTypeNVS:      "nvs",
	PartitionSubTypeCoredump: "coredump",
	PartitionSubTypeNvsKeys:  "nvs_keys",
	PartitionSubTypeEfuse:    "efuse",
	PartitionSubTypeEspHttpd: "esphttpd",
	PartitionSubTypeFAT:      "fat",
	PartitionSubTypeSpiffs:   "spiffs",
}

var partitionSubTypeToUint8 = map[PartitionSubType]uint8{
	PartitionSubTypeFactory:  0x00,
	PartitionSubTypeOTA0:     0x10,
	PartitionSubTypeOTA1:     0x11,
	PartitionSubTypeOTA2:     0x12,
	PartitionSubTypeOTA3:     0x13,
	PartitionSubTypeOTA4:     0x14,
	PartitionSubTypeOTA5:     0x15,
	PartitionSubTypeOTA6:     0x16,
	PartitionSubTypeOTA7:     0x17,
	PartitionSubTypeOTA8:     0x18,
	PartitionSubTypeOTA9:     0x19,
	PartitionSubTypeOTA10:    0x1a,
	PartitionSubTypeOTA11:    0x1b,
	PartitionSubTypeOTA12:    0x1c,
	PartitionSubTypeOTA13:    0x1d,
	PartitionSubTypeOTA14:    0x1e,
	PartitionSubTypeOTA15:    0x1f,
	PartitionSubTypeTest:     0x20,
	PartitionSubTypePHY:      0x01,
	PartitionSubTypeNVS:      0x02,
	PartitionSubTypeCoredump: 0x03,
	PartitionSubTypeNvsKeys:  0x04,
	PartitionSubTypeEfuse:    0x05,
	PartitionSubTypeEspHttpd: 0x80,
	PartitionSubTypeFAT:      0x81,
	PartitionSubTypeSpiffs:   0x82,
}

// String returns a string representation of the given PartitionSubType
func (p PartitionSubType) String() string {
	name, found := PartitionSubTypeToString[p]

	if found {
		return name
	}
	return strconv.Itoa(int(p))
}

func (p PartitionSubType) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// ToUint8 returns a uint8 representation as defined in the ESP-IDF
func (p PartitionSubType) ToUint8() uint8 {
	value, found := partitionSubTypeToUint8[p]
	if found {
		return value
	}
	return uint8(p)
}

// We don't need to know the partition type here, because no sub types (currently) overlap
func PartitionSubTypeFromUint8(value uint8) PartitionSubType {
	for partitionSubtype, rawValue := range partitionSubTypeToUint8 {
		if rawValue == value {
			return partitionSubtype
		}
	}
	return PartitionSubType(value)
}

// ParsePartitionSubType parses a string into a PartitionSubType
func ParsePartitionSubType(value string) (PartitionSubType, error) {
	if numericValue, err := strconv.ParseInt(value, 0, 32); err == nil {
		return PartitionSubType(numericValue), nil
	}
	for partitionSubType, name := range PartitionSubTypeToString {
		if strings.ToLower(name) == value {
			return partitionSubType, nil
		}
	}

	return PartitionSubType(0), fmt.Errorf("Illegal partition sub type '%s'", value)
}

func (p PartitionSubType) UnmarshalJSON(data []byte) error {
	partitionSubType, err := ParsePartitionSubType(string(data))
	p = partitionSubType
	return err
}

type Partition struct {
	Name    string           `json:"name"`
	Type    PartitionType    `json:"type"`
	SubType PartitionSubType `json:"subtype"`
	Offset  int              `json:"offset'`
	Size    int              `json:"size"`
	//Flags   PartitionFlags   `json:"flags"`
}

func (p *Partition) String() string {
	return fmt.Sprintf("'%-16s' (%8s:%8s) from %6X to %6X", p.Name, p.Type.String(), p.SubType.String(), p.Offset, p.Size)
}

type PartitionList []Partition

func (p PartitionList) String() string {
	builder := &strings.Builder{}
	for _, partition := range p {
		fmt.Println(partition.String())
	}
	return builder.String()
}
