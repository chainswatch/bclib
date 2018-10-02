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
* OP_DUP, OP_HASH160, OP_PUBKEYHASH, OP_EQUALVERIFY, OP_CHECKSIG
* Script[0] = OP_DUP (0x76)
* Script[1] = OP_HASH160 (0xA9)
* Script[2] = 20 (A length value equal to the size of a public key hash)
* Script[3-22]=The 20 byte public key address
* Script[23]=OP_EQUALVERIFY (0x88)
* Script[24]=OP_CHECKSIG (0xAC)
*/
func isPubkeyHash(script []byte) ([]byte, error) {
  var err error = nil

  if script[0] != OP_DUP {
    err = fmt.Errorf("Other: Should start with OP_DUP (0x76) opcode. Found 0x%X", script[0])
  }
  if script[1] != OP_HASH160 {
    err = fmt.Errorf("Other: Expected OP_HASH160 (0xA9) opcode. Found 0x%X", script[1])
  }
  if script[2] != 20 {
    err = fmt.Errorf("Other: Expected length 20. Found ", script[2])
  }
  if script[23] != OP_EQUALVERIFY {
    err = fmt.Errorf("Other: Expected OP_EQUALVERIFY (0x88) opcode. Found 0x%X", script[23])
  }
  if script[24] != OP_CHECKSIG {
    err = fmt.Errorf("Other: Should end with an OP_CHECKSIG (0xAC) opcode. Found 0x%X", script[24])
  }
  return script[3:23], err
}

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

    log.Debug("getNumOps: i=", i)
    if (opCode < OP_PUSHDATA1) {
      dataLength = uint32(opCode)
    } else if (opCode == OP_PUSHDATA1) {
      dataLength = uint32(script[i])
    } else if (opCode == OP_PUSHDATA2) {
      dataLength = binary.LittleEndian.Uint32(script[i:(i+2)])
    } else if (opCode == OP_PUSHDATA4) {
      dataLength = binary.LittleEndian.Uint32(script[i:(i+4)])
    } else {
      ops = append(ops, []byte{script[i]})
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
  log.Debug(fmt.Sprintf("Last opCode: %#x", script[i]))
  return ops, nil
}

func getAddress(script []byte) {
  _, err := getNumOps(script)
  if err != nil {
    log.Debug(err)
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
