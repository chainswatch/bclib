package parser

import (
  "bytes"
  "encoding/binary"
)

func NewBuffer(b []byte) *Buffer {
	return &Buffer{b, 0}
}

func (buf *Buffer) Reset() {
	buf.pos = 0
}

func (buf *Buffer) Peek(length int) ([]byte, error) {
  return buf.b[buf.pos:(buf.pos + uint64(length))], nil
}

func (buf *Buffer) ReadByte() byte {
  buf.pos++
  return buf.b[buf.pos - 1]
}

func (buf *Buffer) ReadBytes(length uint64) []byte {
  buf.pos += length
  return buf.b[(buf.pos - length):buf.pos]
}

func (buf *Buffer) ReadUint16() uint16 {
  buf.pos += 2
  val := buf.b[(buf.pos - 2):buf.pos]
	return binary.LittleEndian.Uint16(val)
}

func (buf *Buffer) ReadInt32() int32 {
  raw := buf.b[buf.pos:(buf.pos + 4)]
  buf.pos += 4
	var val int32
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

func (buf *Buffer) ReadUint32() uint32 {
  val := buf.b[buf.pos:(buf.pos + 4)]
  buf.pos += 4
	return binary.LittleEndian.Uint32(val)
}

func (buf *Buffer) ReadUint64() uint64 {
  val := buf.b[buf.pos:(buf.pos + 8)]
  buf.pos += 8
	return binary.LittleEndian.Uint64(val)
}

func (buf *Buffer) ReadVarint() uint64 {
	intType := buf.ReadByte()
	if intType == 0xFF {
		return buf.ReadUint64()
	} else if intType == 0xFE {
		return uint64(buf.ReadUint32())
	} else if intType == 0xFD {
		return uint64(buf.ReadUint16())
	}

	return uint64(intType)
}

/*
 * WARNING: Methods listed below belonger to DataBuf
 * They have been assigned to Buffer
 * Which may cause some unknown behavior
*/

func (buf *Buffer) Seek(pos uint64) {
	buf.pos = pos
}

func (buf *Buffer) ShiftByte() byte {
	val := buf.b[buf.pos : buf.pos+1]
	buf.pos++
	return val[0]
}

func (buf *Buffer) ShiftBytes(length uint64) []byte {
	val := buf.b[buf.pos : buf.pos+length]
	buf.pos += length
	return val
}

func (buf *Buffer) Shift16bit() uint16 {
	val := binary.LittleEndian.Uint16(buf.b[buf.pos : buf.pos+2])
	buf.pos += 2
	return val
}

func (buf *Buffer) ShiftU64bit() uint64 {
	val := binary.LittleEndian.Uint64(buf.b[buf.pos : buf.pos+8])
	buf.pos += 8
	return val
}

func (buf *Buffer) Shift64bit() int64 {
	val := binary.LittleEndian.Uint64(buf.b[buf.pos : buf.pos+8])
	buf.pos += 8
	return int64(val)
}

func (buf *Buffer) ShiftU32bit() uint32 {
	val := binary.LittleEndian.Uint32(buf.b[buf.pos : buf.pos+4])
	buf.pos += 4
	return val
}

func (buf *Buffer) Shift32bit() int32 {
	val := binary.LittleEndian.Uint32(buf.b[buf.pos : buf.pos+4])
	buf.pos += 4
	return int32(val)
}

func (buf *Buffer) ShiftVarint() uint64 {
	var n uint64
	for true {
		b := buf.b[buf.pos : buf.pos+1][0]
		buf.pos++
		n = (n << uint64(7)) | uint64(b&uint8(0x7F))
		if b&uint8(0x80) > 0 {
			n++
		} else {
			return n
		}
	}

	return n
}
