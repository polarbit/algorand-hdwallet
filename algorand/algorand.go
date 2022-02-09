package algorand

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"fmt"

	"github.com/algorand/go-algorand-sdk/mnemonic"
	"golang.org/x/crypto/hkdf"
)

func GenerateMnemonic() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return mnemonic.FromKey(key)
}

func GetSeedFromMnemonic(mnemonic_ string) ([]byte, error) {
	key, err := mnemonic.ToMasterDerivationKey(mnemonic_)
	if err != nil {
		return nil, err
	}
	return key[:], nil
}

func DeriveAccount(mnemonic_ string, ix uint32) (string, string, error) {
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
