package parser

/*
* Functions needed to read a block from a file
 */

import (
	"bytes"
	"encoding/binary"
)

// Type returns "file"
func (file *File) Type() string {
	return "file"
}

// Close file
func (file *File) Close() {
	file.f.Close()
}

// Reset sets cursor's position to 0
func (file *File) Reset() {
	// TODO: To implement
}

// Seek moves cursor's position to offset
func (file *File) Seek(offset int64, whence int) (int64, error) {
	return file.f.Seek(offset, whence)
}

// Size returns the size of a file
func (file *File) Size() (int64, error) {
	fInfo, err := file.f.Stat()
	if err != nil {
		return 0, err
	}
	return fInfo.Size(), err
}

// Peek read length bytes without moving cursor
func (file *File) Peek(length int) ([]byte, error) {
	pos, err := file.Seek(0, 1)
	if err != nil {
		return nil, err
	}
	val := make([]byte, length)
	file.f.Read(val)
	_, err = file.Seek(pos, 0)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// ReadByte reads next one byte of data
func (file *File) ReadByte() byte {
	val := make([]byte, 1)
	file.f.Read(val)
	return val[0]
}

// ReadBytes reads next length bytes of data
func (file *File) ReadBytes(length uint64) []byte {
	val := make([]byte, length)
	file.f.Read(val)
	return val
}

// ReadUint16 reads next 4 bytes of data as uint16, LE
func (file *File) ReadUint16() uint16 {
	val := make([]byte, 2)
	file.f.Read(val)
	return binary.LittleEndian.Uint16(val)
}

// ReadInt32 reads next 8 bytes of data as int32, LE
func (file *File) ReadInt32() int32 {
	raw := make([]byte, 4)
	file.f.Read(raw)
	var val int32
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

// ReadUint32 reads next 8 bytes of data as uint32, LE
func (file *File) ReadUint32() uint32 {
	val := make([]byte, 4)
	file.f.Read(val)
	return binary.LittleEndian.Uint32(val)
}

// ReadInt64 reads next 16 bytes of data as int64, LE
func (file *File) ReadInt64() int64 {
	raw := make([]byte, 8)
	file.f.Read(raw)
	var val int64
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

// ReadUint64 reads next 16 bytes of data as uint64, LE
func (file *File) ReadUint64() uint64 {
	val := make([]byte, 8)
	file.f.Read(val)
	return binary.LittleEndian.Uint64(val)
}

// ReadCompactSize reads N byte of data as uint64, LE.
// N depends on the first byte
func (file *File) ReadCompactSize() uint64 {
	intType := file.ReadByte()
	if intType == 0xFF {
		return file.ReadUint64()
	} else if intType == 0xFE {
		return uint64(file.ReadUint32())
	} else if intType == 0xFD {
		return uint64(file.ReadUint16())
	}
	return uint64(intType)
}

// ReadVarint does not work for file
// TODO: Implement it
func (file *File) ReadVarint() uint64 {
	return 0xFFFFFF
}
