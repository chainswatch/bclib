package btc

import (
  "git.posc.in/cw/bclib/serial"
  "git.posc.in/cw/bclib/parser"
	"testing"
	"fmt"
  log "github.com/sirupsen/logrus"
)

func TestTransaction(t *testing.T) {
  log.SetLevel(log.DebugLevel)
	tests := []struct {
		name    	string
		txhash  	string
		rawtx			[]byte
		encoded 	[]string
		addrtype	[]uint8
    nvin    	int
	}{
		{
			name:		  "mainnet simple",
			txhash:	  "07df89107e2f4a03c3e055d745169ef1f8a35053202e26d40513a71fcff81372",
			rawtx:	  []byte("0200000001a938b4d9736f16e9da93e8bce68b7b9565356cbd82113677ac64757fe15181a0010000006b483045022100fbdde0b32524c27a4f89abc0c0c625d0eb9e9d0626ea88160db18aaae58d2f0e022016f4aa8bfa88348d4ce730611b4437a14857df674979efa9c7ae3ad2af424f81012102eb8ce6402d63c5ccebdc401a239ccddb61ce4571272b3c58f98dc464a1162e7efeffffff04d0599900000000001976a914d1a72b95a14ac1d07ea044c94433f3c78c737ca888acf0ab3f01000000001976a9145c76196c569fe8990555602c6039253a373adc9588acc9acac33000000001976a91488903133b8b562a38bc00210e93895b2d43cf6fb88ac98347202000000001976a91423981be1300d3cf2cb0d5c3efaa0a07ce43b403d88ac9b510800"),
			encoded:	[]string{"1L7YVHYiNiRcoRJ14q4oFotU4sXULhz1Kq",
			"19RteLGkZ9PPjDovVz4MJyySaeWi5QhzUL",
			"1DT5djphMGn95wzftze6BjJeurM1KcAaVS",
			"14FCsRiFTHuraupmhbrnLq2y2W7JtAHJTc"},
			addrtype:	[]uint8{txP2pkh,txP2pkh,txP2pkh,txP2pkh},
      nvin:     1,
		},
		{
			name:		  "mainnet simple with 3* output",
			txhash:   "65a3235950069fe3f8cb5428ac72960ec0ea8ed09f77847a6622fff249c9d967",
			rawtx:	  []byte("01000000017d171d1fb676322e54684890fa2963dee483e024db6049319565841974fcb460010000006a47304402204c37cd1f21ac7016fd11e8b743a867a8f663d59f1796044f90950bea00429d5c022055e828cc7de3e0106139e30fc5a233c27fca2ff349f3e5f06cd56568dcc9db580121033dab59fca27b7f6ee99fb46fc8f8afe3e0841094ab1c55163644c2051a5a0885ffffffff02ecee0e00000000001976a914a028c390ffb66fe0164bb0a9154c4e76aafb988b88ac65dd17000000000017a9145d4674e8466725f60e8379a259948931ca73f5be8700000000"),
			encoded:	[]string{"1Fbqx2BqxNPcD6Gbnpd6XAGXpWSZ3P8xxq",
			"3ACDAoJhzHxLEcDgf5vxH2VeJftNY6NANY"},
			addrtype:	[]uint8{txP2pkh,txP2sh},
      nvin:     1,
		},
		{
			name:		"mainnet only 3* outputs",
			txhash:	"576534d167cf564fe395c16ebef382a1e9e553a10ffdae6eddee360606ee37f6",
			rawtx:	[]byte("02000000000101f817ede6389ce1633b04afe3ac2011e3b91a8646c2e3adfc424cb4819bc7b5be01000000171600146cb6b34581cd008dd0d3248b6e9ff2b4285c1019feffffff02fa2ec2000000000017a91444b6c8fa2b9d3e646c70887651f2861e2cb4ff7f87f18705000000000017a914f6902c9bd92f0c089a9f94341d7aba56634534668702483045022100b6c42756abf43ea7e7d27c1b5021e39d928d78f452758b6b7cfae49438d7a22202206f2b00b932133db3be1236498bc6cea4c42dd7fc585680631a4dc17d88771dba012102e91f9ee10015228fa9d79d7599720e9a637553b219d42742677d7500e2b4b603d84d0800"),
			encoded:	[]string{"37xLqiZQYD4WXg2BaVgYfBS2NeTdHwtcuY",
			"3QAis8tJTf8Ahr29jYM5Znrb4oE9u9TUNs"},
			addrtype:	[]uint8{txP2sh,txP2sh},
      nvin:     1,
		},
		{
			name:		"mainnet bench32 output",
			txhash:	"1e52891604b9ec5ccec8d94a404b3d08ef2cc6195137c7437d28c1fc52e3b663",
			rawtx:	[]byte("0100000000010181ec646c82f7e9a5a8d9adfb46162415172e891c37c4412ee60f8aaefd599dda0100000000ffffffff04200b2000000000001976a9147db7cf2faff951ee575ce8cbc2afadfc6f83986488acc0c62d00000000001976a9147dde2cdaef4c0be4149618fea926c50ee0e5016388ac70991400000000001976a9146b5dcb791a7ecc715aa4f30a6d24abb06351228d88ac694db50000000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d0400483045022100a16308296bbb33b92f82ef1d4b75b3b18a4aea8d74899476b80016d1ec2e107d0220756c99a26759c2e4417014030176fb59d022566562744fbb2eb33d70892a2bb201483045022100f39bda511ee4e07918bc3d16c903e3d0837672182682b2b18491a2b2d4f099a802205e85688ac3d8332353f29c8b79dde7e7e3a28494b648f1edf3c2dc6e5858827d016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000"),
			encoded:	[]string{"1CTjfGuuquBv2BFCWXeMGzKanWHXJ2FQr6",
			"1CUXcz3DQML6WmgNmF1iA4Bc7Mh99k2acD",
			"1AnhfWELiznEuFnKLQ2E1b3z99cjstkBfg",
			"bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej"},
			addrtype:	[]uint8{txP2pkh,txP2pkh,txP2pkh,txP2wsh},
      nvin:     1,
		},
		{
			name:		"mainnet OP_RETURN 1",
			txhash:	"8bae12b5f4c088d940733dcd1455efc6a3a69cf9340e17a981286d3778615684",
			rawtx:	[]byte("0100000001c858ba5f607d762fe5be1dfe97ddc121827895c2562c4348d69d02b91dbb408e010000008b4830450220446df4e6b875af246800c8c976de7cd6d7d95016c4a8f7bcdbba81679cbda242022100c1ccfacfeb5e83087894aa8d9e37b11f5c054a75d030d5bfd94d17c5bc953d4a0141045901f6367ea950a5665335065342b952c5d5d60607b3cdc6c69a03df1a6b915aa02eb5e07095a2548a98dcdd84d875c6a3e130bafadfd45e694a3474e71405a4ffffffff020000000000000000156a13636861726c6579206c6f766573206865696469400d0300000000001976a914b8268ce4d481413c4e848ff353cd16104291c45b88ac00000000"),
			encoded:	[]string{"636861726c6579206c6f766573206865696469",
			"1HnhWpkMHMjgt167kvgcPyurMmsCQ2WPgg"},
			addrtype:	[]uint8{txOpreturn,txP2pkh},
      nvin:     1,
		},
		{
			name:		"mainnet OP_RETURN 2",
			txhash:	"db0ea7f6b294afd5965b7b25d49814ff77a576c7ea3a71ccfabb042417bb5c45",
			rawtx:	[]byte("0100000001595081f57893476892696add35c72a767540c47e5f348fdda06381fa5483329a000000008b483045022100a022c5e1511f343e6b474408d8805eabc053b97ea880c26f5af2b7f8fca7e8c102204ec552fd10380e4f09585a4d2993651d07a1fc8c7b276f1ba7427aecc9e34f63014104ed87d130c464cd64f55618f7d44538b8746b16e42a517386e60fb9602cb4fe828cbcb4c0ab7315d0b4e0c5127c924430579a4f1cd8b9f1c761f8b0a889f608e4ffffffff0388080000000000001976a914dc552c27adf78437b51b7ba89fc6c6c106506b2a88ac0000000000000000166a146f6d6e69000000000000000100000000173eed8022020000000000001976a914fa0692278afe508514b5ffee8fe5e97732ce066988ac00000000"),
			encoded:	[]string{"1M61h6F3gcx8SnhqsKcYxauj7CJChwQNjY",
			"6f6d6e69000000000000000100000000173eed80",
			"1Po1oWkD2LmodfkBYiAktwh76vkF93LKnh"},
			addrtype:	[]uint8{txP2pkh,txOpreturn,txP2pkh},
      nvin:     1,
		},
		{
			name:		"mainnet multisig",
			txhash:	"e834fccee45aa04a0658f1dd363abdfc6902a17014952eacf3e658caa3030ea9",
			rawtx:	[]byte("01000000076477ec238c6a1a4d2657ff125768443b0d640fa1d293d9d8bbf1022b20f5eb2c00000000fdab0100473044022056c6b2b52df8ed27c23fc7ed0d5f30eb15a2a43ed28fe1014588f331bb9520fd02204ecf50c114edaf77be197a7560e0ad3f39fffac7e962cc4c56febdee48d52ace014730440220105102ba67659b4b4b3142051e65881991df6ccbeb5f4adf7e0ac55de6281742022060402c03db698ad76ad6cd7669f6ee4cae2ae2e37a2e06019c9452be55721bb301483045022100945f3b102b1cdcade2427bec332cc7caeda4deb34793c48f93ccf5f41a4f4bed02204c061162c76c3f88750557bf85c4c526ad8dab33dec0768e55c194726557eb8a014ccf5321021ecd2e5eb5dbd7c8e59f66e37da2ae95f7d61a07f4b2567c3bb10bbb1b2ec95321023bd78b0e7606fc1205721e4403355dfc0dbe4f1b15712cbbb17b1dc323cc8c0b2102afa49972b95496b39e7adc13437239ded698d81c85e9d029debb88641733528d2102b63fe474a5daac88eb74fdc9ce0ec69a8f8b81d2d89ac8d518a2f54d4bcaf4a52102fb394aaf232e114c06b1d1ca15f97602d2377c33e6fe5a1287421b09b08a5a3e2103fedb540dd71a0211170b1857a3888d9f950231ecd0fcc7a37ffe094721ca151f56aeffffffffad93c468a4ff920899b79053ce4b0a5d5c52af99e6815e6817880adb64214a8000000000fdac0100483045022100c0fbbee05a2570db503c80bd4c6849a3d4dd120dd4f8fcd1a799b5e97ec8c8a802201e61b5683a15a3a9141fdd4158d2131e6da5a50a1019a6af46ea3902936f39c7014830450221008e1d19a8eb2611ff12c7b052b1e37401164bc97324493f6547677ab303c8beec02201156688dd3102196073f43ff4656956ce60996724ef2aaa72713d6643d357e8b014730440220184a595ed5fb1d4ca899b784b95fe649745ab68a2c8fc8d86e1f3fcf617107da02200ebf95f8b20415a3f4809336cbdf63422bc10bc12c5e2d6d1659d67631730a29014ccf5321021ecd2e5eb5dbd7c8e59f66e37da2ae95f7d61a07f4b2567c3bb10bbb1b2ec95321023bd78b0e7606fc1205721e4403355dfc0dbe4f1b15712cbbb17b1dc323cc8c0b2102afa49972b95496b39e7adc13437239ded698d81c85e9d029debb88641733528d2102b63fe474a5daac88eb74fdc9ce0ec69a8f8b81d2d89ac8d518a2f54d4bcaf4a52102fb394aaf232e114c06b1d1ca15f97602d2377c33e6fe5a1287421b09b08a5a3e2103fedb540dd71a0211170b1857a3888d9f950231ecd0fcc7a37ffe094721ca151f56aeffffffffd425bc33d111779b01ecb735a4073f23d9182d1727cfc4031649a573da43cffe00000000fdaa010047304402203a0cf052d16dbbf201feb75db45ed91dad127e3b217ab7efdd3807ccbad46790022064b865af424b63270e8c6e451e7d1b3986123cc508d56061f4671a58f9c9bd700147304402200e4be57c41faaca855dd6ea8db7975fbd9b3378be47fa4094472362abdb688c7022044c03518f84e044655a76903679f3e61ef857ff6181e8b9f17eebfcc45da9ef70147304402201c60e5b38f1846a1786380b7e24ce2b646bd27b50a8ca0ad9e00325c13ce75c802207f9ccc705e1599ecc86dbf681602256891cca9c754dae88b07a359e38e4c048e014ccf5321021ecd2e5eb5dbd7c8e59f66e37da2ae95f7d61a07f4b2567c3bb10bbb1b2ec95321023bd78b0e7606fc1205721e4403355dfc0dbe4f1b15712cbbb17b1dc323cc8c0b2102afa49972b95496b39e7adc13437239ded698d81c85e9d029debb88641733528d2102b63fe474a5daac88eb74fdc9ce0ec69a8f8b81d2d89ac8d518a2f54d4bcaf4a52102fb394aaf232e114c06b1d1ca15f97602d2377c33e6fe5a1287421b09b08a5a3e2103fedb540dd71a0211170b1857a3888d9f950231ecd0fcc7a37ffe094721ca151f56aefffffffff892e0271e63b9bbe4d08b3e0dcf8b2becee0c362bd4175c809dc4aa3186636c00000000fdad0100483045022100e232eb7f5c5ca84afb732af8c75f32a9f05bb200863dca3b553057cd2d83e9080220260916967d80e82e308c88d82b3e6ab1ad1459d11ca58c9b5864678d4420ad7401483045022100a21ad8da4d6b4a27bbe14b7249dd29b601ad453b4a5a7217877ef9aa0acac6d102203072c77ba3237d256af04a4dccc1fcb5b35cec35eb272f3432efe2dd9edbdf9001483045022100b75cb7b1f48c86097b75b0120c330f5d927dbda5b37c982015b3486f847c877d02203610616d7473dfd0194e10fbb01370e714db95476669a70dbef31e18bd3d8688014ccf5321021ecd2e5eb5dbd7c8e59f66e37da2ae95f7d61a07f4b2567c3bb10bbb1b2ec95321023bd78b0e7606fc1205721e4403355dfc0dbe4f1b15712cbbb17b1dc323cc8c0b2102afa49972b95496b39e7adc13437239ded698d81c85e9d029debb88641733528d2102b63fe474a5daac88eb74fdc9ce0ec69a8f8b81d2d89ac8d518a2f54d4bcaf4a52102fb394aaf232e114c06b1d1ca15f97602d2377c33e6fe5a1287421b09b08a5a3e2103fedb540dd71a0211170b1857a3888d9f950231ecd0fcc7a37ffe094721ca151f56aeffffffffc2f6494c0bdd2a60c7fea7b77f3f3c4504ff9a39db482175bf7a38288910e37b00000000fdaa010047304402206f45ec261361d0016c55b0e0a40a6dcf7c8aacaecda5c09fe0b9db545e9361cb02201fab320c1f6bf3a4f47964b4a5832f5a246c8c5cd2adafdb5bd0b11c593f17560147304402207310de246ea207f75da7b46f8dc9dce48b64c2d884757edb7450dccf35a594290220055ac9c3936f563306c5db7ff1b3f894b024e39d0f4f1273c6f816361ea88e7e014730440220731b230d66c09dc96d07763547f902a5e7fa5d088d6a11b9be4ac69d17455acb0220372ddee0e6c13c78ed563978c97df41038206d999a72836ad2aba4fc132908f3014ccf5321021ecd2e5eb5dbd7c8e59f66e37da2ae95f7d61a07f4b2567c3bb10bbb1b2ec95321023bd78b0e7606fc1205721e4403355dfc0dbe4f1b15712cbbb17b1dc323cc8c0b2102afa49972b95496b39e7adc13437239ded698d81c85e9d029debb88641733528d2102b63fe474a5daac88eb74fdc9ce0ec69a8f8b81d2d89ac8d518a2f54d4bcaf4a52102fb394aaf232e114c06b1d1ca15f97602d2377c33e6fe5a1287421b09b08a5a3e2103fedb540dd71a0211170b1857a3888d9f950231ecd0fcc7a37ffe094721ca151f56aeffffffff94b8a3a02452c056777119d07c626b5d09a89ec02d413c55e0728496b9dccf1d00000000fdaa010047304402202708731dc8f09e929647caee64034462712059b81657a9ba52984413678e50ab0220689b03b7edda7739a535031a2446ed9adb6902785dd046fb5c5ee6a66f8fc90801473044022035bd4ba4400c8eb7df82459e3c774effa159d7517b6e404ef127dac1ee70008e02206ca8164bc51d3fffe20dceb75278f1d58e067ac9caece5f260b78e1e97737ffe01473044022053394edf8f9a12f693651462c5838d0943fc7b084da46503a273755382222a55022072005367df4ca7feaeaa25ea75cd1aa414fe6414d696ba17d64d17d862d8a54d014ccf5321021ecd2e5eb5dbd7c8e59f66e37da2ae95f7d61a07f4b2567c3bb10bbb1b2ec95321023bd78b0e7606fc1205721e4403355dfc0dbe4f1b15712cbbb17b1dc323cc8c0b2102afa49972b95496b39e7adc13437239ded698d81c85e9d029debb88641733528d2102b63fe474a5daac88eb74fdc9ce0ec69a8f8b81d2d89ac8d518a2f54d4bcaf4a52102fb394aaf232e114c06b1d1ca15f97602d2377c33e6fe5a1287421b09b08a5a3e2103fedb540dd71a0211170b1857a3888d9f950231ecd0fcc7a37ffe094721ca151f56aeffffffffd378ede06c0568060b56ec690011db0ca07e973c3bbd033a8251d938685f37a600000000fdac0100483045022100e8c6ecb79c1794fe15fb63fe2353b6fae94e6e205d244dc53d513fe54bbf854002207bdd2791251dc343481067ec16ff5af5c43e9045e35288bcc9d23e85a83c969a01473044022030be5d02ac4b8f886dfa4aa302d7f7a6ea581c98d6de68dd1bda99cb834d3aaa022074a629f6c0efc7703de76c2f33c16be3f8d154844c3cc51876bb44dbec3500c001483045022100805053c6136ad7a61eed898a479eed50ec7fe832d73f80058fade3a5be4fcee302204394c91c418eb0c3d8168aef555ecd4a5eb42d0cddc8a64c7ca9fd7fc23f0e1f014ccf5321021ecd2e5eb5dbd7c8e59f66e37da2ae95f7d61a07f4b2567c3bb10bbb1b2ec95321023bd78b0e7606fc1205721e4403355dfc0dbe4f1b15712cbbb17b1dc323cc8c0b2102afa49972b95496b39e7adc13437239ded698d81c85e9d029debb88641733528d2102b63fe474a5daac88eb74fdc9ce0ec69a8f8b81d2d89ac8d518a2f54d4bcaf4a52102fb394aaf232e114c06b1d1ca15f97602d2377c33e6fe5a1287421b09b08a5a3e2103fedb540dd71a0211170b1857a3888d9f950231ecd0fcc7a37ffe094721ca151f56aeffffffff02646fcc780400000017a9147c6775e20e3e938d2d7e9d79ac310108ba501ddb8700d0ed902e0000001976a914cebb2851a9c7cfe2582c12ecaf7f3ff4383d1dc088ac00000000"),
			encoded:	[]string{"3D2oetdNuZUqQHPJmcMDDHYoqkyNVsFk9r",
			"1Kr6QSydW9bFQG1mXiPNNu6WpJGmUa9i1g"},
			addrtype:	[]uint8{txP2sh,txP2pkh},
      nvin:     7,
		},
    {
			name:		"mainnet NEW",
			txhash:	"97d1b00fcef1f19531a19bb1722635341a9f2ad261ecf6eed89eca2cbd3bb3ee",
			rawtx:	[]byte("0100000001e507cb947464fc74540a9c197f815aa283ba9db74185ac08449c38491a8c34ac00000000fdfe0000483045022100ca17e5614fdf80c170b16f67c650046a40c7b7563b25aaad1dd08cd28c22141d02204359bdde6a171ef094f91173a36978f48d917a094c811aa2cee597cc8d6b9507014830450221008d86ee7406d98716c50579b3c7f171deeb2d5a065b00f84907aae9a2ed05220102205277477975e1d819ca05c4d05e7a575a5bb1bca0594a58173087394463b7914f014c69522102551ecd379bac7fe3374df7b50478301a26d34dfd4094d909ec6f9b0a40217d1c2103a7dd8cf968258d25dfef47adf1ef616ee10c77be3cb1f2fb1a7110856ba6a5a72102eab43b82bbe8c482abbdfd1b443084e48b8d5be232280af23ac8ad78e6e3591853aeffffffff02a0d21e00000000001976a914c4d2c1a1fa246e6d55610975565bfeedd83a1e5488ace0930400000000001976a914200031455ce0dd3ad265a2686a314352e67f58c388ac00000000"),
			encoded:	[]string{"1JwhvD6mwpRzAhVKihWcxXPJNeFKJGgHUh",
			"13vCrWerFRP2rULtNUpV5bTipDfubVtm7U"},
			addrtype:	[]uint8{txP2pkh,txP2pkh},
      nvin:     1,
    },
    {
			name:		"buffer overflow (unable to parse)",
			txhash:	"9740e7d646f5278603c04706a366716e5e87212c57395e0d24761c0ae784b2c6",
			rawtx:	[]byte("010000000121eb234bbd61b7c3d31034762126a64ff046b074963bf359eaa0da0ab59203a0010000008b4830450220263325fcbd579f5a3d0c49aa96538d9562ee41dc690d50dcc5a0af4ba2b9efcf022100fd8d53c6be9b3f68c74eed559cca314e718df437b5c5c57668c5930e1414050201410452eca3b9b42d8fac888f4e6a962197a386a8e1c423e852dfbc58466a8021110ec5f1588cec8b4ebfc4be8a4d920812a39303727a90d53e82a70adcd3f3d15f09ffffffff01a0860100000000006b4c684c554b452d4a522049532041205045444f5048494c4521204f682c20616e6420676f642069736e2774207265616c2c207375636b612e2053746f7020706f6c6c7574696e672074686520626c6f636b636861696e207769746820796f7572206e6f6e73656e73652eac00000000"),
			encoded:	[]string{""},
			addrtype:	[]uint8{txParseErr},
      nvin:     1,
    },
    {
			name:		"mainnet bench32 ouput 2",
			txhash:	"5841ff53611ce55facbc57d18c0563576af9e5453f2dd1406f4324a0cee02a18",
			rawtx:	[]byte("0100000000010106325bac2f2e7ca67fa46c8304fb3b747e5578df1eef0394349ce2cdd744f7f10100000000ffffffff02db355202000000001976a91489ea1263056ac068adba4844efb376a3a19635ad88ac43b72f0700000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d0400483045022100b4c3ac1d0d785a75d7e0e21b1054f426deac0604635bef79010cc8c961bddec1022043ceac1de07f7011b1c922afc46c19a4aff1d6c5ec24c035330472f6de973f7c0147304402201ed600cde0e2ef4b48b4be8144b26cf91ca62a778c89a27cb8340d99551fbe8b02207fdc4eedb12aeba6fc1eea500000e39dda02487c01837eb78fad5c5e6de2d88e016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000"),
			encoded:	[]string{"1DaDyspUHk5GUFTWRXUogitzSKKf47ecFp",
			"bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej"},
			addrtype:	[]uint8{txP2pkh,txP2wsh},
      nvin:     1,
    },
    {
			name:		"mainnet (long) block reward",
			txhash:	"8a3357ccbcd82d309e1788757fcbccc4803d1570051e242226625b5b029cbc63",
			rawtx:	[]byte("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff510346d7030d00456c69676975730052089d66fabe6d6d779ca5b9a3af6c66c21db522bb5577722b3d526cd525e298e329b2fdb1b1c2670400000000000000002f7373322f005011e2090000000000001deaffffffff2e0f7f0c12000000001976a914cd9be734639e0738692e569412a9fc65ce249fae88ac16bf0a0b000000001976a9148d72578af4c8094ce919789a9201003dfd461aeb88ac97ec3f0a000000001976a914001eecd0f4102ce18b9c4bf1c1c547ddad43a71988ac00018909000000001976a914217a637f55892e4327694184dd0d20536186236888ac57b43907000000001976a9148d20ff772a9ac36f789b96e2fef13ac08b21f63b88acaa552807000000001976a9146c8a135febf4ab61458ed2d3de1c74b1c2749f5d88ac24d12307000000001976a914dd98b5c3f3be71942db8b4bec0cd2502ec7d9f2588ac28a91007000000001976a914534764fb2f6e4fcab01789036135c5b6eb1ffd1e88ac1fb70507000000001976a9143fd9ec42d8f5bf65faf11cd1297ed5dbd0e91ab888ac87fb0307000000001976a9149538a23567d9233e3268cd461df9389528c1308b88ac0365f005000000001976a914857c240b87bb0ddce6d8e5ffaffdd7a70be7e4ea88acbfc5ee05000000001976a91483afb6e18fed3c86278f18d4d78afad8c504580988ac2fc20905000000001976a914a6183c12a78e2a318f831a99e887689b69b927b188acc8628604000000001976a91475bed242e9b2cf5b4d45ef9f052be369067eae6888ac45bb9303000000001976a9142e5aa66d38ca6c395969b1aac1f614349221124388ac06451603000000001976a914f344947d74437ff7ccafbdf8c2b9b2d0e5d5b21c88ac5aa65902000000001976a914d2d9d17985077c429493b73cd6d9e9ae7bf30e5688ac68334802000000001976a914f1e57140974995ddaa75c2d996c89293816f13e788acda071002000000001976a914472427e6c7bbf78ad5ace308f94315b4d27c313e88ac5010ba01000000001976a91498a044d09bd08cc51a915333b8d3a45bb4e6b73888ac9d536e01000000001976a91471e8822a58e9a693f90c277c1314dc718528176988ace4350e01000000001976a914226c1eee79f5ed89ea65ed078751af96c14e7ac388ac57100a01000000001976a91417ad8ff03859468be931fcb934af69919f5abbbb88acb4ac0401000000001976a914aaf382157ae5cf1a570dd8ff901c4af20e638bd588ac9dc30301000000001976a914553d1ac6a4070c0e51ad68e9eb786f81cc3299bb88ace6fd0201000000001976a914a975dcbdef5c9360d831963bb3923411a495c0d788ac3cf10201000000001976a9149032d26c5176cb8f6316e872074a3d9a0a10624f88acf2070201000000001976a9144500d06d486b08090cf9ffcc651b8072b57ec14188acdf160101000000001976a9146d8096de881bca5fd3385655376b868e3201ef2488ac2cadff00000000001976a914b6a067391eb2d6eab154b40cceb8b2ac783e9acc88acd953ff00000000001976a91406b0a951d0ec510fd050add21b9ca2c64124b13e88ac08f1fe00000000001976a9148996ba6c437ed65cc40f465c36cf72775bc3902888acc980fe00000000001976a9140731cf6b5e608dceb61349ee4eb79f1df963de6788acb169f600000000001976a91481c5e2e5869d49b0ac9592f48561004c96aafe7f88aca743f000000000001976a914a46ab0b1936b2e087ac201f9fe366618c17a58dd88ac2d50e300000000001976a914cfda421589a13697ab3404416a813e474a63bab088ace0a2c600000000001976a914b5ba7443e2b051d6152891fc8aa8924623a71fee88acf325a800000000001976a9147ddba907674dc4d627c1d589972c2876b755de3a88acf5309400000000001976a914ec8b0add810bd1533c8c3e38da2b1b10fe9f08f788ac35859100000000001976a914907204535ba24f5c9669acfb8a56ba9562cc99f288acbb558700000000001976a914d59149b012bb32c64dfb6f76daba3335249ac5f888acfed57200000000001976a914e3e46a767f0b8e55534cc845f7de4d373440318188ac55b67200000000001976a91418181f50b1b1eb3abcae05883f9b6485be4e522988ac673a7000000000001976a9149054471ea5771a3e95178e637911e9ad68edd95988ac624a6200000000001976a9147e516766bb50f3cd0301a7c87505ac4b006d70a388ac01000000000000001976a9145399c3093d31e4b0af4be1215d59b857b861ad5d88ac00000000"),

			encoded:	[]string{"1KkAHfbpzXUaYdVX85kGuxfdGWENk2Zpy6",
			"1DtuFfgtaRYrhwCt6KpGVcwBp6UocTf2My",
			"11e3iF3q3Zhm8jyx1nCRUdT7EyCjgoViM",
			"1441v3Y1NYx4rCi6kRwfyYW4XAhsGfpTKG",
			"1DsDonRjgiuVDJCcB4MpZjB4sNrvCnhojZ",
			"1AtuPPHjwesiEEN9DonRFvucjstYFtJc4U",
			"1MChGsd9FpVsoiNgEdwrSqhWESfe1qLuao",
			"18bLcVkviErQi75zB8X39jZXxHNpSZggdC",
			"16pcddk8c6MFKyQvL2NvuWEguGw1R8B3VL",
			"1Ec1XZFtRLE9ZchGeAGNVSvmGZcbMMwafU",
			"1DAoay7hfohT9TdUjB4gf2JNj8dTnhLHoP",
			"1D1J1yn16eJ9hz8VSEbt88mZ8S3Zb67gTy",
			"1G9EChTTs7LPXQqXrg424HkNDsYJJFDboG",
			"1Bjaewgk3bXujiJUqLcrkk7aKSS1f9kUV3",
			"15E6gJdA4ivRWFw7uLxAq9Tvt9cCDwKKRr",
			"1PBHMe68qdinR9692eWoozLuJRFSjRdiqv",
			"1LDsqZd5XNqFquWh1yT14UbjsNastvGiof",
			"1P42iAk97dB8qAv96crLFiDVVUJSm38aHV",
			"17VAHms6irdfFycJG9Ma7vFzhNgcXcAywh",
			"1Ev1hq9uSQZoPe7AGYExp5NuZqPufQX6iD",
			"1BPHtea52ATP9N45dio2S7ZgFCANekJaKF",
			"1491VtZRuZFuYaKBaR3V9Xw6TVakeBTu6j",
			"13ACTjJK6Bcbzf4bTob5fL526YBZ5gYKZz",
			"1GauarVjtcWeiDCFpojf3Z8FVHt3i6V7wX",
			"18mhdvn1mdESTzT1TMtViCGtAbWQpZYCqZ",
			"1GT2PdkbK3TkoB522aMGQHbd8ARvXhLT8Z",
			"1E9TBo9X1yRhoE3EWUctM53bqfGUWXgxSE",
			"17HrbojD3C99PbxgcpCYEFu4cE9RhqyArv",
			"1AyzhTGvhqpC7CUXJHJ7eho3HZBVei2oDN",
			"1Hee8vb5QkMPg1d6kVqtkBxjpVPhsTicVo",
			"1cNgoS7MzNzPgM9AguXtFZdKPkCqmh9Rd",
			"1DYW94XE6HeKLJgzBarMyKw1irUX3Jx7ik",
			"1f3QDfjRv3sVvHrpmULEs6xAH8cYtyfqT",
			"1CqBE8SfT69aH2PTLvjGjxmy31Moatg6MQ",
			"1FzMdRhm2Bm2WxVo1wbgTnZTJvVuDCaRbr",
			"1Kx2LYvYmrg311zcdLHpzPca4Rh133nDne",
			"1HZtfnSm2aJrVVH6zXTA6RfM4GASk1LT8V",
			"1CUUcF2fNS38VakmJh5QPSrzN9q6TnvjR9",
			"1NZj37PrVPE7Cmcxznct2am8BkMLbiSnP1",
			"1EAktf5LhyWNLWzfmEkD24savdzc4hwCZ9",
			"1LUEydu72PiEqEYFtFaahRiHaPQabQ9j3u",
			"1Mmz2ABwmKZaG3YvJWSB8hjJkAaVnBiYMG",
			"13CQ7eqsF7h6dT2oRmaQMAPLQ6Evy8tjSL",
			"1EA9GMDbhRMj5fPPYXpg4Sfz1WovhbBWMv",
			"1CWufDPNkmWLjdLgRHSkLQkVL1WBqSCj7e",
			"18d3HV2bm94UyY4a9DrPfoZ17sXuiDQq2B"},
			addrtype:	[]uint8{txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh,txP2pkh},
      nvin:     1,
    },
	}

	for _, test := range tests {
		rawtx := serial.HexToBinary(test.rawtx)
		rawTx, err := parser.New(rawtx)
		if err != nil {
			t.Fatal(err)
		}
		tx := DecodeTx(rawTx)
		txHash := fmt.Sprintf("%x", serial.ReverseHex(tx.Hash))
		if txHash != test.txhash {
      t.Errorf("%v: Transaction hash does not match: %v != %v",
			test.name, txHash, test.txhash)
		}
    if (int(tx.NVin) != test.nvin) {
			t.Error("Wrong number of input. Should be tx.NVin =", test.nvin)
    }
		if (int(tx.NVout) != len(test.encoded)) {
			t.Error("Wrong number of output. Should be tx.NVout =", len(test.encoded))
		}
		if len(test.addrtype) != len(test.encoded) {
			t.Errorf("Error in tests: %d != %d", len(test.addrtype), len(test.encoded))
		}
		for i, vout := range tx.Vout {
			if vout.AddrType != test.addrtype[i] {
				t.Errorf("%s: Wrong TxType: %d != %d", test.name, vout.AddrType, test.addrtype[i])
			}
			decoded, err := DecodeAddr(vout.AddrType, vout.Hash)
			if err != nil {
				t.Error("DecodeAddr: ", err)
			}
			if test.encoded[i] != decoded {
				t.Errorf("%v: String on decoded value does not match expected value: %v != %v",
				test.name, decoded, test.encoded[i])
			}
		}
	}
}
