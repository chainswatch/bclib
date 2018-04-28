// +build ignore
package chains

import (
	"fmt"
	"log"
	"bytes"
	//"time"
	"encoding/json"
	//"github.com/golang/protobuf/proto"
	"app/db"
	//"app/generated/qrl"
)

func binToHstr(vec []byte) string {
	var buffer bytes.Buffer
	for _,c := range(vec) {
		buffer.WriteString(fmt.Sprintf("%02x", c))
	}
	return buffer.String()
}

func getHeaderhashByHeight(stateDb *db.StateDb, blockHeight []byte) ([]byte, error) {
	data, err := stateDb.Get(blockHeight, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	blockNumberMapping := &qrl.BlockNumberMapping {}
	if err := json.Unmarshal(data, blockNumberMapping); err != nil {
		return nil, err
	}
	// fmt.Println(bin2hstr(block_number_mapping.Headerhash))
	return block_number_mapping.Headerhash, nil
}

func getBlockByHeaderHash(stateDb *db.StateDb, headerHash []byte) (*qrl.Block, error) {
	data, err := stateDb.Get([]byte(binToHstr(headerHash)), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
	block := &qrl.Block {}
	if err := json.Unmarshal(data, block); err != nil {
		log.Fatal(err)
	}
	return block, nil
}

func qrlWatcher(dataDir string) {
	stateDb, _ := db.OpenStateDb(dataDir)
	defer stateDb.Close()

	//json_data = self._db.get_raw(bin2hstr(header_hash).encode())
	//data, err := indexDb.Get(append([]byte("t"), txHash...), nil)
	headerHash,_ := getHeaderhashByHeight(stateDb, []byte("1000"))
	getBlockByHeaderHash(stateDb, headerHash)
	//fmt.Println(block)
	/*
	s := bin2hstr(block_number_mapping.Headerhash)
	fmt.Println(s)
	data, err = StateDb.Get([]byte(s), nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}


	iter := StateDb.NewIterator(nil, nil)
	block_count := 0
	total := 0
	for iter.Next() {
		total++
		//value := iter.Key()
		value := iter.Value()
		block := &qrl.Block {}
		if err := proto.Unmarshal(value, block); err == nil {
			b_header := &qrl.BlockHeader {}
			if err := proto.Unmarshal(block.header, b_header); err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println(block.Transactions)
			if (len(block.Transactions) > 0) {
				fmt.Println("TRANSACTION")
			}
			block_count++
		}
		//timer1 := time.NewTimer(10 * time.Millisecond)
		//<-timer1.C
	}
	fmt.Println("DONE: ", block_count, "/", total)
	iter.Release()
	*/
}
