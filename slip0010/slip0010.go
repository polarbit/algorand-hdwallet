package slip0010

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"errors"

	"github.com/polarbit/hdwallet/bip32path"
	"github.com/polarbit/hdwallet/utils"
	"github.com/tyler-smith/go-bip39"
)

const (
	CURVE_ED25519   string = "ed25519 seed"
	CURVE_SECP256K1        = "Bitcoin seed"
	CURVE_NIST256P1        = "Nist256p1 seed"
)

type ExtendedKey struct {
	PrivateKey      []byte
	ChainCode       []byte
	CurvePrivateKey []byte
	CurvePublicKey  []byte
	PublicKey       []byte
	Fingerprint     []byte
	ParentKey       *ExtendedKey
}

func DeriveAccount(curve string, path *bip32path.Bip32Path, xmaster *ExtendedKey) (xchd *ExtendedKey, err error) {
	if curve != CURVE_ED25519 {
		return nil, errors.New("Only ed25519 is supported")
	}

	xchd = xmaster
	for _, s := range path.Segments {
		xchd, err = CKD(curve, xchd, s.Value)
	}

	return
}

func CKD(curve string, xpar *ExtendedKey, i uint32) (xchd *ExtendedKey, err error) {
	if len(xpar.PrivateKey) != 32 {
		panic("Invalid xparent.Key")
	}
	if len(xpar.ChainCode) != 32 {
		panic("Ivalid xparent.ChainCode")
	}
	if i < bip32path.HARDENED_OFFSET {
		return nil, errors.New("Only hardened keys are supported")
	}

	vbuf := make([]byte, 4)
	binary.BigEndian.PutUint32(vbuf, i)

	key := xpar.PrivateKey
	chainCode := xpar.ChainCode

	buf := make([]byte, 1)
	buf = append(buf, key...)
	buf = append(buf, vbuf...)

	hash := hmac.New(sha512.New, chainCode)
	_, err = hash.Write(buf)
	if err != nil {
		return
	}

	I := hash.Sum(nil)
	xchd, err = ExtendKey(xpar, I[:32], I[32:])

	return
}

func ExtendKey(xpar *ExtendedKey, key []byte, chainCode []byte) (*ExtendedKey, error) {
	pub, priv, err := GetCurveKeyPair_Ed25519(key)
	if err != nil {
		return nil, err
	}

	xpub := append([]byte{0x00}, pub...) // Padding is reqired for ED25519 acc. to SLIP-0010
	xkey := &ExtendedKey{
		PrivateKey:      key,
		ChainCode:       chainCode,
		CurvePrivateKey: priv, // !!! Actually curve private key's first half [:32] and extended private key are equal for ED25519
		CurvePublicKey:  pub,
		PublicKey:       xpub,
		ParentKey:       xpar,
	}

	if xpar != nil {
		keyId, err := utils.Hash160(xpar.PublicKey)
		if err != nil {
			return nil, err
		}
		xkey.Fingerprint = keyId[:4]
	} else {
		xkey.Fingerprint = []byte{0, 0, 0, 0}
	}

	return xkey, nil
}

func GetCurveKeyPair_Ed25519(key []byte) ([]byte, []byte, error) {
	priv := ed25519.NewKeyFromSeed(key)
	pub := priv.Public()

	buf := make([]byte, 32)
	copy(buf, pub.(ed25519.PublicKey))

	return buf, priv, nil
}

func GenerateMasterKey(curve string, seed []byte) (xmaster *ExtendedKey, err error) {
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

	xmaster, err = ExtendKey(nil, I[:32], I[32:])

	return
}

func MnemonicToSeed(mnemonic string) ([]byte, error) {
	return bip39.NewSeedWithErrorChecking(mnemonic, "")
}
