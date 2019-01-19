package serial

/*
* Get tests from:
* https://github.com/btcsuite/btcutil/blob/master/address_test.go
 */

import (
	"testing"
)

func TestSerialization(t *testing.T) {
	tests := []struct {
		name     string
		hashtype string
		addr     []byte
		encoded  string
	}{
		{
			name:     "mainnet p2pk SEC compressed (0x01)",
			hashtype: "SEC",
			addr:     []byte("02192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4"),
			encoded:  "13CG6SJ3yHUXo4Cr2RY4THLLJrNFuG3gUg",
		},
		{
			name:     "mainnet p2pk SEC compressed (0x02)",
			hashtype: "SEC",
			addr:     []byte("03b0bd634234abbb1ba1e986e884185c61cf43e001f9137f23c2c409273eb16e65"),
			encoded:  "15sHANNUBSh6nDp8XkDPmQcW6n3EFwmvE6",
		},
		{
			name:     "mainnet p2pk SEC uncompressed",
			hashtype: "SEC",
			addr:     []byte("0411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2" + "e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3"),
			encoded:  "12cbQLTFMXRnSzktFkuoG3eHoMeFtpTu3S",
		},
		{
			name:     "mainnet p2pk SEC hybrid (0x06)",
			hashtype: "SEC",
			addr:     []byte("06192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4" + "0d45264838c0bd96852662ce6a847b197376830160c6d2eb5e6a4c44d33f453e"),
			encoded:  "1Ja5rs7XBZnK88EuLVcFqYGMEbBitzchmX",
		},
		{
			name:     "mainnet p2pk SEC hybrid (0x07)",
			hashtype: "SEC",
			addr:     []byte("07b0bd634234abbb1ba1e986e884185c61cf43e001f9137f23c2c409273eb16e65" + "37a576782eba668a7ef8bd3b3cfb1edb7117ab65129b8a2e681f3c1e0908ef7b"),
			encoded:  "1ExqMmf6yMxcBMzHjbj41wbqYuqoX6uBLG",
		},
	}

	var decoded string
	for _, test := range tests {
		// if strings.Compare(test.hashtype, "SEC") == 0 {
		if test.hashtype == "SEC" {
			data, err := HexToBinary(test.addr)
			if err != nil {
				t.Fatal(err)
			}
			decoded = SecToAddress(data)
		} else {
			t.Error("Hash type not found")
		}

		if decoded != test.encoded {
			t.Errorf("%v: String on decoded value does not match expected value: %v != %v",
				test.name, test.addr, decoded)
		}

		// TODO: Encode again and compare against the original
	}
}
