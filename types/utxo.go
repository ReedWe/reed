// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

type UTXO struct {
	OutputId  Hash   `json:"outputId"`
	SoureId   Hash   `json:"sourceId"`
	SourcePos uint64 `json:"sourcePos"`
	Amount    uint64 `json:"amount"`
	Address   []byte `json:"address"`
	ScriptPk  []byte `json:"scriptPK"`
}
