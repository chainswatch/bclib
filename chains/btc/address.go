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

func readWeirdEcdsa(script []byte) ([]byte, error) {
  var err error = nil

  if script[65] != 0xAC {
    err = fmt.Errorf("WeirdECDSA: Should end with an OP_CHECKSIG (0xAC) opcode. Found 0x%X", script[66])
  }
  // TODO: Convert 65 bytes ECDSA to 25 byte publick key hash form
  return misc.EcdsaToPkeyHash(script[0:65]), err
}

func readEcdsa(script []byte) ([]byte, error) {
  var err error = nil

  if script[0] != 65 {
    err = fmt.Errorf("ECDSA: ECDSA 65 byte public key address should have a length of 65. Found ", script[0])
  }
  if script[66] != 0xAC {
    err = fmt.Errorf("ECDSA: Should end with an OP_CHECKSIG (0xAC) opcode. Found 0x%X", script[66])
  }
  // TODO: Convert 65 bytes ECDSA to 25 byte publick key hash form
  return misc.EcdsaToPkeyHash(script[1:66]), err
}

/*
* ScriptLength == 25
* Script[0] = OP_DUP (0x76)
* Script[1] = OP_HASH160 (0xA9)
* Script[2] = 20 (The length of the public key hash address which follows)
* Script[3-24] = The 20 byte public key address.
*/
func readShortEcdsa(script []byte) ([]byte, error) {
  var err error = nil

  if script[0] != OP_DUP { // OP_DUP
    err = fmt.Errorf("Short ECDSA: Should start with OP_DUP (0x76) opcode. Found 0x%X", script[0])
  }
  if script[1] != OP_HASH160 { // OP_HASH160
    err = fmt.Errorf("Short ECDSA: Expected OP_HASH160 (0xA9) opcode. Found 0x%X", script[1])
  }
  if script[2] != 20 {
    err = fmt.Errorf("Short ECDSA: Expected length 20. Found %d", script[2])
  }
  // TODO: OP_CHECKSIG ??
  return script[3:25], err
}

func readError(script []byte) error {
  var err error = nil

  if script[0] != 0x76 {
    err = fmt.Errorf("Error: The script should start with OP_DUP (0x76) opcode. Found 0x%X", script[0])
  }
  if script[1] != 0xA9 {
    err = fmt.Errorf("Error: Expected OP_HASH160 (0xA9) opcode. Found 0x%X", script[1])
  }
  if script[2] != 0 {
    err = fmt.Errorf("Error: Expected length 0. Found ", script[2])
  }
  if script[4] != 0xAC {
    err = fmt.Errorf("Error: Should end with an OP_CHECKSIG (0xAC) opcode. Found 0x%X", script[4])
  }

  return err
}

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
func isPubkeyHash(ops [][]byte) bool {
  if len(ops) == 5 {
    if ops[0][0] == OP_DUP &&
    ops[1][0] == OP_HASH160 &&
    isOpPubkeyhash(ops[2]) &&
    ops[3][0] == OP_EQUALVERIFY &&
    ops[4][0] == OP_CHECKSIG {
      return true
    }
  }
  return false
}

// P2SH: OP_HASH160, OP_PUBKEYHASH, OP_EQUAL
func isScriptHash(ops [][]byte) bool {
  if len(ops) == 3 {
    if ops[0][0] == OP_HASH160 &&
    isOpPubkeyhash(ops[1]) &&
    ops[2][0] == OP_EQUAL {
      return true
    }
  }
  return false
}

// P2PK: OP_PUBKEY, OP_CHECKSIG
func isPubkey(ops [][]byte) bool {
  if len(ops) == 2 {
    if ops[0][0] == OP_CHECKSIG && isOpPubkey(ops[1]) {
      return true
    }
  }
  return false
}

func isMultiSig(ops [][]byte) bool {
  opLength := len(ops)
  if opLength < 3 || opLength > (16 + 3) {
    return false
  }
  return false
}

func getAddress(script []byte) {
  ops, err := getNumOps(script)
  if err != nil {
    log.Info(err)
  }
  opsLength := len(ops)
  log.Info("Number of ops: ", opsLength)
  for i := 0; i < opsLength; i++ {
    log.Info(fmt.Sprintf("%#x", ops[i]))
  }
  if isPubkeyHash(ops) {
    log.Info("Format: PubkeyHash")
  } else if isScriptHash(ops) {
    log.Info("Format: ScriptHash")
  } else if isPubkey(ops) {
    log.Info("Format: Pubkey")
  } else {
    log.Info("Format: NOT FOUND")
  }
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
