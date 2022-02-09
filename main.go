package main

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/hkdf"

	mnemonic "github.com/algorand/go-algorand-sdk/mnemonic"
)

const testMnemonic = "mask wear topple mixture steel cupboard gain satoshi chuckle analyst spoil borrow melody punch start ivory resource olympic sibling conduct stairs manual curtain absorb citizen"
const testVector2Seed = "000102030405060708090a0b0c0d0e0f"

const (
	ED25519_CURVE   = "ed25519 seed"
	HARDENED_OFFSET = 0x80000000
)

func main() {
	fmt.Println("Hello, world.")

	testIt()

	m, err := generateMnemonic()
	if err != nil {
		panic(err)
	}

	fmt.Println("\n\nGenerated Mnemonic")
	fmt.Println(m)

	fmt.Println("\n\n=== Generate Child Keys ===")

	for i := uint32(1); i <= 10; i++ {
		fmt.Printf("\nAccount %2d\n", i)

		addr, mnemonic, err := deriveAccount(testMnemonic, i)

		if err != nil {
			panic(err)
		}

		fmt.Println(addr)
		fmt.Println(mnemonic)
	}
}

func testIt() {

	key, chaincode := getMasterKeyFromHexSeed(testVector2Seed)

	fmt.Println("TestVector2 Seed")
	fmt.Println(testVector2Seed)
	fmt.Println("chaincode: ", hex.EncodeToString(chaincode))
	fmt.Println("private: ", hex.EncodeToString(key))

	genK := ed25519.NewKeyFromSeed(key)
	genP := genK.Public()

	bufx := make([]byte, 32)
	copy(bufx, genP.(ed25519.PublicKey))
	bufx = append(make([]byte, 1), bufx...)

	hexP := hex.EncodeToString(bufx)

	fmt.Println("public: ", hexP)

}

// bip32
func getMasterKeyFromMnemonic(mnemonic_ string) ([]byte, []byte) {
	key, err := mnemonic.ToMasterDerivationKey(mnemonic_)
	if err != nil {
		panic(err)
	}

	return getMasterKey(key[:])
}

func getMasterKeyFromHexSeed(seed string) ([]byte, []byte) {
	key, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}

	return getMasterKey(key)
}

func getMasterKey(secret []byte) ([]byte, []byte) {
	salt := []byte(ED25519_CURVE)
	hash := hmac.New(sha512.New, salt)
	_, err := hash.Write(secret[:])
	if err != nil {
		panic(err)
	}

	I := hash.Sum(nil)
	IL := I[:32]
	IR := I[32:]

	return IL, IR
}

// bip32
func ckd(key []byte, chainCode []byte, ix uint32) ([]byte, []byte) {
	vbuf := make([]byte, 4)
	binary.BigEndian.PutUint32(vbuf, ix)

	buf := make([]byte, 1)
	buf = append(buf, key...)
	buf = append(buf, vbuf...)

	hash := hmac.New(sha512.New, chainCode)
	_, err := hash.Write(buf)
	if err != nil {
		panic(err)
	}

	I := hash.Sum(nil)
	IL := I[:32]
	IR := I[32:]

	return IL, IR
}

func generateMnemonic() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return mnemonic.FromKey(key)
}

func deriveAccount(mnemonic_ string, ix uint32) (string, string, error) {
	info := []byte(fmt.Sprintf("AlgorandDeterministicKey-%d", ix))

	key, err := mnemonic.ToMasterDerivationKey(mnemonic_)

	keystream := hkdf.Expand(sha512.New512_256, key[:], info)

	pub, priv, err := ed25519.GenerateKey(keystream)
	if err != nil {
		return "", "", err
	}

	chksum := sha512.Sum512_256(pub[:])
	checksumAddress := append(pub[:], chksum[28:]...)
	a := base32.StdEncoding.EncodeToString(checksumAddress)[:58]

	m, err := mnemonic.FromPrivateKey(priv)
	if err != nil {
		return "", "", err
	}

	return a, m, nil
}
