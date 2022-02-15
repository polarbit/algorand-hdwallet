package algorand

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/polarbit/algorand-hdwallet/bip32path"
	"github.com/polarbit/algorand-hdwallet/hdwallet"
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

func DeriveHDWalletAccount(mnemonic string, derivationPath string) (a string, m string, xkey *hdwallet.ExtendedKey, err error) {
	path, err := bip32path.Parse(derivationPath)
	if err != nil {
		return
	}

	seed, err := MnemonicToSeed(mnemonic)
	if err != nil {
		return
	}

	xmaster, err := hdwallet.GenerateMasterKey(hdwallet.CURVE_ED25519, seed)
	if err != nil {
		return
	}

	xkey, err = hdwallet.DeriveAccount(hdwallet.CURVE_ED25519, path, xmaster)
	if err != nil {
		return
	}

	// !!! IMPORTANT !!!
	// 'hdwallet' creates extended public keys padded with zero byte from left.
	// But we need a 32 byte public key here; so we need the original one.
	a, err = PublicKeyToAddress(xkey.CurvePublicKey)
	if err != nil {
		return
	}

	// We can also use xkey.CurvePrivateKey here;
	// Because xkey.CurvePrivateKey[:32] == xkey.PrivateKey
	// Also xkey.CurvePrivateKey[32:] == xkey.CurvePublicKey
	m, err = PrivateKeyToMnemomic(xkey.PrivateKey)
	if err != nil {
		return
	}

	return
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
