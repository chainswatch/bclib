package parser

import (
  "os"
  "fmt"
)

// blockReader is an interface used to decode blocks and transactions
// it allows to apply the same functions to files and buffers
type Reader interface {
  Peek(int) ([]byte, error)
  ReadByte() byte
  ReadBytes(uint64) []byte
  ReadVarint() uint64
  ReadUint32() uint32
  ReadUint64() uint64
  ReadInt32() int32
}

// file.go
type File struct {
  file      *os.File
  fileNum   uint32
}

// buffer.go
type Buffer struct {
  b      []byte
  pos    uint64
}

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

