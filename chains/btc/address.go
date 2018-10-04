package btc

/*

Functions used to discriminate between addresses

*/

import (
  "app/misc"
  "fmt"
  "encoding/binary"
  log "github.com/sirupsen/logrus"
)

/*
* Get Ops from Script
*/
func getNumOps(script []byte) ([][]byte, error) {
  var err error = nil
  scriptLength := uint32(len(script))
  log.Debug("Script Length: ", scriptLength)
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

    if (opCode < OP_PUSHDATA1) {
      dataLength = uint32(opCode)
    } else if (opCode == OP_PUSHDATA1) {
      dataLength = uint32(script[i])
    } else if (opCode == OP_PUSHDATA2) {
      dataLength = binary.LittleEndian.Uint32(script[i:(i+2)])
    } else if (opCode == OP_PUSHDATA4) {
      dataLength = binary.LittleEndian.Uint32(script[i:(i+4)])
    } else {
      ops = append(ops, []byte{opCode})
      log.Debug("getNumOps: i=", i," dataLength=", dataLength, fmt.Sprintf(", opCode: %#x", opCode))
      continue
    }
    log.Debug("getNumOps: i=", i," dataLength=", dataLength, fmt.Sprintf(", opCode: %#x", opCode))

    // don't alloc a push buffer if there is no more data available
    if (i + dataLength >= scriptLength) {
      err = fmt.Errorf("Buffer overflow")
      return nil, err
    }

    ops = append(ops, script[i:(i+dataLength)])
    i += dataLength
  }
  ops = append(ops, []byte{script[i]})
  log.Debug(fmt.Sprintf("Last opCode: %#x", script[i]))
  return ops, nil
}

/*
* Check if OP is a PubkeyHash (length == 20)
*/
func isOpPubkeyhash(op []byte) bool {
  // TOPO: OP_PUSHDATA4
  if len(op) != 20 {
    return false
  }
  return true
}

/*
*
*/
func isOpPubkey(op []byte) bool {
  // TOPO: OP_PUSHDATA4
  dataLength := len(op)
  if (dataLength != BTC_ECKEY_COMPRESSED_LENGTH && dataLength != BTC_ECKEY_UNCOMPRESSED_LENGTH) {
    return false
  }
  return true
}

// P2PKH: OP_DUP, OP_HASH160, OP_PUBKEYHASH, OP_EQUALVERIFY, OP_CHECKSIG
func scriptIsPubkeyHash(ops [][]byte) []byte {
  if len(ops) == 5 {
    if ops[0][0] == OP_DUP &&
    ops[1][0] == OP_HASH160 &&
    isOpPubkeyhash(ops[2]) &&
    ops[3][0] == OP_EQUALVERIFY &&
    ops[4][0] == OP_CHECKSIG {
      return ops[2]
    }
  }
  return nil
}

// P2SH: OP_HASH160, OP_PUBKEYHASH, OP_EQUAL
func scriptIsScriptHash(ops [][]byte) []byte {
  if len(ops) == 3 {
    if ops[0][0] == OP_HASH160 &&
    isOpPubkeyhash(ops[1]) &&
    ops[2][0] == OP_EQUAL {
      return ops[1]
    }
  }
  return nil
}

// P2PK: OP_PUBKEY, OP_CHECKSIG
func scriptIsPubkey(ops [][]byte) []byte {
  if len(ops) == 2 {
    if ops[0][0] == OP_CHECKSIG && isOpPubkey(ops[1]) {
      return ops[1]
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

/*
* A witness program is any valid script that consists of a 1-byte push opcode
* followed by a data push between 2 and 40 bytes
*/
func scriptIsWitnessprogram(script []byte, version int32) bool {
  if (version != 0) {
    return false
  }
  lengthScript := len(script)
  if (lengthScript < 4 || lengthScript > 42) {
    return false
  }
  if (script[0] != OP_0 && (script[0] < OP_1 || script[0] > OP_16)) {
    return false
  }
  if (int(script[1] + 2) == lengthScript) {
    log.Debug("WITNESS")
  }
  return false
}

func getAddress(script []byte, version int32) {
  ops, err := getNumOps(script)
  if err != nil {
    log.Info(err)
  }
  opsLength := len(ops)
  log.Info("Number of ops: ", opsLength)
  for i := 0; i < opsLength; i++ {
    log.Info(fmt.Sprintf("%#x", ops[i]))
  }
	var hash []byte
  if hash = scriptIsPubkeyHash(ops); hash != nil {
    log.Info("Script: PubkeyHash, ", misc.SecToAddress(hash))
		// btc_pubkey_get_hash160
  } else if hash = scriptIsScriptHash(ops); hash != nil {
    log.Info("Script: ScriptHash")
  } else if hash = scriptIsPubkey(ops); hash != nil {
    log.Info("Script: Pubkey, ", misc.SecToAddress(hash))
  } else if hash = scriptIsMultiSig(ops); hash != nil {
    log.Info("Script: Multisig, ", len(hash))
  } else {
    log.Info("Script: NOT FOUND")
  }

  scriptIsWitnessprogram(script, version)
}

/*
func getAddress(script []byte) {
  scriptLength := len(script)

  var res []byte
  var err error
  if scriptLength == 67 {
    res, err = readEcdsa(script)
  } else if scriptLength == 66 {
    res, err = readWeirdEcdsa(script)
  } else if scriptLength >= 25 { // Most common format
    res, err = readShortEcdsa(script)
  } else if scriptLength == 5 {
    err = readError(script)
  }
  if err != nil {
    log.Fatal(err)
  }
  log.Debug(fmt.Sprintf("Output Address: %x", res))
  if (len(res) == 100) {
   log.Info(res)
 }
}
*/
