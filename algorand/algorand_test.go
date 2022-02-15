package algorand

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testMnemonic = "mask wear topple mixture steel cupboard gain satoshi chuckle analyst spoil borrow melody punch start ivory resource olympic sibling conduct stairs manual curtain absorb citizen"
const testDerivedPub = "KI7RS3I3IHRR5T3WJ2YH2BILXF652Z2XNS45X6AXL4UVGZ76FKWKZYB4UM"
const testDerivedMnemonic = "finish when fun spatial art feed scare bomb fame hold measure hurt hill hope way warrior satisfy country inflict father option flee enlist abandon ice"

var testDerivedSeedBytes = []byte{184, 138, 62, 188, 8, 109, 134, 82, 9, 56, 25, 147, 34, 219, 19, 253, 214, 181, 181, 1, 95, 247, 250, 69, 204, 230, 58, 213, 77, 99, 89, 9}

func TestGenerateMnemonic(t *testing.T) {
	m, err := GenerateMnemonic()
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Validate with signing

	t.Logf("Mnemonic:\n%s", m)
}

func TestMnemonicToSeed(t *testing.T) {
	seed, err := MnemonicToSeed(testDerivedMnemonic)
	if err != nil {
		t.Fatal(t)
	}

	assert.Equal(t, testDerivedSeedBytes, seed)
}

func TestDeriveBasicAccount(t *testing.T) {
	for i := uint32(1); i <= 10; i++ {
		t.Run(fmt.Sprintf("account-%v", i), func(t *testing.T) {
			addr, mnemonic, err := DeriveBasicAccount(testMnemonic, i)
			if err != nil {
				t.Error(err)
			}

			if i == 10 {
				assert.Equal(t, testDerivedPub, addr)
				assert.Equal(t, testDerivedMnemonic, mnemonic)
			}

			t.Logf("Address: %s", addr)
			t.Logf("Mnemonic: %s", mnemonic)
		})
	}
}
