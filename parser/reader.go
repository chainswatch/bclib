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

func NewFile(blockchainDataDir string, fileNum uint32) (*File, error) {
	filepath := fmt.Sprintf(blockchainDataDir + "/blocks/blk%05d.dat", fileNum)
	//fmt.Printf("Opening file %s...\n", filepath)

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &File{file: file, FileNum: fileNum}, nil
}

// TODO: Change *Btc into *file
func (f *File) Close() {
	f.file.Close()
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *File) Size() (int64, error) {
	fInfo, err := f.file.Stat()
	if err != nil {
		return 0, err
	}
	return fInfo.Size(), err
}

/*
* BlockReader interface
*/
func (f *File) Peek(length int) ([]byte, error) {
	pos, err := f.Seek(0, 1)
	if err != nil {
		return nil, err
	}
	val := make([]byte, length)
	f.file.Read(val)
	_, err = f.Seek(pos, 0)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (tx *Buffer) Peek(length int) ([]byte, error) {
  return tx.Body[tx.Pos:(tx.Pos + uint32(length))], nil
}

/*
* BlockReader interface
*/
func (f *File) ReadByte() byte {
	val := make([]byte, 1)
	f.file.Read(val)
	return val[0]
}

func (tx *Buffer) ReadByte() byte {
  tx.Pos++
  return tx.Body[tx.Pos - 1]
}

/*
* BlockReader interface
*/
func (f *File) ReadBytes(length uint64) []byte {
	val := make([]byte, length)
	f.file.Read(val)
	return val
}

func (tx *Buffer) ReadBytes(length uint64) []byte {
  tx.Pos += uint32(length)
  return tx.Body[(tx.Pos - uint32(length)):tx.Pos]
}

func (f *File) ReadUint16() uint16 {
	val := make([]byte, 2)
	f.file.Read(val)
	return binary.LittleEndian.Uint16(val)
}

func (tx *Buffer) ReadUint16() uint16 {
  tx.Pos += 2
  val := tx.Body[(tx.Pos - 2):tx.Pos]
	return binary.LittleEndian.Uint16(val)
}

/*
* BlockReader interface
*/
func (f *File) ReadInt32() int32 {
	raw := make([]byte, 4)
	f.file.Read(raw)
	var val int32
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

func (tx *Buffer) ReadInt32() int32 {
  raw := tx.Body[tx.Pos:(tx.Pos + 4)]
  tx.Pos += 4
	var val int32
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

/*
* BlockReader interface
*/
func (f *File) ReadUint32() uint32 {
	val := make([]byte, 4)
	f.file.Read(val)
	return binary.LittleEndian.Uint32(val)
}

func (tx *Buffer) ReadUint32() uint32 {
  val := tx.Body[tx.Pos:(tx.Pos + 4)]
  tx.Pos += 4
	return binary.LittleEndian.Uint32(val)
}

func (f *File) ReadInt64() int64 {
	raw := make([]byte, 8)
	f.file.Read(raw)
	var val int64
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

/*
* BlockReader interface
*/
func (f *File) ReadUint64() uint64 {
	val := make([]byte, 8)
	f.file.Read(val)
	return binary.LittleEndian.Uint64(val)
}

func (tx *Buffer) ReadUint64() uint64 {
  val := tx.Body[tx.Pos:(tx.Pos + 8)]
  tx.Pos += 8
	return binary.LittleEndian.Uint64(val)
}

//func (f *File) ReadVarint() uint64 {
//	chSize := f.ReadByte()
//	if chSize < 253 {
//		return uint64(chSize)
//	} else if chSize == 253 {
//		return uint64(f.ReadUint16())
//	} else if chSize == 254 {
//		return uint64(f.ReadUint32())
//	} else {
//		return f.ReadUint64()
//	}
//}

/*
* BlockReader interface
*/
func (f *File) ReadVarint() uint64 {
	intType := f.ReadByte()
	if intType == 0xFF {
		return f.ReadUint64()
	} else if intType == 0xFE {
		return uint64(f.ReadUint32())
	} else if intType == 0xFD {
		return uint64(f.ReadUint16())
	}

	return uint64(intType)
}

func (tx *Buffer) ReadVarint() uint64 {
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
