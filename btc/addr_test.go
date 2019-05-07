package btc

import (
	"github.com/chainswatch/bclib/serial"

	"bytes"
	"strings"
	"testing"
)

func TestAddr(t *testing.T) {
	// TODO: txP2wpkh, txMultisig
	tests := []struct {
		name					string
		addr					string
		pkey					[]byte
		scriptPkey		[]byte
		txtype				uint8
		isWitness			bool
	}{
		{
			name:     			"DUP HASH160 P20 <pkey> EQVERIFY CHECKSIG",
			addr: 					"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX",
			pkey:     			[]byte("119b098e2e980a229e139a9ed01a469e518e6f26"),
			scriptPkey:     []byte("76a914119b098e2e980a229e139a9ed01a469e518e6f2688ac"),
			txtype:					txP2pkh,
			isWitness:			false,
		},
		{
			name:     			"HASH160 P20 <pkey> EQUAL",
			addr: 					"37xLqiZQYD4WXg2BaVgYfBS2NeTdHwtcuY",
			pkey:     			[]byte("44b6c8fa2b9d3e646c70887651f2861e2cb4ff7f"),
			scriptPkey:    	[]byte("a91444b6c8fa2b9d3e646c70887651f2861e2cb4ff7f87"),
			txtype:					txP2sh,
			isWitness:			false,
		},
		{
			name:     			"P65 <pkey> CHECKSIG",
			addr: 					"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX",
			pkey:     			[]byte("119b098e2e980a229e139a9ed01a469e518e6f26"),
			scriptPkey:    	[]byte("410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac"),
			txtype:					txP2pk,
			isWitness:			false,
		},
		{
			name:     			"0 P32 <pkey>",
			addr: 					"bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej",
			pkey:     			[]byte("701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d"),
			scriptPkey:    	[]byte("0020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d"),
			txtype:					txP2wsh, // V0_P2WSH
			isWitness:			false,
		},
		{
			name:     			"0 P20 <pkey>",
			addr: 					"bc1qj7qmp2h92ls935e4wcnt3zj89uu2sdjwq6hpag",
			pkey:     			[]byte("9781b0aae557e058d3357626b88a472f38a8364e"),
			scriptPkey:    	[]byte("00149781b0aae557e058d3357626b88a472f38a8364e"),
			txtype:					txP2wpkh, // V0_P2WSH
			isWitness:			false,
		},
	}

	for _, test := range tests {
		script, err := serial.HexToBinary(test.scriptPkey)
		if err != nil {
			t.Errorf("%s: %s", test.name, err)
		}
		txtype, pkey := getPkeyFromScript(script)
		if txtype != test.txtype {
			t.Errorf("%s: Wrong txType: %d != %d", test.name, txtype, test.txtype)
		}
		tpkey, _ := serial.HexToBinary(test.pkey)
		if bytes.Compare(pkey, tpkey) != 0 {
			t.Errorf("%s: Wrong pkey from script: %x != %x", test.name, pkey, tpkey)
		}

		// Encoding
		addr, err := EncodeAddr(test.txtype, tpkey)
		if err != nil {
			t.Errorf("%s: %s", test.name, err)
		}
		if strings.Compare(addr, test.addr) != 0 {
			t.Errorf("%s: Wrong encoding: %s != %s", test.name, addr, test.addr)
		}

		// Decoding
		pkey, err = DecodeAddr(test.addr)
		if err != nil {
			t.Errorf("%s: %s", test.name, err)
		}
		if bytes.Compare(pkey, tpkey) != 0 {
			t.Errorf("%s: Wrong decoding: %x != %x", test.name, pkey, tpkey)
		}
	}
}

func TestSegwit(t *testing.T) {
	tests := []struct {
		encoded					string
		decoded					[]byte
	}{
		{
			encoded: 					"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
			decoded:     			[]byte("0014751e76e8199196d454941c45d1b3a323f1433bd6"),
		},
		{
			encoded: 					"tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sl5k7",
			decoded:     			[]byte("00201863143c14c5166804bd19203356da136c985678cd4d27a1b8c6329604903262"),
		},
		{
			encoded: 					"bc1pw508d6qejxtdg4y5r3zarvary0c5xw7kw508d6qejxtdg4y5r3zarvary0c5xw7k7grplx",
			decoded:     			[]byte("5128751e76e8199196d454941c45d1b3a323f1433bd6751e76e8199196d454941c45d1b3a323f1433bd6"),
		},
		{
			encoded: 					"bc1sw50qa3jx3s",
			decoded:     			[]byte("6002751e"),
		},
		{
			encoded: 					"bc1zw508d6qejxtdg4y5r3zarvaryvg6kdaj",
			decoded:     			[]byte("5210751e76e8199196d454941c45d1b3a323"),
		},
		{
			encoded: 					"tb1qqqqqp399et2xygdj5xreqhjjvcmzhxw4aywxecjdzew6hylgvsesrxh6hy",
			decoded:     			[]byte("0020000000c4a5cad46221b2a187905e5266362b99d5e91c6ce24d165dab93e86433"),
		},
	}

	for _, test := range tests {
		tdecoded, _ := serial.HexToBinary(test.decoded)

		// TODO: Upper Case

		// Decoding
		version, output, err := segwitAddrDecode(test.encoded[:2], test.encoded)
		if err != nil {
			t.Error(err)
		}
		decoded := segwitScriptpubkey(version, output)
		if bytes.Compare(decoded, tdecoded) != 0 {
			t.Errorf("%x %x != %x", version, decoded, tdecoded)
		}

		// Encoding
		encoded, err := segwitAddrEncode(test.encoded[:2], version, tdecoded[2:])
		if err != nil {
			t.Error(err)
		}
		if strings.Compare(encoded, test.encoded) != 0 {
			t.Errorf("%s != %s", encoded, test.encoded)
		}

	}
}
