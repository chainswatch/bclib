package btc

// TODO: Rename to script.go

import (
	"fmt"
	"github.com/chainswatch/bclib/serial"
	"github.com/chainswatch/bclib/serial/bech32"
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
			return serial.Hash160(ops[0])
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

// TODO: Improve
// A witness program is any valid script that consists of a 1-byte push opcode
// followed by a data push between 2 and 40 bytes
func scriptIsWitnessProgram(ops [][]byte) bool {
	if len(ops) != 2 {
		return false
	}
	if ops[0][0] != op0 && (ops[0][0] < op1 || ops[0][0] > op16) {
		return false
	}
	return true
}

func segwitScriptpubkey(version byte, program []byte) []byte {
	if version != 0 {
		version += 0x50
	}
	return append(append([]byte{version}, byte(len(program))), program...)
}

// segwitAddrDecode decodes hrp(human-readable part) Segwit Address(string), returns version(int) and data(bytes array) / or error
func segwitAddrDecode(hrp, addr string) (byte, []byte, error) {
	dechrp, data, err := bech32.Decode(addr)
	if err != nil {
		return 0, nil, err
	}
	if dechrp != hrp {
		return 0, nil, fmt.Errorf("invalid human-readable part : %s != %s", hrp, dechrp)
	}
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("invalid decode data length : %d", len(data))
	}
	if data[0] > 16 {
		return 0, nil, fmt.Errorf("invalid witness version : %d", data[0])
	}
	pkey, err := bech32.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return 0, nil, err
	}
	if len(pkey) < 2 || len(pkey) > 40 {
		return 0, nil, fmt.Errorf("invalid convertbits length : %d", len(pkey))
	}
	if data[0] == 0 && len(pkey) != 20 && len(pkey) != 32 {
		return 0, nil, fmt.Errorf("invalid program length for witness version 0 (per BIP141) : %d", len(pkey))
	}
	return data[0], pkey, nil
}

// segwitAddrEncode encodes hrp(human-readable part), version and data(bytes array), returns Segwit Address / or error
func segwitAddrEncode(hrp string, version byte, pkey []byte) (string, error) {
	if version > 16 {
		return "", fmt.Errorf("invalid witness version : %d", version)
	}
	if len(pkey) < 2 || len(pkey) > 40 {
		return "", fmt.Errorf("invalid pkey length : %d", len(pkey))
	}
	if version == 0 && len(pkey) != 20 && len(pkey) != 32 {
		return "", fmt.Errorf("invalid pkey length for witness version 0 (per BIP141) : %d", len(pkey))
	}
	data, err := bech32.ConvertBits(pkey, 8, 5, true)
	if err != nil {
		return "", err
	}
	addr, err := bech32.Encode(hrp, append([]byte{version}, data...))
	if err != nil {
		return "", err
	}
	return addr, nil
}

// AddrEncode encodes address from pkey
func AddrEncode(txType uint8, pkey []byte) (string, error) {
	var addr string
	switch txType {
	case txP2pkh:
		addr = serial.Hash160ToAddress(pkey, []byte{0x00})
	case txP2sh:
		addr = serial.Hash160ToAddress(pkey, []byte{0x05})
	case txP2pk:
		addr = serial.Hash160ToAddress(pkey, []byte{0x00})
	case txMultisig:
		return "", fmt.Errorf("Script: Multisig, %d", len(pkey))
	case txP2wpkh:
		addr, _ = segwitAddrEncode("bc", 0x00, pkey)
	case txP2wsh:
		addr, _ = segwitAddrEncode("bc", 0x00, pkey)
	case txOpreturn:
		addr = fmt.Sprintf("%x", pkey)
	case txParseErr:
		addr = ""
	case txUnknown:
		addr = ""
	default:
		return "", fmt.Errorf("EncodeAddr: Unable to encode addr from pkeyscript")
	}
	return addr, nil
}

// AddrDecode accepts an encoded address (P2PKH or P2SH, human readable)
// returns its public key
func AddrDecode(addr string) ([]byte, error) {
	switch {
	case addr[0] == '1':
		data := serial.DecodeBase58(addr)
		if data[0] != 0x00 {
			return nil, fmt.Errorf("Address must start with byte 0x00")
		}
		return data[1:21], nil
	case addr[0] == '3':
		data := serial.DecodeBase58(addr)
		if data[0] != 0x05 {
			return nil, fmt.Errorf("Address must start with byte 0x05")
		}
		return data[1:21], nil
	case addr[0] == 'b' && addr[1] == 'c':
		log.Info(fmt.Sprintf("%s", addr))
		_, data, err := segwitAddrDecode("bc", addr)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, fmt.Errorf("EncodeAddr: Unable to decode pkey from addr")
	// TODO: Check DoubleSha256(data[:21])[:4] == data[-4:]
}

func getPkeyFromOps(ops [][]byte) (txType uint8, pkey []byte) {
	if pkey = scriptIsPubkeyHash(ops); pkey != nil {
		txType = txP2pkh
	} else if pkey = scriptIsScriptHash(ops); pkey != nil {
		txType = txP2sh
	} else if pkey = scriptIsPubkey(ops); pkey != nil {
		txType = txP2pk
	} else if pkey = scriptIsMultiSig(ops); pkey != nil {
		txType = txMultisig // TODO: MULTISIG
	} else if scriptIsWitnessProgram(ops) {
		pkey = ops[1]
		if len(pkey) == 20 { // TODO: Improve
			txType = txP2wpkh
		} else if len(pkey) == 32 {
			txType = txP2wsh
		} else {
			pkey = nil
			txType = txUnknown
		}
	} else if pkey = scriptIsOpReturn(ops); pkey != nil {
		txType = txOpreturn
	} else {
		pkey = nil
		txType = txUnknown
	}
	return txType, pkey
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
	return getPkeyFromOps(ops)
}
