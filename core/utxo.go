package core

import (
	"github.com/tybc/blockchain"
	"github.com/tybc/core/types"
	"github.com/tybc/errors"
)


func GetUtxoByOutputId(store *blockchain.Store, outputId types.Hash) (*types.UTXO, error) {
	id := outputId[:]
	utxo, err := (*store).GetUtxo(id)
	if utxo == nil || err != nil {
		return nil, errors.Wrapf(err, "utxo(outputId:%x) does not exists", id)
	}
	return utxo, nil
}
