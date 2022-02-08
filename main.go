package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"fmt"

	"golang.org/x/crypto/hkdf"

	mnemonic "github.com/algorand/go-algorand-sdk/mnemonic"
)

const testMnemonic = "mask wear topple mixture steel cupboard gain satoshi chuckle analyst spoil borrow melody punch start ivory resource olympic sibling conduct stairs manual curtain absorb citizen"

func main() {
	fmt.Println("Hello, world.")

	m, err := generateMnemonic()

	if err != nil {
		panic(err)
	}

	fmt.Println("Master Mnemonic")
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
