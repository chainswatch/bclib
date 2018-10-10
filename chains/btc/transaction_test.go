package btc

import (
	"app/misc"
	"testing"
	"app/chains/parser"
	"fmt"
)

func TestTransaction(t *testing.T) {
	tests := []struct {
		name    string
		txhash  string
		rawtx		[]byte
		encoded []string
	}{
		{
			name:		"Find a name",
			txhash: "65a3235950069fe3f8cb5428ac72960ec0ea8ed09f77847a6622fff249c9d967",
			rawtx:	[]byte("01000000017d171d1fb676322e54684890fa2963dee483e024db6049319565841974fcb460010000006a47304402204c37cd1f21ac7016fd11e8b743a867a8f663d59f1796044f90950bea00429d5c022055e828cc7de3e0106139e30fc5a233c27fca2ff349f3e5f06cd56568dcc9db580121033dab59fca27b7f6ee99fb46fc8f8afe3e0841094ab1c55163644c2051a5a0885ffffffff02ecee0e00000000001976a914a028c390ffb66fe0164bb0a9154c4e76aafb988b88ac65dd17000000000017a9145d4674e8466725f60e8379a259948931ca73f5be8700000000"),
			encoded:	[]string{"1Fbqx2BqxNPcD6Gbnpd6XAGXpWSZ3P8xxq", "3ACDAoJhzHxLEcDgf5vxH2VeJftNY6NANY"},
		},
		{
			name:		"Find a name",
			txhash:	"576534d167cf564fe395c16ebef382a1e9e553a10ffdae6eddee360606ee37f6",
			rawtx:	[]byte("02000000000101f817ede6389ce1633b04afe3ac2011e3b91a8646c2e3adfc424cb4819bc7b5be01000000171600146cb6b34581cd008dd0d3248b6e9ff2b4285c1019feffffff02fa2ec2000000000017a91444b6c8fa2b9d3e646c70887651f2861e2cb4ff7f87f18705000000000017a914f6902c9bd92f0c089a9f94341d7aba56634534668702483045022100b6c42756abf43ea7e7d27c1b5021e39d928d78f452758b6b7cfae49438d7a22202206f2b00b932133db3be1236498bc6cea4c42dd7fc585680631a4dc17d88771dba012102e91f9ee10015228fa9d79d7599720e9a637553b219d42742677d7500e2b4b603d84d0800"),
			encoded:	[]string{"37xLqiZQYD4WXg2BaVgYfBS2NeTdHwtcuY", "3QAis8tJTf8Ahr29jYM5Znrb4oE9u9TUNs"},
		},
		{
			name:		"Bench32",
			txhash:	"1e52891604b9ec5ccec8d94a404b3d08ef2cc6195137c7437d28c1fc52e3b663",
			rawtx:	[]byte("0100000000010181ec646c82f7e9a5a8d9adfb46162415172e891c37c4412ee60f8aaefd599dda0100000000ffffffff04200b2000000000001976a9147db7cf2faff951ee575ce8cbc2afadfc6f83986488acc0c62d00000000001976a9147dde2cdaef4c0be4149618fea926c50ee0e5016388ac70991400000000001976a9146b5dcb791a7ecc715aa4f30a6d24abb06351228d88ac694db50000000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d0400483045022100a16308296bbb33b92f82ef1d4b75b3b18a4aea8d74899476b80016d1ec2e107d0220756c99a26759c2e4417014030176fb59d022566562744fbb2eb33d70892a2bb201483045022100f39bda511ee4e07918bc3d16c903e3d0837672182682b2b18491a2b2d4f099a802205e85688ac3d8332353f29c8b79dde7e7e3a28494b648f1edf3c2dc6e5858827d016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000"),
			encoded:	[]string{"1CTjfGuuquBv2BFCWXeMGzKanWHXJ2FQr6",
			"1CUXcz3DQML6WmgNmF1iA4Bc7Mh99k2acD",
			"1AnhfWELiznEuFnKLQ2E1b3z99cjstkBfg",
			"bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej"},
		},
	}

	for _, test := range tests {
		rawtx := misc.HexToBinary(test.rawtx)
		rawTx := &parser.RawTx{Body: rawtx, Pos: 0}
		tx, _ := parseTransaction(rawTx)
		txHash := fmt.Sprintf("%x", misc.ReverseHex(tx.Hash))
		if txHash != test.txhash {
			t.Errorf("%v: String on decoded value does not match expected value: %v != %v",
			test.name, test.txhash, txHash)
		}

		if (int(tx.NVout) != len(test.encoded)) {
			t.Error("Wrong number of output. Should be tx.NVout = 2")
		}
		for idx, vout := range tx.Vout {
			txType, hash := getAddressFromScript(vout.Script)
			decoded := getPublicAddress(txType, hash)
			if test.encoded[idx] != decoded {
				t.Errorf("%v: String on decoded value does not match expected value: %v != %v",
				test.name, test.encoded[idx], decoded)
			}
		}
	}
}
