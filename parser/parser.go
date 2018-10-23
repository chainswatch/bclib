package parser

import (
  "os"
)

type File struct {
  file      *os.File
  FileNum   uint32
}

type Buffer struct {
  Body      []byte
  Pos       uint32
}

// TODO: Improve this to use Buffer interface
type DataBuf struct {
	b   []byte
	pos uint64
}
