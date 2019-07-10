package btc

import (
	"math"
)

// getHashRate computes the network hash rate given
// the number of block founds, the number of expected
// blocks and the difficulty
/*
func getHashRate(found, expected , difficulty uint32) {
  return (found / expected) * difficulty * 2 ** 32 / 600
}
*/

// http://charts.woobull.com/bitcoin-nvt-ratio/
// Market Cap / daily USD volume transmitted
// getNVT

// GetBlockReward returns block reward value from height
func GetBlockReward(height uint32) uint64 {
	return uint64(50E8 * math.Pow(0.5, float64(height/210000)))
}

// GetSupply returns total supply from height
func GetTotalSupply(height uint32) (s uint64) {
	for h := uint32(0); h <= height; h++ {
		s += GetBlockReward(h)
	}
	return s
}

// GetHeight returns height from circulating supply
func GetHeight(target uint64) (h uint32) {
	var s uint64

	for h := uint32(0); s <= target; h++ {
		s += uint64(50E8 * math.Pow(0.5, float64(h/210000)))
		if s == target {
			return h
		}
	}
	return 0
}
