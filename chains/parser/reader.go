package parser

/*
* Functions needed to read a block from a file
*/

import (
  "bytes"
  "fmt"
  "encoding/binary"
  "os"
)

type BlockFile struct {
  file      *os.File
  FileNum   uint32
}

type RawTx struct {
  Body      []byte
  Pos       uint32
}

func NewBlockFile(blockchainDataDir string, fileNum uint32) (*BlockFile, error) {
	filepath := fmt.Sprintf(blockchainDataDir + "/blocks/blk%05d.dat", fileNum)
	//fmt.Printf("Opening file %s...\n", filepath)

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &BlockFile{file: file, FileNum: fileNum}, nil
}

// TODO: Change *Btc into *file
func (blockFile *BlockFile) Close() {
	blockFile.file.Close()
}

func (blockFile *BlockFile) Seek(offset int64, whence int) (int64, error) {
	return blockFile.file.Seek(offset, whence)
}

func (blockFile *BlockFile) Size() (int64, error) {
	blockFileInfo, err := blockFile.file.Stat()
	if err != nil {
		return 0, err
	}
	return blockFileInfo.Size(), err
}

/*
* BlockReader interface
*/
func (blockFile *BlockFile) Peek(length int) ([]byte, error) {
	pos, err := blockFile.Seek(0, 1)
	if err != nil {
		return nil, err
	}
	val := make([]byte, length)
	blockFile.file.Read(val)
	_, err = blockFile.Seek(pos, 0)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (tx *RawTx) Peek(length int) ([]byte, error) {
  return tx.Body[tx.Pos:(tx.Pos + uint32(length))], nil
}

/*
* BlockReader interface
*/
func (blockFile *BlockFile) ReadByte() byte {
	val := make([]byte, 1)
	blockFile.file.Read(val)
	return val[0]
}

func (tx *RawTx) ReadByte() byte {
  tx.Pos += 1
  return tx.Body[tx.Pos - 1]
}

/*
* BlockReader interface
*/
func (blockFile *BlockFile) ReadBytes(length uint64) []byte {
	val := make([]byte, length)
	blockFile.file.Read(val)
	return val
}

func (tx *RawTx) ReadBytes(length uint64) []byte {
  tx.Pos += uint32(length)
  return tx.Body[(tx.Pos - uint32(length)):tx.Pos]
}

func (blockFile *BlockFile) ReadUint16() uint16 {
	val := make([]byte, 2)
	blockFile.file.Read(val)
	return binary.LittleEndian.Uint16(val)
}

func (tx *RawTx) ReadUint16() uint16 {
  tx.Pos += 2
  val := tx.Body[(tx.Pos - 2):tx.Pos]
	return binary.LittleEndian.Uint16(val)
}

/*
* BlockReader interface
*/
func (blockFile *BlockFile) ReadInt32() int32 {
	raw := make([]byte, 4)
	blockFile.file.Read(raw)
	var val int32
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

func (tx *RawTx) ReadInt32() int32 {
  raw := tx.Body[tx.Pos:(tx.Pos + 4)]
  tx.Pos += 4
	var val int32
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

/*
* BlockReader interface
*/
func (blockFile *BlockFile) ReadUint32() uint32 {
	val := make([]byte, 4)
	blockFile.file.Read(val)
	return binary.LittleEndian.Uint32(val)
}

func (tx *RawTx) ReadUint32() uint32 {
  val := tx.Body[tx.Pos:(tx.Pos + 4)]
  tx.Pos += 4
	return binary.LittleEndian.Uint32(val)
}

func (blockFile *BlockFile) ReadInt64() int64 {
	raw := make([]byte, 8)
	blockFile.file.Read(raw)
	var val int64
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

/*
* BlockReader interface
*/
func (blockFile *BlockFile) ReadUint64() uint64 {
	val := make([]byte, 8)
	blockFile.file.Read(val)
	return binary.LittleEndian.Uint64(val)
}

func (tx *RawTx) ReadUint64() uint64 {
  val := tx.Body[tx.Pos:(tx.Pos + 8)]
  tx.Pos += 8
	return binary.LittleEndian.Uint64(val)
}

//func (blockFile *BlockFile) ReadVarint() uint64 {
//	chSize := blockFile.ReadByte()
//	if chSize < 253 {
//		return uint64(chSize)
//	} else if chSize == 253 {
//		return uint64(blockFile.ReadUint16())
//	} else if chSize == 254 {
//		return uint64(blockFile.ReadUint32())
//	} else {
//		return blockFile.ReadUint64()
//	}
//}

/*
* BlockReader interface
*/
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

func (tx *RawTx) ReadVarint() uint64 {
	intType := tx.ReadByte()
	if intType == 0xFF {
		return tx.ReadUint64()
	} else if intType == 0xFE {
		return uint64(tx.ReadUint32())
	} else if intType == 0xFD {
		return uint64(tx.ReadUint16())
	}

	return uint64(intType)
}
