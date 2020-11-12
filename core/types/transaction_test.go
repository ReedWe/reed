package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/tybc/crypto"
	"testing"
)

func TestTx_GenerateID(t *testing.T) {
	id1, _ := hex.DecodeString("416c7cfb8a0836d51517b6ae32e0dee579a554f41c10448da847f0888b76881e")
	id2, _ := hex.DecodeString("1a946d5a05761732ae162a886a767fdefa3ce3ee66f9cc0352481e0ae751db9c")
	id3, _ := hex.DecodeString("f3e023bebab7e9b6f83a94ac4064f0f4e2ea5f67af77f6857a6c0fe9121d359f")

	inp1 := &TxInput{
		ID: BytesToHash(id1),
	}
	inp2 := &TxInput{
		ID: BytesToHash(id2),
	}
	inp3 := &TxInput{
		ID: BytesToHash(id3),
	}

	var inps []TxInput
	inps = append(inps, *inp1, *inp2, *inp3)

	fmt.Printf("len %d \n", len(inps))

	emptyTX := &Tx{

	}
	if _, err := emptyTX.GenerateID(); err == nil {
		t.Error("empty inputs would not pass")
	}

	tx := &Tx{
		TxInput: inps,
	}

	txId, err := tx.GenerateID()
	if err != nil {
		t.Error(err)
	}

	b := bytes.Join([][]byte{
		id1, id2, id3,
	}, []byte{})

	if !bytes.Equal((*txId).Bytes(), crypto.Sha256(b)) {
		t.Error("GenerateID error")
	}

}

func mockJustIDTxInput() TxInput {
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
