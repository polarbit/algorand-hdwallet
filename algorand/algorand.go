package algorand

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"errors"
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

func MnemonicToSeed(mnemonic_ string) ([]byte, error) {
	key, err := mnemonic.ToMasterDerivationKey(mnemonic_)
	if err != nil {
		return nil, err
	}
	return key[:], nil
}

func DeriveBasicAccount(mnemonic_ string, ix uint32) (string, string, error) {
	info := []byte(fmt.Sprintf("AlgorandDeterministicKey-%d", ix))

	key, err := mnemonic.ToMasterDerivationKey(mnemonic_)

	keystream := hkdf.Expand(sha512.New512_256, key[:], info)

	pub, priv, err := ed25519.GenerateKey(keystream)
	if err != nil {
		return "", "", err
	}

	a, err := PublicKeyToAddress(pub)
	if err != nil {
		return "", "", err
	}

	m, err := PrivateKeyToMnemomic(priv)
	if err != nil {
		return "", "", err
	}

	return a, m, nil
}

func PublicKeyToAddress(pub []byte) (string, error) {
	if len(pub) != 32 {
		return "", errors.New("Invalid public key length")
	}
	chksum := sha512.Sum512_256(pub[:])
	checksumAddress := append(pub[:], chksum[28:]...)
	return base32.StdEncoding.EncodeToString(checksumAddress)[:58], nil
}

func PrivateKeyToMnemomic(priv []byte) (string, error) {
	return mnemonic.FromPrivateKey(priv)
}
