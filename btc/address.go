package btc

// TODO: Rename to script.go

import (
	"fmt"
	"github.com/chainswatch/bclib/serial"
	log "github.com/sirupsen/logrus"
)

// Check if OP is a PubkeyHash (length == 20)
func isOpPubkeyhash(op []byte) bool {
	// TODO: OP_PUSHDATA4
	if len(op) != 20 {
		return false
	}
	return true
}

func isOpPubkey(op []byte) bool {
	// TODO: OP_PUSHDATA4
	dataLength := len(op)
	if dataLength != btcEckeyCompressedLength && dataLength != btcEckeyUncompressedLength {
		return false
	}
	return true
}

// P2PKH
func scriptIsPubkeyHash(ops [][]byte) []byte {
	if len(ops) == 5 {
		if ops[0][0] == opDup &&
			ops[1][0] == opHash160 &&
			isOpPubkeyhash(ops[2]) &&
			ops[3][0] == opEqualverify &&
			ops[4][0] == opChecksig {
			return ops[2]
		}
	}
	return nil
}

// P2SH
func scriptIsScriptHash(ops [][]byte) []byte {
	if len(ops) == 3 {
		if ops[0][0] == opHash160 &&
			isOpPubkeyhash(ops[1]) &&
			ops[2][0] == opEqual {
			return ops[1]
		}
	}
	return nil
}

// P2PK
func scriptIsPubkey(ops [][]byte) []byte {
	if len(ops) == 2 {
		if ops[1][0] == opChecksig && isOpPubkey(ops[0]) {
			return ops[0]
		}
	}
	return nil
}

func scriptIsMultiSig(ops [][]byte) []byte {
	opLength := len(ops)
	if opLength < 3 || opLength > (16+3) {
		return nil
	}
	return nil
}

func scriptIsOpReturn(ops [][]byte) []byte {
	if len(ops) == 2 && ops[0][0] == opReturn && len(ops[1]) <= 20 {
		return ops[1]
	}
	return nil
}

// A witness program is any valid script that consists of a 1-byte push opcode
// followed by a data push between 2 and 40 bytes
func scriptIsWitnessProgram(script []byte) bool {
	lengthScript := len(script)
	if lengthScript < 4 || lengthScript > 42 {
		return false
	}
	if script[0] != op0 && (script[0] < op1 || script[0] > op16) {
		return false
	}
	if int(script[1]+2) == lengthScript {
		return true
	}
	return false
}

// DecodeAddr decodes address from hash
func DecodeAddr(txType uint8, hash []byte) (string, error) {
	var address string
	switch txType {
	case txP2pkh:
		address = serial.Hash160ToAddress(hash, []byte{0x00})
	case txP2sh:
		address = serial.Hash160ToAddress(hash, []byte{0x05})
	case txP2pk:
		address = serial.SecToAddress(hash)
	case txMultisig:
		log.Info("Script: Multisig, ", len(hash))
		return "", nil
	case txP2wpkh:
		address, _ = serial.EncodeBench32("bc", hash)
	case txP2wsh:
		address, _ = serial.EncodeBench32("bc", hash)
	case txOpreturn:
		address = fmt.Sprintf("%x", hash)
	case txParseErr:
		address = ""
	case txUnknown:
		address = ""
	default:
		return "", fmt.Errorf("DecodeAddr: Unable to decode addr from script")
	}
	return address, nil
}

// FIXME: WTF is it useful for?
func getVersion(op int32) int32 {
	if op == op0 {
		return 0
	}
	if op >= op1 && op <= op16 {
		log.Info("Error in getVersion ", op)
	}
	return op - (op1 - 1)
}

/*
* script:
* version:
* Return hash and hash type (P2PKH,P2SH...) from output script
 */
func getPkeyFromScript(script []byte) (txType uint8, hash []byte) {
	ops, err := getOps(script)
	if err != nil {
		return txParseErr, nil
	}
	/*
		opsLength := len(ops)
		var outputScript string
		for i := 0; i < opsLength; i++ {
			outputScript += fmt.Sprintf("%#x ", ops[i])
		}
	*/

	if hash = scriptIsPubkeyHash(ops); hash != nil {
		txType = txP2pkh
	} else if hash = scriptIsScriptHash(ops); hash != nil {
		txType = txP2sh
	} else if hash = scriptIsPubkey(ops); hash != nil {
		txType = txP2pk
	} else if hash = scriptIsMultiSig(ops); hash != nil {
		txType = txMultisig // TODO: MULTISIG
	} else if scriptIsWitnessProgram(script) {
		hash = append(ops[0], ops[1]...)
		if len(hash) == 20+1 {
			txType = txP2wpkh
		} else if len(hash) == 32+1 {
			txType = txP2wsh
		} else {
			txType = txUnknown
		}
	} else if hash = scriptIsOpReturn(ops); hash != nil {
		txType = txOpreturn
	} else {
		hash = nil
		txType = txUnknown
	}

	return txType, hash
}
