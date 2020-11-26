package vm

import (
	"bytes"
	"encoding/hex"
	"github.com/reed/crypto"
	"github.com/reed/types"
	"github.com/reed/wallet"
	"testing"
)

var (
	pub, _  = hex.DecodeString("b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	priv, _ = hex.DecodeString("b4af5bb2bb1fb9086d2cac65a667f1810dfb0ddd904f2edc947227271fdcaba5b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	addr, _ = hex.DecodeString("b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")

	inpId, _ = hex.DecodeString("cba11e58e81e87bb04b1d8406898a54f071f22f36fb6a49adb75f0ee95ee8af7")
	txId, _  = hex.DecodeString("77ae8f0351299d59bffdf7d69160e78069e0b9bbe5bc9c2403169cd2ec37d542")
)

func TestVM_Run(t *testing.T) {

	signFunc := func(pk []byte, sig []byte) bool {
		message := bytes.Join([][]byte{
			inpId,
			txId,
		}, []byte{})

		return crypto.Verify(pk, message, sig)
	}

	vm := NewVirtualMachine(mockScriptData(), mockScriptPKData(), signFunc)

	vm.Run()
}

func mockScriptData() []byte {
	wt := &wallet.Wallet{
		Pub:  pub,
		Priv: priv,
	}

	txIdHash := types.BytesToHash(txId)

	inp := &types.TxInput{
		ID: types.BytesToHash(inpId),
	}
	sig, _ := inp.GenerateScriptSig(wt, &txIdHash)
	return *sig
}

func mockScriptPKData() []byte {

	iop := &types.TxOutput{
		Address: addr,
	}
	return iop.GenerateLockingScript()
}
