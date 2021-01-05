// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package leveldb

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/reed/types"
	dbm "github.com/tendermint/tmlibs/db"
	"os"
	"testing"
)

func TestSaveTx(t *testing.T) {
	inputs := make([]*types.TxInput, 0)

	inpId, _ := hex.DecodeString("f5cb9eb86f0ee72f00f88c769ca2cb9b635072edb2427ccf1bcf751d788c43ab")
	spoutId, _ := hex.DecodeString("8298e85d0e8465187310a24df198cd844e4cec6e1c146a1e23639def8c91bfe8")
	spsrcId, _ := hex.DecodeString("c7a30af288b381dde540d0edca910a846f3ea4bf1a45b3ce2e0ba0540d0459f3")
	scriptPk, _ := hex.DecodeString("b4af5bb2bb1fb9086d2cac65a667f1810dfb0ddd904f2edc947227271fdcaba5b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	scriptSig, _ := hex.DecodeString("3fe4fed59e6201b3537748514e1b89f03d5a7b5140cb723e8348fc9ef90fd442290b91b0754c32063c5bd725205aec870e9989a9eda5fb7960eff7fd82a4d50f")

	spend := types.Spend{
		SpendOutputId: types.BytesToHash(spoutId),
		SoureId:       types.BytesToHash(spsrcId),
		SourcePos:     0,
		Amount:        10,
		ScriptPk:      scriptPk,
	}

	input := &types.TxInput{
		ID:        types.BytesToHash(inpId),
		Spend:     spend,
		ScriptSig: scriptSig,
	}

	inputs = append(inputs, input)

	//=============== OUT PUT ===============

	outputs := make([]*types.TxOutput, 0)

	outId, _ := hex.DecodeString("4eaa97d3fb3e5659bc4f9c805a1ae747c69a39b57b77cb5b87addb060abdd623")
	outAddr, _ := hex.DecodeString("d75999e54ad60ac7d01c1f7d1fc6339daf43907ba7cd817a0cf28c0a36e68acc")
	outScriptPK, _ := hex.DecodeString("b4af5bb2bb1fb9086d2cac65a667f1810dfb0ddd904f2edc947227271fdcaba5b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	output := &types.TxOutput{
		ID:         types.BytesToHash(outId),
		IsCoinBase: false,
		Amount:     15,
		Address:    outAddr,
		ScriptPk:   outScriptPK,
	}
	outputs = append(outputs, output)

	tx := &types.Tx{
		TxInput:  inputs,
		TxOutput: outputs,
	}
	marshal, err := json.Marshal(tx)
	if err != nil {
		t.Error("marshal error")
	}
	store := dbm.NewDB("core", dbm.LevelDBBackend, os.Getenv("GOPATH")+"/src/github.com/reed/database/file/")

	store.Set([]byte("tx123"), marshal)

	b := store.Get([]byte("tx123"))

	tx2 := &types.Tx{}
	err = json.Unmarshal(b, tx2)
	if err != nil {
		t.Errorf("Unmarshal error %s", err.Error())
	}

	fmt.Printf("%x", tx2.TxInput[0].ScriptPk)

}
