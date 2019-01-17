package parser

import (
	"encoding/binary"
)

// Type returns type of reader
func (buf *Buffer) Type() string {
	return "buffer"
}

// Reset cursor to position 0
func (buf *Buffer) Reset() {
	buf.pos = 0
}

// Close buffer
// TODO: Is it relevant for Buffer?
func (buf *Buffer) Close() {
}

// Peek up to length without moving cursor
func (buf *Buffer) Peek(length int) ([]byte, error) {
	return buf.b[buf.pos:(buf.pos + uint64(length))], nil
}

// Seek moves cursor tu position pos
func (buf *Buffer) Seek(pos int64, whence int) (int64, error) {
	switch whence {
	case 0:
		buf.pos = uint64(pos)
	case 1:
		buf.pos += uint64(pos)
		// TODO: case 2
	}
	return pos, nil
}

// ReadByte reads next one byte of data
func (buf *Buffer) ReadByte() byte {
	val := buf.b[buf.pos : buf.pos+1]
	buf.pos++
	return val[0]
}

// ReadBytes reads next length bytes of data
func (buf *Buffer) ReadBytes(length uint64) []byte {
	val := buf.b[buf.pos : buf.pos+length]
	buf.pos += length
	return val
}

// ReadUint16 reads next 4 bytes of data as uint16, LE
func (buf *Buffer) ReadUint16() uint16 {
	val := binary.LittleEndian.Uint16(buf.b[buf.pos : buf.pos+2])
	buf.pos += 2
	return val
}

// ReadInt32 reads next 8 bytes of data as int32, LE
func (buf *Buffer) ReadInt32() int32 {
	val := binary.LittleEndian.Uint32(buf.b[buf.pos : buf.pos+4])
	buf.pos += 4
	return int32(val)
}

// ReadUint32 reads next 8 bytes of data as uint32, LE
func (buf *Buffer) ReadUint32() uint32 {
	val := binary.LittleEndian.Uint32(buf.b[buf.pos : buf.pos+4])
	buf.pos += 4
	return val
}

// ReadInt64 reads next 16 bytes of data as int64, LE
func (buf *Buffer) ReadInt64() int64 {
	val := binary.LittleEndian.Uint64(buf.b[buf.pos : buf.pos+8])
	buf.pos += 8
	return int64(val)
}

// ReadUint64 reads next 16 bytes of data as uint64, LE
func (buf *Buffer) ReadUint64() uint64 {
	val := binary.LittleEndian.Uint64(buf.b[buf.pos : buf.pos+8])
	buf.pos += 8
	return val
}

// ReadCompactSize reads N byte of data as uint64, LE.
// N depends on the first byte
func (buf *Buffer) ReadCompactSize() uint64 {
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

// ReadVarint reads N byte of data as uint64, LE.
// N depends on the first byte
func (buf *Buffer) ReadVarint() uint64 {
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
