package misc

import (
  "crypto/sha256"
  "golang.org/x/crypto/ripemd160"
	"hash"
)

// Calculate the hash of hasher over buf.
func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

func DoubleSha256(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), sha256.New())
}

// Hash160 calculates the hash ripemd160(sha256(b)).
func Hash160(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), ripemd160.New())
}

/*
* Input: SEC format (compressed or uncompressed)
* Returns a human-readable payment address.
* Used in P2PKH and P2SH (Tested)
*/
func SecToAddress(sec []byte) string {
	prefix := []byte{0x00}
	hash160 := Hash160(sec) // SEC to hash160: ripemd160(sha256(SEC))
  b := append(prefix, hash160...)
	chksum := DoubleSha256(b)[:4] // 
  b = append(b, chksum...)

	return EncodeBase58(b)
}

// p2pkh
func Hash160ToAddress(hash160 []byte, prefix []byte) string {
  b := append(prefix, hash160...)
	chksum := DoubleSha256(b)[:4] // 
  b = append(b, chksum...)

	return EncodeBase58(b)
}

// p2psh p2wpkh
