package hdwallet

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"errors"

	"github.com/polarbit/algorand-hdwallet/bip32path"
	"github.com/polarbit/algorand-hdwallet/utils"
	"github.com/tyler-smith/go-bip39"
)

const (
	CURVE_ED25519   = "ed25519 seed"
	CURVE_SECP256K1 = "Bitcoin seed"
	CURVE_NIST256P1 = "Nist256p1 seed"
)

type ExtendedKey struct {
	Key         []byte
	ChainCode   []byte
	PrivateKey  []byte
	PublicKey   []byte
	Fingerprint []byte
	ParentKey   *ExtendedKey
}

func GenerateAccount(path *bip32path.Bip32Path, masterXKey *ExtendedKey) (xkey *ExtendedKey, err error) {
	xkey = masterXKey

	for _, s := range path.Segments {
		xkey, err = CKD(xkey, s.Value)
	}

	return
}

func BuildExtendedKey(xparent *ExtendedKey, key []byte, chainCode []byte) (*ExtendedKey, error) {
	priv, pub, err := GetCurveKeyPair_Ed25519(key)
	if err != nil {
		return nil, err
	}

	xkey := &ExtendedKey{
		Key:        key,
		ChainCode:  chainCode,
		PrivateKey: priv,
		PublicKey:  pub,
		ParentKey:  xparent,
	}

	if xparent != nil {
		keyId, err := utils.Hash160(xparent.PublicKey)
		if err != nil {
			return nil, err
		}
		xkey.Fingerprint = keyId[:4]
	} else {
		xkey.Fingerprint = []byte{0, 0, 0, 0}
	}

	return xkey, nil
}

func CKD(xparent *ExtendedKey, ix uint32) (xkey *ExtendedKey, err error) {
	vbuf := make([]byte, 4)
	binary.BigEndian.PutUint32(vbuf, ix)

	key := xparent.Key
	chainCode := xparent.ChainCode

	buf := make([]byte, 1)
	buf = append(buf, key...)
	buf = append(buf, vbuf...)

	hash := hmac.New(sha512.New, chainCode)
	_, err = hash.Write(buf)
	if err != nil {
		return
	}

	I := hash.Sum(nil)
	xkey, err = BuildExtendedKey(xparent, I[:32], I[32:])

	return
}

func GetCurveKeyPair_Ed25519(key []byte) ([]byte, []byte, error) {
	priv := ed25519.NewKeyFromSeed(key)
	pub := priv.Public()

	buf := make([]byte, 32)
	copy(buf, pub.(ed25519.PublicKey))

	pubKey := append(make([]byte, 1), buf...)

	return priv, pubKey, nil
}

func GenerateMasterKey(curve string, seed []byte) (xkey *ExtendedKey, err error) {
	if len(seed) < 16 || len(seed) > 64 {
		return nil, errors.New("Invalid  seed length")
	}

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
		// TODO: Check invalid values for other curves (SECP256K1)
		// If curve is not ed25519 and IL is 0 or ≥ n (invalid key)
		// Set S := I and continue at step 2.
		// Ref: BIP32
	}

	xkey, err = BuildExtendedKey(nil, I[:32], I[32:])

	return
}

func MnemonicToSeed(mnemonic string) ([]byte, error) {
	return bip39.NewSeedWithErrorChecking(mnemonic, "")
}
