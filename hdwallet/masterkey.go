package hdwallet

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
)

type FnParseMnemonic func(m string) ([]byte, error)

func GenerateMasterKey(curve string, seed []byte) (xkey *ExtendedKey, err error) {
	if len(seed) < 16 || len(seed) > 64 {
		return nil, errors.New("Invalid  seed length")
	}

	// TODO:
	if curve != CURVE_ED25519 {
		panic("Curve not supported")
	}

	salt := []byte(curve)
	hash := hmac.New(sha512.New, salt)
	_, err = hash.Write(seed[:])
	if err != nil {
		return
	}

	I := hash.Sum(nil)

	if curve != CURVE_ED25519 {
		// TODO: Check invalid values
		// If curve is not ed25519 and IL is 0 or ≥ n (invalid key)
		// Set S := I and continue at step 2.
	}

	xkey, err = BuildExtendedKey(nil, I[:32], I[32:])

	return
}

func GenerateMasterKeyFromHexSeed(curve string, seed string) (*ExtendedKey, error) {
	key, err := hex.DecodeString(seed)
	if err != nil {
		return nil, err
	}

	return GenerateMasterKey(curve, key)
}

func GenerateMasterKeyFromMnemonic(curve string, fnParseMnemonic FnParseMnemonic, mnemonic string) (*ExtendedKey, error) {
	key, err := fnParseMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	return GenerateMasterKey(curve, key)
}
