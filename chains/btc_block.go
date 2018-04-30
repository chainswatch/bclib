package chains

import (
  "os"
  "fmt"
  "time"
)

const (
	//! Magic numbers to identify start of block
	BLOCK_MAGIC_ID_BITCOIN MagicID = 0xd9b4bef9
	BLOCK_MAGIC_ID_TESTNET MagicID = 0x0709110b
)

type BlockFile struct {
	file    *os.File
	FileNum uint32
}

func (btc *BtcBlock) NewBlockFile(fileNum uint32) (*BlockFile, error) {
	filepath := fmt.Sprintf(btc.DataDir+"/blocks/blk%05d.dat", fileNum)
	//fmt.Printf("Opening file %s...\n", filepath)

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &BlockFile{file: file, FileNum: fileNum}, nil
}

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func (btc *BtcBlock) ParseBlockHeaderFromFile(blockFile *BlockFile) {
	btc.Length = blockFile.ReadUint32()
	btc.NVersion = blockFile.ReadInt32()
	btc.HashPrevBlock = blockFile.ReadBytes(32)
	btc.HashMerkleRoot = blockFile.ReadBytes(32)
	btc.NTime = time.Unix(int64(blockFile.ReadUint32()), 0)
	btc.TargetDifficulty = blockFile.ReadUint32() // TODO: Parse this as mantissa?
	btc.NNonce = blockFile.ReadUint32()
}

func (btc *BtcBlock) ParseBlockFromFile(blockFile *BlockFile) {
	curPos, err := blockFile.Seek(0, 1)
	if err != nil {
		fmt.Println(err)
	}

	// Read and validate Magic ID
	btc.MagicID = MagicID(blockFile.ReadUint32())
	if btc.MagicID != BLOCK_MAGIC_ID_BITCOIN {
		blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    fmt.Println("Invalid block header: Can't find Magic ID")
	}

	// Read header fields
	btc.ParseBlockHeaderFromFile(blockFile)

	// Parse transactions
	err = btc.ParseBlockTransactionsFromFile(blockFile)
	if err != nil {
		blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    fmt.Println(err)
	}
}

func (btc *BtcBlock) NewBlockFromFile(num uint32, pos uint32) {
	// Open file for reading
	blockFile, err := btc.NewBlockFile(num)

	if err != nil {
    fmt.Println(err)
	}
	defer blockFile.Close()

	// Seek to pos - 8 to start reading from block header
	fmt.Printf("Seeking to block at %d...\n", pos)
	_, err = blockFile.Seek(int64(pos-8), 0)
	if err != nil {
    fmt.Println(err)
	}

	btc.ParseBlockFromFile(blockFile)
}

func (btc *BtcBlock) getBlockMeta(nHeight uint32) {
  // btc.pgDb.Where("n_height = ?", nHeight)
}

func (btc *BtcBlock) getBlock(nHeight uint32) {
  // Get block infos from... (in particular NFile, NDataPos and NUndoPos)
  btc.getBlockMeta(nHeight)
}
