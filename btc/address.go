package btc

import (
  "git.posc.in/cw/bclib/serial"
  "fmt"
  "encoding/binary"
  log "github.com/sirupsen/logrus"
)

// Get Ops from Script
func getNumOps(script []byte) ([][]byte, error) {
  var err error
  scriptLength := uint32(len(script))
  if scriptLength == 0 {
    err = fmt.Errorf("Script of length 0")
    return nil, err
  }
  ops := [][]byte{}
  var i, dataLength uint32
  for i = 0; i < scriptLength - 1; {
    dataLength = 0
    opCode := script[i]
    i++

    if opCode < opPushdata1 && opCode > op0 {
      dataLength = uint32(opCode)
    } else if opCode == opPushdata1 {
      dataLength = uint32(script[i])
    } else if opCode == opPushdata2 {
      dataLength = binary.LittleEndian.Uint32(script[i:(i+2)])
    } else if opCode == opPushdata4 {
      dataLength = binary.LittleEndian.Uint32(script[i:(i+4)])
    } else {
      ops = append(ops, []byte{opCode})
      continue
    }

    // don't alloc a push buffer if there is no more data available
    if (i + dataLength > scriptLength) {
      err = fmt.Errorf("Buffer overflow: %d + %d > %d", i, dataLength, scriptLength)
      return nil, err
    }
    ops = append(ops, script[i:(i+dataLength)])
    i += dataLength
  }
  if i < scriptLength {
    ops = append(ops, []byte{script[i]})
  }
  return ops, nil
}

// Check if OP is a PubkeyHash (length == 20)
func isOpPubkeyhash(op []byte) bool {
  // TODO: OP_PUSHDATA4
  if len(op) != 20 {
    return false
  }
  return true
}

/*
*
*/
func isOpPubkey(op []byte) bool {
  // TODO: OP_PUSHDATA4
  dataLength := len(op)
  if (dataLength != btcEckeyCompressedLength && dataLength != btcEckeyUncompressedLength) {
    return false
  }
  return true
}

// P2PKH: OP_DUP, OP_HASH160, OP_PUBKEYHASH, OP_EQUALVERIFY, OP_CHECKSIG
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

// P2SH: OP_HASH160, OP_PUBKEYHASH, OP_EQUAL
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

// P2PK: OP_PUBKEY, OP_CHECKSIG
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
  if opLength < 3 || opLength > (16 + 3) {
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
func scriptIsWitnessProgram(script []byte, version int32) bool {
  if (version != 0) {
    return false
  }
  lengthScript := len(script)
  if (lengthScript < 4 || lengthScript > 42) {
    return false
  }
  if (script[0] != op0 && (script[0] < op1 || script[0] > op16)) {
    return false
  }
  if (int(script[1] + 2) == lengthScript) {
    return true
  }
  return false
}

func decodeAddress() () {
}

func getPublicAddress(txType uint8, hash []byte) string {
  var address string
  if txType == txP2pkh {
    address = serial.Hash160ToAddress(hash, []byte{0x00})
  } else if txType == txP2sh {
    address = serial.Hash160ToAddress(hash, []byte{0x05})
  } else if txType == txP2pk {
    address = serial.SecToAddress(hash)
  } else if txType == txMultisig {
    log.Info("Script: Multisig, ", len(hash))
    return ""
  } else if txType == txP2wpkh {
    address, _ = serial.EncodeBench32("bc", hash)
  } else if txType == txP2wsh {
    address, _ = serial.EncodeBench32("bc", hash)
  } else if txType == txOpreturn {
    address = fmt.Sprintf("%x", hash)
  } else {
    log.Info("Script: NOT FOUND")
    return ""
  }
  return address
}

func getVersion(op int32) int32 {
  if (op == op0) {
    return 0;
  }
  if (op >= op1 && op <= op16) {
    log.Fatal("Error in getVersion ", op)
  }
  return op - (op1 - 1)
}


/*
* script:
* version:
* Return hash and hash type (P2PKH,P2SH...) from output script
*/
func getAddressFromScript(script []byte) (uint8, []byte) {
  ops, err := getNumOps(script)
  if err != nil {
    log.Info(err)
  }
  opsLength := len(ops)
  version := getVersion(int32(ops[0][0]))
	var outputScript string
  for i := 0; i < opsLength; i++ {
    outputScript += fmt.Sprintf("%#x ", ops[i])
  }
  var hash []byte
  var txType uint8
  if hash = scriptIsPubkeyHash(ops); hash != nil {
    txType = txP2pkh
  } else if hash = scriptIsScriptHash(ops); hash != nil {
    txType = txP2sh
  } else if hash = scriptIsPubkey(ops); hash != nil {
    txType = txP2pk
  } else if hash = scriptIsMultiSig(ops); hash != nil {
    txType = txMultisig
    return 0, nil
  } else if scriptIsWitnessProgram(script, version) {
    hash = append(ops[0], ops[1]...)
    if len(hash) == 20 + 1 {
      txType = txP2wpkh
    } else if len(hash) == 32 + 1 {
      txType = txP2wsh
    }
  } else if hash = scriptIsOpReturn(ops); hash != nil {
    txType = txOpreturn
  } else {
    hash = nil
    txType = txUnknown
  }

  return txType, hash
}
