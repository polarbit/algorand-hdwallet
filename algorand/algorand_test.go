package algorand

import (
	"fmt"
	"testing"
)

const testMnemonic = "mask wear topple mixture steel cupboard gain satoshi chuckle analyst spoil borrow melody punch start ivory resource olympic sibling conduct stairs manual curtain absorb citizen"

func TestGenerateMnemonic(t *testing.T) {
	m, err := GenerateMnemonic()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Mnemonic:\n%s", m)
}

func TestDeriveAccount(t *testing.T) {
	for i := uint32(1); i <= 10; i++ {
		t.Run(fmt.Sprintf("account-%v", i), func(t *testing.T) {
			addr, mnemonic, err := DeriveAccount(testMnemonic, i)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Address: %s", addr)
			t.Logf("Mnemonic: %s", mnemonic)
		})
	}
}
