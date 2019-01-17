package btc

import (
	"encoding/binary"
	"fmt"
)

func getNextOp(script []byte) ([]byte, []byte) {
	dataLength := uint32(0)
	switch {
	case script[0] < opPushdata1 && script[0] > op0:
		dataLength = uint32(script[0])
	case script[0] == opPushdata1 && len(script) > 1:
		dataLength = uint32(script[1])
	case script[0] == opPushdata2 && len(script) > 2:
		dataLength = binary.LittleEndian.Uint32(append([]byte{0, 0}, script[1:3]...))
	case script[0] == opPushdata4 && len(script) > 4:
		dataLength = binary.LittleEndian.Uint32(script[1:5])
	default:
		return script[:1], script[1:]
	}
	if dataLength >= uint32(len(script)) {
		return script[1:], nil
	}
	return script[1 : 1+dataLength], script[1+dataLength:]
}

// Get Ops from Script
func getOps(raw []byte) (ops [][]byte, err error) {
	script := make([]byte, len(raw))
	copy(script, raw)
	var op []byte
	for len(script) > 0 {
		if op, script = getNextOp(script); script == nil {
			return ops, fmt.Errorf("Overflow")
		}
		ops = append(ops, op)
	}
	return ops, nil
}
