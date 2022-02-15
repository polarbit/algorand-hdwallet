package hdwallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testMnemonic string = "segment inhale symptom olive cheese tissue vacuum lazy sketch salt enroll wink oyster hen glory food weasel comic glow legal cute diet fun real"
)

func TestMnemonicEquality(t *testing.T) {
	path := "m/44'/283'/0'"
	_, x, err := DeriveAddressFromMnemonic(testMnemonic, path)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, x.PrivateKey, x.CurvePrivateKey[:32])
}

func TestDeriveAddress(t *testing.T) {
	path := "m/44'/283'/0'"
	a, _, err := DeriveAddressFromMnemonic(testMnemonic, path)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Address: %v", a)
}
