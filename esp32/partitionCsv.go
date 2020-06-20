package esp32

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const partitionTableSize = 0x1000
const partitionTableOffset = 0x8000

// PartitionCSVReader reads a CSV file with partition entries
type PartitionCSVReader struct {
	csvReader       *csv.Reader
	lastEnd         int
	partitionsCount int
}

// NewPartitionCSVReader initializes a PartitionCSVReader
func NewPartitionCSVReader(reader io.Reader) *PartitionCSVReader {
	partitionCSVReader := &PartitionCSVReader{
		csvReader: csv.NewReader(reader),
		lastEnd:   partitionTableSize + partitionTableOffset,
	}
	partitionCSVReader.csvReader.Comment = '#'
	return partitionCSVReader
}

// Read tries to read the next partition
func (p *PartitionCSVReader) Read() (*Partition, error) {
	row, err := p.csvReader.Read()
	if err != nil {
		return nil, err
	}

	partition, err := partitionFromRow(row)
	if err != nil {
		return nil, err
	}

	err = p.sanitizePartition(partition)
	if err != nil {
		return nil, err
	}

	return partition, nil
}

func (p *PartitionCSVReader) ReadAll() (partitionList PartitionList, err error) {
	for {
		var partition *Partition
		partition, err = p.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				return
			}
			if strings.HasPrefix(err.Error(), "Invalid row length") {
				continue
			}
			return
		}
		partitionList = append(partitionList, *partition)
	}
	return
}

func partitionFromRow(row []string) (*Partition, error) {
	if len(row) != 6 {
		return nil, fmt.Errorf("Invalid row length %d", len(row))
	}

	for index, column := range row {
		row[index] = os.ExpandEnv(strings.TrimSpace(column))
	}

	partName := row[0]
	partType, err := ParsePartitionType(row[1])
	if err != nil {
		return nil, err
	}
	partSubType, err := ParsePartitionSubType(row[2])
	if err != nil {
		return nil, err
	}
	partOffset, err := parseNumeric(row[3])
	if err != nil {
		return nil, err
	}
	partSize, err := parseNumeric(row[4])
	if err != nil {
		return nil, err
	}
	// TODO partition flags

	return &Partition{
		Name:    partName,
		Type:    partType,
		SubType: partSubType,
		Offset:  partOffset,
		Size:    partSize,
	}, nil
}

// sanitizePartition fixes missing offsets and negative sizes
func (p *PartitionCSVReader) sanitizePartition(partition *Partition) error {
	if partition.Offset != 0 && partition.Offset < p.lastEnd {
		if p.partitionsCount == 0 {
			return fmt.Errorf("First partition overlaps with end of partition table")
		}
		return fmt.Errorf("Partition %d overlaps with previous one", p.partitionsCount)
	}
	p.partitionsCount++

	if partition.Offset == 0 {
		padTo := 0x10000
		if partition.Type == PartitionTypeData {
			padTo = 4
		}
		if p.lastEnd%padTo != 0 {
			p.lastEnd = padTo - (p.lastEnd % padTo)
		}
		partition.Offset = p.lastEnd
	}

	if partition.Size < 0 {
		partition.Size = -partition.Size - partition.Offset
	}

	p.lastEnd = partition.Offset + partition.Size
	return nil
}

// parseNumeric parses generic integer values with provision vor k/m/K/M suffixes and 0x/0b prefixes
func parseNumeric(value string) (int, error) {
	temp := strings.ToLower(value)
	multiplier := 1
	if strings.HasSuffix(temp, "k") {
		multiplier = 1024
		temp = temp[:len(temp)-1]
	}
	if strings.HasSuffix(temp, "m") {
		multiplier = 1024 * 1024
		temp = temp[:len(temp)-1]
	}

	if len(temp) == 0 {
		return 0, nil
	}

	numericValue, err := strconv.ParseInt(temp, 0, 32)
	if err != nil {
		return 0, err
	}

	return int(numericValue) * multiplier, nil
}

type PartitionCSVWriter struct {
	baseWriter io.Writer
	csvWriter  *csv.Writer
}

func NewPartitionCSVWriter(writer io.Writer) *PartitionCSVWriter {
	return &PartitionCSVWriter{
		baseWriter: writer,
		csvWriter:  csv.NewWriter(writer),
	}
}

func (p *PartitionCSVWriter) WriteAll(partitionList PartitionList) error {
	io.WriteString(p.baseWriter, "# name, partition, type, subtype, offset, size, flags\n")
	for _, partition := range partitionList {
		if err := p.Write(&partition); err != nil {
			return err
		}
	}
	p.csvWriter.Flush()
	return nil
}

func (p *PartitionCSVWriter) Write(partition *Partition) error {
	return p.csvWriter.Write([]string{
		partition.Name,
		partition.Type.String(),
		partition.SubType.String(),
		strconv.Itoa(partition.Offset),
		strconv.Itoa(partition.Size),
		"0",
	})
}
