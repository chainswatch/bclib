package btc

import (
  "app/misc"
  log "github.com/sirupsen/logrus"
  "fmt"
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

func readShortEcdsa(script []byte) ([]byte, error) {
  var err error = nil

  if script[0] != 0x76 {
    err = fmt.Errorf("Short ECDSA: Should start with OP_DUP (0x76) opcode. Found 0x%X", script[0])
  }
  if script[1] != 0xA9 {
    err = fmt.Errorf("Short ECDSA: Expected OP_HASH160 (0xA9) opcode. Found 0x%X", script[1])
  }
  if script[2] != 20 {
    err = fmt.Errorf("Short ECDSA: Expected length 20. Found %s", script[2])
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

func lastChance(script []byte) ([]byte, error) {
  var err error = nil

  if script[0] != 0x76 {
    err = fmt.Errorf("Other: Should start with OP_DUP (0x76) opcode. Found 0x%X", script[0])
  }
  if script[1] != 0xA9 {
    err = fmt.Errorf("Other: Expected OP_HASH160 (0xA9) opcode. Found 0x%X", script[1])
  }
  if script[2] != 20 {
    err = fmt.Errorf("Other: Expected length 20. Found ", script[2])
  }
  if script[23] != 0x88 {
    err = fmt.Errorf("Other: Expected OP_EQUALVERIFY (0x88) opcode. Found 0x%X", script[23])
  }
  if script[24] != 0x88 {
    err = fmt.Errorf("Other: Should end with an OP_CHECKSIG (0xAC) opcode. Found 0x%X", script[24])
  }

  return script[3:23], err
}

func getAddress(script []byte) {
  scriptLength := len(script)
  // log.Info("Script Length:", scriptLength)

  var res []byte
  var err error
  if scriptLength == 67 {
    res, err = readEcdsa(script)
  } else if scriptLength == 66 {
    res, err = readWeirdEcdsa(script)
  } else if scriptLength >= 25 {
    res, err = readShortEcdsa(script)
  } else if scriptLength == 5 {
    err = readError(script)
  }
  if err != nil {
    log.Fatal(err)
  }
  /*
  if res == nil {
    res = lastChance(script)
  }
  */
  if (len(res) == 100) {
   log.Info(res)
 }
}
