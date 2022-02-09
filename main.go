package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/polarbit/algorand-hdwallet/algorand"
	"github.com/polarbit/algorand-hdwallet/bip32path"
	"github.com/polarbit/algorand-hdwallet/hdwallet"
	"github.com/polarbit/algorand-hdwallet/utils"
)

const testVector2Seed = "000102030405060708090a0b0c0d0e0f"
const testVector2Path = "m/0'"

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

	wallet := hdwallet.HdWalletOptions{
		Curve:           hdwallet.CURVE_ED25519,
		FnParseMnemonic: algorand.GetSeedFromMnemonic,
	}

	mkey, mchaincode, err := wallet.GenerateMasterKeyFromHexSeed(testVector2Seed)
	if err != nil {
		panic(err)
	}

	key, chaincode := generateAccount(accountPath, mkey, mchaincode)

	fmt.Println("chaincode: ", hex.EncodeToString(chaincode))
	fmt.Println("private: ", hex.EncodeToString(key))

	hexPub := getPublicKey(key)

	fmt.Println("public: ", hexPub)

	mpub := getPublicKey(mkey)
	mpubbytes, _ := hex.DecodeString(mpub)
	keyID, _ := utils.Hash160(mpubbytes)
	fingerPrintPub := hex.EncodeToString(keyID[:4])
	fmt.Println("fingerprint: ", fingerPrintPub)
}

func generateAccount(path *bip32path.Bip32Path, masterKey []byte, masterChainCode []byte) ([]byte, []byte) {
	key, chainCode := masterKey, masterChainCode

	for _, s := range path.Segments {
		key, chainCode = hdwallet.CKD(key, chainCode, s.Value)
	}

	return key, chainCode
}

func getPublicKey(key []byte) string {
	priv := ed25519.NewKeyFromSeed(key)
	pub := priv.Public()

	buf := make([]byte, 32)
	copy(buf, pub.(ed25519.PublicKey))

	buf = append(make([]byte, 1), buf...)
	hexPub := hex.EncodeToString(buf)

	return hexPub
}
