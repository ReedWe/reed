package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/tybc/crypto"
	"github.com/tybc/vm"
	"testing"
)

func TestTxOutput_GenerateID(t *testing.T) {
	var addr, _ = hex.DecodeString("c27f26c8bf818e5509abacfc20206d43fc0db6a415f20d48726eb8cd2888f68e")
	var scriptPK, _ = hex.DecodeString("bf8776efb3367228d115c325a623b3fe6b359a87e45d25c98a506e203b0ec5b1fc0db6a415f20d48726eb8cd2888f68")
	amt := uint64(199)
	icb := false
	output := &TxOutput{
		IsCoinBase: icb,
		Amount:     amt,
		Address:    addr,
		ScriptPk:   scriptPK,
	}

	id := output.GenerateID()

	var amtByte = make([]byte, 8)
	binary.LittleEndian.PutUint64(amtByte, amt)

	var datas [][]byte
	split := []byte(":")
	datas = append(datas, []byte{0}, split, addr, split, amtByte, split, scriptPK)

	h := crypto.Sha256(datas...)

	if !bytes.Equal(id.Bytes(), h) {
		t.Fatalf("GenerateID error")
	}

}

func TestTxOutput_SetLockingScript(t *testing.T) {
	var pub, _ = hex.DecodeString("b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	o := TxOutput{
		Address: pub,
	}

	script := o.GenerateLockingScript()

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
