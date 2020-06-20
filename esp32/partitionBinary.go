package esp32

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
)

var partitionMagicBytes = []byte{0xAA, 0x50}
var partitionMD5Begin = []byte{0xEB, 0xEB, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
var partitionTableMaxSize = 0xC00

type PartitionBinaryReader struct {
	reader io.Reader
	md5    hash.Hash
}

func NewPartitionBinaryReader(reader io.Reader) *PartitionBinaryReader {
	md5 := md5.New()
	return &PartitionBinaryReader{
		reader: io.TeeReader(reader, md5),
		md5:    md5,
	}
}

func (p *PartitionBinaryReader) ReadAll() (partitionList PartitionList, err error) {
	for {
		partitionHeader := make([]byte, 2)
		_, err = p.reader.Read(partitionHeader)
		if err != nil {
			return
		}
		if bytes.Equal(partitionHeader, partitionMD5Begin[:2]) {
			// TODO check md5sum
			return
		}
		if !bytes.Equal(partitionHeader, partitionMagicBytes) {
			err = fmt.Errorf("Illegal start of partition header: %v", partitionHeader)
			return
		}

		var partition *Partition
		partition, err = p.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		partitionList = append(partitionList, *partition)
	}
}

func (p *PartitionBinaryReader) Read() (partition *Partition, err error) {

	typeRaw := make([]byte, 1)
	_, err = p.reader.Read(typeRaw)
	if err != nil {
		return
	}

	subTypeRaw := make([]byte, 1)
	_, err = p.reader.Read(subTypeRaw)
	if err != nil {
		return
	}

	var offset uint32
	err = binary.Read(p.reader, binary.LittleEndian, &offset)
	if err != nil {
		return
	}

	var size uint32
	err = binary.Read(p.reader, binary.LittleEndian, &size)
	if err != nil {
		return
	}

	nameRaw := make([]byte, 16)
	_, err = p.reader.Read(nameRaw)
	nameRaw = bytes.Trim(nameRaw, "\x00")

	var flags uint32
	err = binary.Read(p.reader, binary.LittleEndian, &flags)
	if err != nil {
		return
	}

	partition = &Partition{
		Name:    string(nameRaw),
		Type:    PartitionTypeFromUint8(typeRaw[0]),
		SubType: PartitionSubTypeFromUint8(subTypeRaw[0]),
		Offset:  int(offset),
		Size:    int(size),
	}
	return
}

func (p *PartitionBinaryReader) readPartitionHeader() error {
	return nil
}

type PartitionBinaryWriter struct {
	baseWriter io.Writer
}

type CountingWriter struct {
	baseWriter io.Writer
	count      int
}

func (c *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = c.baseWriter.Write(p)
	c.count += len(p)
	return
}

func NewPartitionBinaryWriter(writer io.Writer) *PartitionBinaryWriter {
	return &PartitionBinaryWriter{
		baseWriter: writer,
	}
}

func (p *PartitionBinaryWriter) WriteAll(partitionList PartitionList) error {
	md5 := md5.New()
	countingWriter := CountingWriter{baseWriter: p.baseWriter}
	hashedWriter := io.MultiWriter(md5, &countingWriter)
	for _, partition := range partitionList {
		err := partition.writeBinary(hashedWriter)
		if err != nil {
			return err
		}
	}
	countingWriter.Write(partitionMD5Begin)
	countingWriter.Write(md5.Sum(nil))
	if countingWriter.count > partitionTableMaxSize {
		return fmt.Errorf("Max partition length exceeded. %d > %d", countingWriter.count, partitionTableMaxSize)
	}

	for {
		if countingWriter.count >= partitionTableMaxSize {
			return nil
		}
		if _, err := countingWriter.Write([]byte{0xFF}); err != nil {
			return err
		}
	}
	return nil
}

func (p *Partition) writeBinary(w io.Writer) error {
	w.Write(partitionMagicBytes)
	w.Write([]byte{
		p.Type.ToUint8(),
		p.SubType.ToUint8(),
	})
	binary.Write(w, binary.LittleEndian, uint32(p.Offset))
	binary.Write(w, binary.LittleEndian, uint32(p.Size))
	name := []byte(p.Name)
	if len(name) > 16 {
		name = name[:16]
	}
	w.Write(name)
	w.Write(bytes.Repeat([]byte{0}, 16-len(name)))
	binary.Write(w, binary.LittleEndian, uint32(0))
	return nil
}
