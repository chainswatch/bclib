package misc

import (
  "app/models"
  "crypto/sha256"
  "encoding/hex"
  "fmt"
)

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

func DoubleSha256(data []byte) models.Hash256 {
	hash := sha256.New()
	hash.Write(data)
	firstSha256 := hash.Sum(nil)
	hash.Reset()
	hash.Write(firstSha256)
	return hash.Sum(nil)
}

func HexToByte(hash []byte) []byte {
  blockHashInBytes := make([]byte, hex.DecodedLen(len(hash)))
  n, err := hex.Decode(blockHashInBytes, hash)
  if err != nil {
    fmt.Println(err)
  }
  // Reverse hex to get the LittleEndian order
  return ReverseHex(blockHashInBytes[:n])
}
