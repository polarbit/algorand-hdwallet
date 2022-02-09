package hdwallet

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

const (
	CURVE_ED25519   = "ed25519 seed"
	CURVE_SECP256K1 = "Bitcoin seed"
	CURVE_NIST256P1 = "Nist256p1 seed"
)

type HdWalletOptions struct {
	Curve           string
	FnParseMnemonic func(m string) ([]byte, error)
}

type ExtendedKey struct {
	Key       []byte
	ChainCode string
}

func (o *HdWalletOptions) GenerateMasterKey(seed []byte) ([]byte, []byte, error) {
	if len(seed) < 16 || len(seed) > 64 {
		return nil, nil, errors.New("Invalid  seed length")
	}

	if o.Curve != CURVE_ED25519 {
		return nil, nil, errors.New("Curve not supported")
	}

	salt := []byte(o.Curve)
	hash := hmac.New(sha512.New, salt)
	_, err := hash.Write(seed[:])
	if err != nil {
		panic(err)
	}

	I := hash.Sum(nil)
	IL := I[:32]
	IR := I[32:]

	if o.Curve != CURVE_ED25519 {
		// TODO: Check invalid values
		// If curve is not ed25519 and IL is 0 or ≥ n (invalid key)
		// Set S := I and continue at step 2.
	}

	return IL, IR, nil
}

func (o *HdWalletOptions) GenerateMasterKeyFromHexSeed(seed string) ([]byte, []byte, error) {
	key, err := hex.DecodeString(seed)
	if err != nil {
		return nil, nil, err
	}

	return o.GenerateMasterKey(key)
}

func (o *HdWalletOptions) GenerateMasterKeyFromMnemonic(mnemonic string) ([]byte, []byte, error) {
	key, err := o.FnParseMnemonic(mnemonic)
	if err != nil {
		return nil, nil, err
	}

	return o.GenerateMasterKey(key)
}

func CKD(key []byte, chainCode []byte, ix uint32) ([]byte, []byte) {
	vbuf := make([]byte, 4)
	binary.BigEndian.PutUint32(vbuf, ix)

	buf := make([]byte, 1)
	buf = append(buf, key...)
	buf = append(buf, vbuf...)

	hash := hmac.New(sha512.New, chainCode)
	_, err := hash.Write(buf)
	if err != nil {
		// TODO: Do not panic
		panic(err)
	}

	I := hash.Sum(nil)
	IL := I[:32]
	IR := I[32:]

	return IL, IR
}

func GetPublicKey(key []byte) string {
	priv := ed25519.NewKeyFromSeed(key)
	pub := priv.Public()

	buf := make([]byte, 32)
	copy(buf, pub.(ed25519.PublicKey))

	buf = append(make([]byte, 1), buf...)
	hexPub := hex.EncodeToString(buf)

	return hexPub
}
