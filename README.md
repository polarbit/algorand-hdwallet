# Algorand HDWallet

This a Algorand HDWallet implementation trial, complying with [bip32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki) & [slip-0010](https://github.com/satoshilabs/slips/blob/master/slip-0010.md) for ED25519 curve.

Also, basic **deterministic** key derivation from **25 words** mnemomic is implemented too using native [go sdks](https://github.com/algorand/go-algorand-sdk). (Thanks to Yieldly repo [here](https://github.com/yieldly-finance/yieldly-deterministic-account-generator/tree/main/src/go))


### Key Notes
* Only hardened public keys are supported, because normal/public derivations (extended public key to public key) is not possible for ED25519 curve.
* Algorand coin code is **283** according to [slip-0044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md). So derivation paths may be `m/44'/283'/0'/1'/100'` (But not mandated...)
* ed25519 test vectors 1 & 2 defined with SLIP-0010 are implemented. 


## How To
Run single test: `go test ./... -run "MnemonicToSeed" -v`


## References 

[SLIP-0010](https://github.com/satoshilabs/slips/blob/master/slip-0010.md)

[BIP-32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)

[BIP-44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)

[Yieldly - Algorand Deterministic Account Derivation](https://github.com/yieldly-finance/yieldly-deterministic-account-generator/tree/main/src/go)
