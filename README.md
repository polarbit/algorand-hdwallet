# hdwallet

* HDWallet bip32/slip-0010 implementation for ED25519 curve.
* Also Algorand native **deterministic** key derivation from **25 words** mnemomic is implemented seperately
  * Derived accounts also has their own 25 words mnemonic.


Run single test: \
`go test ./... -run "MnemonicToSeed" -v`
