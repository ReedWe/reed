package types

import (
	"bytes"
	"encoding/hex"
	"github.com/tybc/crypto"
	"github.com/tybc/vm"
	"testing"
)

func TestOutputLockingScript(t *testing.T) {
	var pub, _ = hex.DecodeString("b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	o := TxOutput{
		Address: pub,
	}

	o.SetLockingScript()
	script := o.ScriptPk

	if len(script) != (32 + 4) {
		t.Errorf("script len error;actual len=%d", len(script))
	}
	if script[0] != vm.OpDup {
		t.Error("script first part is not OP_DUP")
	}
	if script[1] != vm.OpHash256 {
		t.Error("script second part is not OP_HASH256")
	}

	pubHash := crypto.Sha256(pub)

	if !bytes.Equal(script[2:34], pubHash) {
		t.Error("script third part is not Hash data")
	}

	if script[34] != vm.OpEqualverify {
		t.Error("script fourth part is not OP_EQUALVERIFY")
	}
	if script[35] != vm.OpChecksig {
		t.Error("script fifth part is not OP_CHECKSIG")
	}

}
