package esp32

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strings"
	"testing"
)

var (
	csv1 = `# Name,   Type, SubType, Offset,  Size, Flags
# Note: if you have increased the bootloader size, make sure to update the offsets to avoid overlap
nvs,      data, nvs,     ,        0x6000,
phy_init, data, phy,     ,        0x1000,
factory,  app,  factory, ,        1M,`

	desired1 = PartitionList{
		Partition{
			Name:    "nvs",
			Type:    PartitionTypeData,
			SubType: PartitionSubTypeNVS,
			Offset:  0x9000,
			Size:    0x6000,
		},
		Partition{
			Name:    "phy_init",
			Type:    PartitionTypeData,
			SubType: PartitionSubTypePHY,
			Offset:  0xF000,
			Size:    0x1000,
		},
		Partition{
			Name:    "factory",
			Type:    PartitionTypeApp,
			SubType: PartitionSubTypeFactory,
			Offset:  0x10000,
			Size:    1024 * 1024,
		},
	}
	md5sum1 = "5d61d196adc3dba01928f264eb169be7"
)

func assertPartition(t *testing.T, desired *Partition, received *Partition) {
	if received.Name != desired.Name {
		t.Errorf("Expected partition name '%s', received '%s'", desired.Name, received.Name)
	}
	if received.Type != desired.Type {
		t.Errorf("Expected partition type %v, received %v", desired.Type, received.Type)
	}
	if received.SubType != desired.SubType {
		t.Errorf("Expected partition sub type %v, received %v", desired.SubType, received.SubType)
	}
	if received.Offset != desired.Offset {
		t.Errorf("Expected partition offset %d, received %d", desired.Offset, received.Offset)
	}
	if received.Size != desired.Size {
		t.Errorf("Expected partition size %d, received %d", desired.Size, received.Size)
	}
}

func assertPartitionList(t *testing.T, desired PartitionList, received PartitionList) {
	if len(desired) != len(received) {
		t.Errorf("Expected partition list of length %d, received %d", len(desired), len(received))
	}
	for index, receivedPartition := range received {
		assertPartition(t, &desired[index], &receivedPartition)
	}
}

func assertReadCSV(t *testing.T, csv string, desired PartitionList) {
	partitionCSVReader := NewPartitionCSVReader(strings.NewReader(csv))

	partitions, err := partitionCSVReader.ReadAll()
	if err != nil {
		t.Errorf("partitionCSVReader.ReadAll errored with: %v", err)
	}

	assertPartitionList(t, desired, partitions)
}

func TestReadCSV(t *testing.T) {
	assertReadCSV(t, csv1, desired1)
}

func TestWriteCSV(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	partitionsCSVWriter := NewPartitionCSVWriter(buf)
	if err := partitionsCSVWriter.WriteAll(desired1); err != nil {
		t.Errorf("partitionsCSVWriter.WriteAll errored with: %v", err)
	}

	partitionCSVReader := NewPartitionCSVReader(bytes.NewReader(buf.Bytes()))
	partitionList, err := partitionCSVReader.ReadAll()
	if err != nil {
		t.Errorf("partitionCSVReader.ReadAll errored with: %v", err)
	}
	assertPartitionList(t, desired1, partitionList)
}

func TestWriteBinary(t *testing.T) {
	hash := md5.New()
	writer := NewPartitionBinaryWriter(hash)
	if err := writer.WriteAll(desired1); err != nil {
		t.Errorf("Failed to write partition list: %v", err)
	}
	md5sum := fmt.Sprintf("%x", hash.Sum(nil))
	if md5sum != md5sum1 {
		t.Errorf("Checksum did not match %s", md5sum)
	}
}

func TestReadBinary(t *testing.T) {
	reader := NewPartitionBinaryReader(bytes.NewReader(binary1))
	partitionList, err := reader.ReadAll()
	if err != nil {
		t.Errorf("Got unexpected error while reading partition table: %v", err)
	}
	assertPartitionList(t, desired1, partitionList)
}
