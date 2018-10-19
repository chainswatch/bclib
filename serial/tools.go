package serial

import (
  "encoding/hex"
  log "github.com/sirupsen/logrus"
)

// ReverseHex reverses the order of bytes input
// TODO: Maybe can optimize
func ReverseHex(b []byte) []byte {
  newb := make([]byte, len(b))
  copy(newb, b)
  for i := len(newb)/2 - 1; i >= 0; i-- {
    opp := len(newb) - 1 - i
    newb[i], newb[opp] = newb[opp], newb[i]
  }

  return newb
}

// HexToBinary transforms human readable address into binary address
func HexToBinary(src []byte) []byte {
  b := make([]byte, hex.DecodedLen(len(src)))
  n, err := hex.Decode(b, src)
  if err != nil {
    log.Warn(err)
  }
	return b[:n]
}
