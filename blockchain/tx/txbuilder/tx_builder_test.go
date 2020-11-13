package txbuilder

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tybc/core/types"
	"github.com/tybc/crypto"
	"github.com/tybc/vm/vmcommon"
	"github.com/tybc/wallet"
	"testing"
)

func TestMapTx(t *testing.T) {
	req := mockReqTx()
	tx, err := mapTx(req)
	if err != nil {
		t.Error("map tx error")
	}
	marshal, err := json.Marshal(tx)
	if err != nil {
		t.Error("marshal error")
	}

	fmt.Println(string(marshal))
}

func TestSetScriptSig(t *testing.T) {

}

func mockUTXO() *types.UTXO {
	var1, _ := hex.DecodeString("5276236b8ff88b9075d87ccbde5da44555f0111cd3789df6c9404f3e2320567a")
	var2, _ := hex.DecodeString("e6f48722821ee92c838577fabeb0adbbad2851f6e63f5ce5cc5a0758d0713406")
	var3, _ := hex.DecodeString("d7eef5283d692d22b1951864a6329f4600b4b420c545ad857e704bacd7c258e2")

	var utxo = &types.UTXO{
		OutputId:  var1,
		SoureId:   var2,
		SourcePos: 0,
		Amount:    19,
		Address:   var3,
		ScriptPk:  vmcommon.BuildP2PKHScript(crypto.Sha256(var3)),
	}

	return utxo
}

func mockReqTx() *types.SubmitTxRequest {
	var reqInps []types.ReqInput
	var reqOnps []types.ReqOutput
	reqInp := &types.ReqInput{
		SpendOutputId: "5fa60a785896d7e9ba0141dfdbe596a01779f28a320ed9f6799918379f97e3f0",
	}
	reqInps = append(reqInps, *reqInp)

	reqOnp := &types.ReqOutput{
		Address: "774471dc03273212f7e4beb45893d2d7dc315d4daf45bbbe445c8c4c5ecf4ba5",
		Amount:  199,
	}
	reqOnps = append(reqOnps, *reqOnp)

	req := &types.SubmitTxRequest{
		Password:  "123",
		TxInputs:  reqInps,
		TxOutputs: reqOnps,
	}

	return req
}

func mockWallet() *wallet.Wallet {
	pub, _ := hex.DecodeString("b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	priv, _ := hex.DecodeString("b4af5bb2bb1fb9086d2cac65a667f1810dfb0ddd904f2edc947227271fdcaba5b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")

	return &wallet.Wallet{
		Pub:  pub,
		Priv: priv,
	}
}
