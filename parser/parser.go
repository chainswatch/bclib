package parser

import (
  "encoding/binary"
  "os"
  "fmt"
)

// Reader is an interface used to decode blocks and transactions
// it allows to apply the same functions to files and buffers
type Reader interface {
  Peek(int) ([]byte, error)
  ReadByte() byte
  ReadBytes(uint64) []byte
  ReadVarint() uint64
  ReadUint32() uint32
  ReadUint64() uint64
  ReadInt32() int32

  ReadUint16() uint16
	ShiftVarint() uint64
}

// File allows to use the Reader interface when reading a file
type File struct {
  file      *os.File
  fileNum   uint32
}

// Buffer allows to use the Reader interface when storing data in memory
type Buffer struct {
  b       []byte
  pos     uint64
}

// New allows to declare a new Reader interface from a file or from raw data
func New(x interface{}) Reader {
  switch x.(type) {
  case []byte:
    return &Buffer{x.([]byte), 0}
  case uint32:
    filepath := fmt.Sprintf("/blocks/blk%05d.dat", x.(uint32))
    file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
    if err != nil {
      return nil
    }
    return &File{file: file, fileNum: x.(uint32)}
  default:
    return nil
  }
}

// Varint convert an int to a series of 1 to 8 bytes
func Varint(n uint64) []byte {
	if n > 0xFFFFFFFF {
		val := make([]byte, 8)
		binary.LittleEndian.PutUint64(val, n)
		return append([]byte{0xFF}, val...)
	} else if n > 0xFFF {
		val := make([]byte, 4)
		binary.LittleEndian.PutUint32(val, uint32(n))
		return append([]byte{0xFE}, val...)
	} else if n > 0xFC {
		val := make([]byte, 2)
		binary.LittleEndian.PutUint16(val, uint16(n))
		return append([]byte{0xFD}, val...)
	} else {
		return []byte{byte(n)}
	}
}
