package hdwallet

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ()

func TestFromMnemonicToSeed(t *testing.T) {
	mnemonic := "craft curve like tool damp voice jaguar sick fit immense pistol able omit define produce"
	expectedSeed := "c2b7c55f5a47068fced1ceb14b88c6fbf2e60f4546a3a8345b71f55d872574db46cba53f8a80d6e792c8dfcfd74c81ae51dbfeb4cab283263f7c98a4982a7f79"

	seed, err := MnemonicToSeed(mnemonic)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, expectedSeed, hex.EncodeToString(seed))
}
