package parser

import (
  "os"
)

// file.go
type File struct {
  file      *os.File
  FileNum   uint32
}

// buffer.go
type Buffer struct {
  b      []byte
  pos    uint64
}

// TODO: Improve this to use Buffer interface
type DataBuf struct {
	b   []byte
	pos uint64
}
