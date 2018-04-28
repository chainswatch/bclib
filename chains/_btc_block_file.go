// +build ignore
package chains

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

//! Magic numbers to identify start of block
const (
	blockMagicIDBitcoin uint32 = 0xd9b4bef9
	blockMagicIDTestnet uint32 = 0x0709110b
)

type BlockFile struct {
	file    *os.File
	FileNum uint32
}

// NewBlockFile opens (RDONLY) the correct blk.dat file
func NewBlockFile(blockchainDataDir string, fileNum uint32) (*BlockFile, error) {
	filepath := fmt.Sprintf(blockchainDataDir+"/blocks/blk%05d.dat", fileNum)
	//fmt.Printf("Opening file %s...\n", filepath)

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &BlockFile{file: file, FileNum: fileNum}, nil
}

// Close used blk.dat file
func (blockFile *BlockFile) Close() {
	blockFile.file.Close()
}

// Seek ???
func (blockFile *BlockFile) Seek(offset int64, whence int) (int64, error) {
	return blockFile.file.Seek(offset, whence)
}

// Size Get blk.dat file size
func (blockFile *BlockFile) Size() (int64, error) {
	fileInfo, err := blockFile.file.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), err
}

// Peek ???
func (blockFile *BlockFile) Peek(length int) ([]byte, error) {
	pos, err := blockFile.file.Seek(0, 1)
	if err != nil {
		return nil, err
	}
	val := make([]byte, length)
	blockFile.file.Read(val)
	_, err = blockFile.file.Seek(pos, 0)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// ReadByte Read on byte from blk file TODO: Delete
func (blockFile *BlockFile) ReadByte() byte {
	val := make([]byte, 1)
	blockFile.file.Read(val)
	return val[0]
}

// ReadBytes multiple bytes from blk file
func (blockFile *BlockFile) ReadBytes(length uint64) []byte {
	val := make([]byte, length)
	blockFile.file.Read(val)
	return val
}

// ReadUint16 Uint16
func (blockFile *BlockFile) ReadUint16() uint16 {
	val := make([]byte, 2)
	blockFile.file.Read(val)
	return binary.LittleEndian.Uint16(val)
}

// ReadInt32 TODO
func (blockFile *BlockFile) ReadInt32() uint32 {
	raw := make([]byte, 4)
	blockFile.file.Read(raw)
	var val uint32
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

// ReadUint32 TODO
func (blockFile *BlockFile) ReadUint32() uint32 {
	val := make([]byte, 4)
	blockFile.file.Read(val)
	return binary.LittleEndian.Uint32(val)
}

// ReadInt64 TODO
func (blockFile *BlockFile) ReadInt64() int64 {
	raw := make([]byte, 8)
	blockFile.file.Read(raw)
	var val int64
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

// ReadUint64 TODO 
func (blockFile *BlockFile) ReadUint64() uint64 {
	val := make([]byte, 8)
	blockFile.file.Read(val)
	return binary.LittleEndian.Uint64(val)
}

// ReadVarint Select size TODO
func (blockFile *BlockFile) ReadVarint() uint64 {
	intType := blockFile.ReadByte()
	if intType == 0xFF {
		return blockFile.ReadUint64()
	} else if intType == 0xFE {
		return uint64(blockFile.ReadUint32())
	} else if intType == 0xFD {
		return uint64(blockFile.ReadUint16())
	}

	return uint64(intType)
}
