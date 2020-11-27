package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/reed/crypto"
	"strconv"
	"testing"
)

func TestTx_GenerateID(t *testing.T) {
	id1, _ := hex.DecodeString("416c7cfb8a0836d51517b6ae32e0dee579a554f41c10448da847f0888b76881e")
	id2, _ := hex.DecodeString("1a946d5a05761732ae162a886a767fdefa3ce3ee66f9cc0352481e0ae751db9c")
	id3, _ := hex.DecodeString("f3e023bebab7e9b6f83a94ac4064f0f4e2ea5f67af77f6857a6c0fe9121d359f")
	inp1 := &TxInput{
		ID:        BytesToHash(id1),
		ScriptSig: []byte("miner"),
	}

	inp2 := &TxInput{
		ID: BytesToHash(id2),
	}
	inp3 := &TxInput{
		ID: BytesToHash(id3),
	}

	var inps []TxInput
	inps = append(inps, *inp1, *inp2, *inp3)

	oid1, _ := hex.DecodeString("4116ac20904e2b01daf6ececf6f0afb960b760bdfce179f3afc3f6c5eea36c82")
	oid2, _ := hex.DecodeString("570c9dcf6705794eee3d50025d4f838ef15b87992962247252f23d8a106d9758")
	oid3, _ := hex.DecodeString("4c66effcf3d375e28bce54d5ae405d47c721fdd7d135d0ae14e487b7bf58f324")
	onp1 := &TxOutput{
		ID: BytesToHash(oid1),
	}
	onp2 := &TxOutput{
		ID: BytesToHash(oid2),
	}
	onp3 := &TxOutput{
		ID: BytesToHash(oid3),
	}

	var onps []TxOutput
	onps = append(onps, *onp1, *onp2, *onp3)

	emptyTX := &Tx{

	}
	if _, err := emptyTX.GenerateID(); err == nil {
		t.Error("empty inputs would not pass")
	}

	tx := &Tx{
		TxInput:  inps,
		TxOutput: onps,
	}

	txId, err := tx.GenerateID()
	if err != nil {
		t.Error(err)
	}

	b := bytes.Join([][]byte{
		id1, id2, id3, oid1, oid2, oid3,
	}, []byte{})

	if !bytes.Equal((*txId).Bytes(), crypto.Sha256(b)) {
		t.Error("GenerateID error")
	}

	var txs []Tx
	txs = append(txs, *tx)

	block := &Block{
		Transactions: &txs,
	}
	incrementExtraNonce(19, block)

	marshal, _ := json.Marshal(tx.TxInput[0])
	fmt.Printf("%s", marshal)
}

func incrementExtraNonce(extraNonce uint64, cblock *Block) {

	txs := *cblock.Transactions

	msg := bytes.Join([][]byte{txs[0].TxInput[0].ScriptSig, []byte(strconv.FormatUint(extraNonce,10))}, []byte{})

	txs[0].TxInput[0].ScriptSig = msg
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
