package txbuilder

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMapTx(t *testing.T) {
	req := mockReqTx()
	tx, err := req.MapTx()
	if err != nil {
		t.Error("map tx error")
	}
	marshal, err := json.Marshal(tx)
	if err != nil {
		t.Error("marshal error")
	}

	fmt.Println(string(marshal))
}

func mockReqTx() *SubmitTxRequest {
	var reqInps []ReqInput
	var reqOnps []ReqOutput
	reqInp := &ReqInput{
		SpendOutputId: "5fa60a785896d7e9ba0141dfdbe596a01779f28a320ed9f6799918379f97e3f0",
	}
	reqInps = append(reqInps, *reqInp)

	reqOnp := &ReqOutput{
		Address: "774471dc03273212f7e4beb45893d2d7dc315d4daf45bbbe445c8c4c5ecf4ba5",
		Amount:  199,
	}
	reqOnps = append(reqOnps, *reqOnp)

	req := &SubmitTxRequest{
		Password:  "123",
		TxInputs:  reqInps,
		TxOutputs: reqOnps,
	}

	return req
}
