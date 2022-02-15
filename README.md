# Algorand HDWallet

This a Algorand HDWallet implementation trial, complying with [bip32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki) & [slip-0010](https://github.com/satoshilabs/slips/blob/master/slip-0010.md) for ED25519 curve.

Also, basic **deterministic** key derivation from **25 words** mnemomic is implemented too using native [go sdks](https://github.com/algorand/go-algorand-sdk). (Thanks to Yieldly repo [here](https://github.com/yieldly-finance/yieldly-deterministic-account-generator/tree/main/src/go))


### Key Notes
* Only hardened public keys are supported, because normal/public derivations (extended public key to public key) is not possible for ED25519 curve.
* Algorand coin code is **283** according to [slip-0044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md). So derivation paths may be `m/44'/283'/0'/1'/100'` or simply `m/1'/0'/0'` (Not mandating any format...)
* ed25519 test vectors 1 & 2 defined with SLIP-0010 are implemented. 

---

## Usage

```go
package main

import (
        "fmt"
        "github.com/polarbit/hdwallet"
)

func main() {
        mnemonic24 := "segment inhale symptom olive cheese tissue vacuum lazy sketch salt enroll wink oyster hen glory food weasel comic glow legal cute diet fun real"
        path := "m/44'/283'/0'"
        a, _, _ := hdwallet.DeriveAddressFromMnemonic(mnemonic24, path)

        fmt.Printf("Address: %s\n", a)
}
```

---

## Notes

### PublicKey & Algorand Address
*slip0010* creates extended public keys padded a zero byte (0x00) from left. \
This 33 bytes public key is not compatible with Alogrand address generation (32 bytes required). \
 So while creating Algorand an address, we ignore the padded *zero* byte.

### PrivateKey
*ed25519* key generation produces a 64 byte private key. \
Second half of this private key is public key of the curve, `pub = priv[32:]` . \
slip0010 uses the first half (32 byest) of that for key generation \
Algorand-Sdk also uses the first half of that private key for mnemonics.

### kmd - Deterministic Address Generation

[This](https://github.com/algorand/go-algorand/blob/04e69d4153d0e67d477d4d4b12faede7ec5331b1/daemon/kmd/wallet/driver/sqlite_crypto.go#L234) is the Algorand kmd code that generates deterministec wallets. I expect same derivations from same seed; but did not tested myself.

### How transactions are signed by Algorand-SDK

Below are two main/low-level methods that helps to understand low levels of Algorand tx signing.

```go
// rawSignTransaction signs the msgpack-encoded tx (with prepended "TX" prefix), and returns the sig and txid
func rawSignTransaction(sk ed25519.PrivateKey, tx types.Transaction) (s types.Signature, txid string, err error) {
	toBeSigned := rawTransactionBytesToSign(tx)

	// Sign the encoded transaction
	signature := ed25519.Sign(sk, toBeSigned)

	// Copy the resulting signature into a Signature, 
    // and check that it's the expected length
	n := copy(s[:], signature)
	if n != len(s) {
		err = errInvalidSignatureReturned
		return
	}
	// Populate txID
	txid = txIDFromRawTxnBytesToSign(toBeSigned)
	return
}

// rawTransactionBytesToSign returns the byte form of the tx that 
// we actually sign and compute txID from.
func RawTransactionBytesToSign(tx types.Transaction) []byte {
	var txidPrefix = []byte("TX")

	// Encode the transaction as msgpack
	encodedTx := msgpack.Encode(tx)

	// Prepend the hashable prefix
	msgParts := [][]byte{txidPrefix, encodedTx}
	return bytes.Join(msgParts, nil)
}
```
[go-algorand/daemon/kmd/wallet/driver/sqlite_crypto.go](https://github.com/algorand/go-algorand-sdk/blob/f09c24dcd1866f04ee84da89e60d12495b388a9b/crypto/crypto.go#L97)

---

## Tests

Run single test: `go test ./... -run "MnemonicToSeed" -v`

Run integration test: `TEST_INT=true go test ./... -run "AlgorandPayment" -v`

---

## References 

[SLIP-0010](https://github.com/satoshilabs/slips/blob/master/slip-0010.md)

[BIP-32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)

[BIP-44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)

[Yieldly - Algorand Deterministic Account Derivation](https://github.com/yieldly-finance/yieldly-deterministic-account-generator)

[Algorand Developer Docs - Offline Signing](https://developer.algorand.org/docs/get-details/transactions/offline_transactions/)

[Algorand Go SDK - crypto.go](https://github.com/algorand/go-algorand-sdk/blob/develop/crypto/crypto.go)

---

## To Do

- Refactor: logs, types ([]byte),  errors (better messages).
- Add TX building, signing and broadcasting capability. (We already do this in payment integration test)
- Import & export xpriv, and generate keys from it.
- Use environment variables in integration tests for mnemomics.
- ... support other curves when Algorand is done. (maybe no or just ed25519 coins!)