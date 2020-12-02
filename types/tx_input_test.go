// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import (
	"bytes"
	"encoding/hex"
	"github.com/reed/common/byteutil/byteconv"
	"github.com/reed/crypto"
	"github.com/reed/vm/vmcommon"
	"github.com/reed/wallet"
	"testing"
)

var (
	inpId, _     = hex.DecodeString("f5cb9eb86f0ee72f00f88c769ca2cb9b635072edb2427ccf1bcf751d788c43ab")
	spoutId, _   = hex.DecodeString("8298e85d0e8465187310a24df198cd844e4cec6e1c146a1e23639def8c91bfe8")
	spsrcId, _   = hex.DecodeString("c7a30af288b381dde540d0edca910a846f3ea4bf1a45b3ce2e0ba0540d0459f3")
	scriptPk, _  = hex.DecodeString("b4af5bb2bb1fb9086d2cac65a667f1810dfb0ddd904f2edc947227271fdcaba5b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	scriptSig, _ = hex.DecodeString("3fe4fed59e6201b3537748514e1b89f03d5a7b5140cb723e8348fc9ef90fd442290b91b0754c32063c5bd725205aec870e9989a9eda5fb7960eff7fd82a4d50f")
)

func TestTxInput_GenerateID(t *testing.T) {
	input := mockTxInput()
	id := input.GenerateID()

	split := []byte(":")

	b := bytes.Join([][]byte{
		spoutId,
		split,
		spsrcId,
		split,
		byteconv.Uint64ToByte(0),
		split,
		byteconv.Uint64ToByte(10),
		split,
		scriptPk,
	}, []byte{})
	h := BytesToHash(crypto.Sha256(b))
	if !bytes.Equal(id.Bytes(), h.Bytes()) {
		t.Error("not equal")
	}

}

func TestTxInput_GenerateScriptSig(t *testing.T) {
	input := mockTxInput()

	pub, _ := hex.DecodeString("b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	priv, _ := hex.DecodeString("b4af5bb2bb1fb9086d2cac65a667f1810dfb0ddd904f2edc947227271fdcaba5b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")

	wt := &wallet.Wallet{
		Pub:  pub,
		Priv: priv,
	}

	txId, _ := hex.DecodeString("9bd5b8b068be0c6ed5cc5e5e160c7a369b869563c0b7fbf58224af19bd425a54")

	txh := BytesToHash(txId)
	scriptSig, err := input.GenerateScriptSig(wt, &txh)
	if err != nil {
		t.Error(err)
	}

	//scriptSig = signature(64) + public key(32) = 96

	if len(*scriptSig) != 96+2 {
		t.Error("invalid scriptSig len")
	}

	signature := (*scriptSig)[1:65]

	if !bytes.Equal((*scriptSig)[(64+2):], wt.Pub) {
		t.Error("scriptSig error:last len(32) not equal public key")
	}

	message := bytes.Join([][]byte{
		inpId,
		txId,
	}, []byte{})

	if !crypto.Verify(wt.Pub, message, signature) {
		t.Error("signature not pass")
	}

	expectSig := crypto.Sign(wt.Priv, message)

	expectScriptSig := bytes.Join([][]byte{
		{byte(vmcommon.OpPushData64)},
		expectSig,
		{byte(vmcommon.OpPushData32)},
		wt.Pub,
	}, []byte{})

	if !bytes.Equal(*scriptSig, expectScriptSig) {
		t.Error("not equal")
	}

}

func mockTxInput() TxInput {
	spend := Spend{
		SpendOutputId: BytesToHash(spoutId),
		SoureId:       BytesToHash(spsrcId),
		SourcePos:     0,
		Amount:        10,
		ScriptPk:      scriptPk,
	}

	return TxInput{
		ID:        BytesToHash(inpId),
		Spend:     spend,
		ScriptSig: scriptSig,
	}
}
