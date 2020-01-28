package keys

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoxSeal(t *testing.T) {
	alice := GenerateX25519Key()
	bob := GenerateX25519Key()

	msg := "Hey bob, it's alice. The passcode is 12345."
	encrypted := BoxSeal([]byte(msg), bob.PublicKey(), alice)

	out, err := BoxOpen(encrypted, alice.PublicKey(), bob)
	require.NoError(t, err)
	require.Equal(t, "Hey bob, it's alice. The passcode is 12345.", string(out))
}

func ExampleBoxSeal() {
	ak := GenerateX25519Key()
	bk := GenerateX25519Key()

	msg := "Hey bob, it's alice. The passcode is 12345."
	encrypted := BoxSeal([]byte(msg), bk.PublicKey(), ak)

	out, err := BoxOpen(encrypted, ak.PublicKey(), bk)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(out))
	// Output:
	// Hey bob, it's alice. The passcode is 12345.
}
