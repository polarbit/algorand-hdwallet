package main

import (
	"encoding/hex"
	"fmt"

	"github.com/polarbit/algorand-hdwallet/bip32path"
	"github.com/polarbit/algorand-hdwallet/hdwallet"
)

const testVector2Seed = "000102030405060708090a0b0c0d0e0f"
const testVector2Path = "m/0'/1'/2'/2'/1000000000'"

const (
	ED25519_CURVE   = "ed25519 seed"
	HARDENED_OFFSET = 0x80000000
)

func main() {
	fmt.Println("Hello, world.")

	testIt()
}

func testIt() {

	accountPath, err := bip32path.Parse(testVector2Path)
	if err != nil {
		panic(err)
	}

	masterXKey, err := hdwallet.GenerateMasterKeyFromHexSeed(hdwallet.CURVE_ED25519, testVector2Seed)
	if err != nil {
		panic(err)
	}

	xkey, err := hdwallet.GenerateAccount(accountPath, masterXKey)

	fmt.Println("chaincode: ", hex.EncodeToString(xkey.ChainCode))
	fmt.Println("private: ", hex.EncodeToString(xkey.PrivateKey))
	fmt.Println("public: ", hex.EncodeToString(xkey.PublicKey))
	fmt.Println("fingerprint: ", hex.EncodeToString(xkey.Fingerprint))
}
