package parser

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Reader is an interface used to decode blocks and transactions
// it allows to apply the same functions to files and buffers
type Reader interface {
	Type() string
	Peek(int) ([]byte, error)
	Seek(int64, int) (int64, error)
	Reset()
	ReadByte() byte
	ReadBytes(uint64) []byte
	ReadUint32() uint32
	ReadUint64() uint64
	ReadInt32() int32
	ReadVarint() uint64
	ReadCompactSize() uint64

	ReadUint16() uint16
	Close()
}

// File allows to use the Reader interface when reading a file
type File struct {
	f     *os.File
	pos   uint64 // position inside file
	NFile uint32 // file number
}

// Buffer allows to use the Reader interface when storing data in memory
type Buffer struct {
	b   []byte
	pos uint64
}

// New allows to declare a new Reader interface from a file or from raw data
func New(x interface{}) (Reader, error) {
	switch x.(type) {
	case []byte:
		return &Buffer{x.([]byte), 0}, nil
	case uint32:
		dataDir := os.Getenv("DATADIR")
		if dataDir == "" {
			return nil, fmt.Errorf("parser: DATADIR missing")
		}
		filepath := fmt.Sprintf("%s/blocks/blk%05d.dat", dataDir, x.(uint32))
		file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
		if err != nil {
			return nil, err
		}
		return &File{f: file, NFile: x.(uint32)}, nil
	default:
		return nil, fmt.Errorf("parser.New(): Unrecognized input type")
	}
}

// CompactSize convert an int to a series of 1 to 8 bytes
// Used for scriptLength, NVin, NVout, witnessCount
func CompactSize(n uint64) []byte {
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
