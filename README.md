# Algorand HDWallet

This a Algorand HDWallet implementation trial, complying with [bip32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki) & [slip-0010](https://github.com/satoshilabs/slips/blob/master/slip-0010.md) for ED25519 curve.

Also, basic **deterministic** key derivation from **25 words** mnemomic is implemented too using native [go sdks](https://github.com/algorand/go-algorand-sdk). (Thanks to Yieldly repo [here](https://github.com/yieldly-finance/yieldly-deterministic-account-generator/tree/main/src/go))


### Key Notes
* Only hardened public keys are supported, because normal/public derivations (extended public key to public key) is not possible for ED25519 curve.
* Algorand coin code is **283** according to [slip-0044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md). So derivation paths may be `m/44'/283'/0'/1'/100'` or simply `m/1'/0'/0'` (Not mandating any format...)
* ed25519 test vectors 1 & 2 defined with SLIP-0010 are implemented. 


## How To
Run single test: `go test ./... -run "MnemonicToSeed" -v`


## References 

[SLIP-0010](https://github.com/satoshilabs/slips/blob/master/slip-0010.md)

[BIP-32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)

[BIP-44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)

[Yieldly - Algorand Deterministic Account Derivation](https://github.com/yieldly-finance/yieldly-deterministic-account-generator)

---

## Implementation Notes

### PublicKey & Algorand Address
*slip0010* creates extended public keys padded a zero byte (0x00) from left. \
This 33 bytes public key is not compatible with Alogrand address generation (32 bytes required). \
 So while creating Algorand an address, we ignore the padded *zero* byte.

### PrivateKey
*ed25519* key generation produces a 64 byte private key. \
Second half of this private key is public key of the curve, `pub = priv[32:]` . \
slip0010 uses the first half (32 byest) of that for key generation \
Algorand-Sdk also uses the first half of that private key for mnemonics.

