package chains

import(
  "bytes"
  "encoding/binary"
)

func (blockFile *BtcBlock) Close() {
	blockFile.file.Close()
}

func (blockFile *BtcBlock) Seek(offset int64, whence int) (int64, error) {
	return blockFile.file.Seek(offset, whence)
}

func (blockFile *BtcBlock) Size() (int64, error) {
	fileInfo, err := blockFile.file.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), err
}

func (blockFile *BtcBlock) Peek(length int) ([]byte, error) {
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

func (blockFile *BtcBlock) ReadByte() byte {
	val := make([]byte, 1)
	blockFile.file.Read(val)
	return val[0]
}

func (blockFile *BtcBlock) ReadBytes(length uint64) []byte {
	val := make([]byte, length)
	blockFile.file.Read(val)
	return val
}

func (blockFile *BtcBlock) ReadUint16() uint16 {
	val := make([]byte, 2)
	blockFile.file.Read(val)
	return binary.LittleEndian.Uint16(val)
}

func (blockFile *BtcBlock) ReadInt32() int32 {
	raw := make([]byte, 4)
	blockFile.file.Read(raw)
	var val int32
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

func (blockFile *BtcBlock) ReadUint32() uint32 {
	val := make([]byte, 4)
	blockFile.file.Read(val)
	return binary.LittleEndian.Uint32(val)
}

func (blockFile *BtcBlock) ReadInt64() int64 {
	raw := make([]byte, 8)
	blockFile.file.Read(raw)
	var val int64
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

func (blockFile *BtcBlock) ReadUint64() uint64 {
	val := make([]byte, 8)
	blockFile.file.Read(val)
	return binary.LittleEndian.Uint64(val)
}

//func (blockFile *BtcBlock) ReadVarint() uint64 {
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

func (blockFile *BtcBlock) ReadVarint() uint64 {
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
