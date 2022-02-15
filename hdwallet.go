package hdwallet

import (
	"github.com/polarbit/hdwallet/algorand"
	"github.com/polarbit/hdwallet/bip32path"
	"github.com/polarbit/hdwallet/slip0010"
	"github.com/tyler-smith/go-bip39"
)

func DeriveAddressFromMnemonic(mnemonic string, derivationPath string) (a string, xkey *slip0010.ExtendedKey, err error) {
	path, err := bip32path.Parse(derivationPath)
	if err != nil {
		return
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return
	}

	xmaster, err := slip0010.GenerateMasterKey(slip0010.CURVE_ED25519, seed)
	if err != nil {
		return
	}

	xkey, err = slip0010.DeriveAccount(slip0010.CURVE_ED25519, path, xmaster)
	if err != nil {
		return
	}

	a, err = algorand.PublicKeyToAddress(xkey.CurvePublicKey)
	if err != nil {
		return
	}

	return
}
