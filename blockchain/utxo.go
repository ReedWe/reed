package blockchain

import (
	"github.com/tybc/errors"
	"github.com/tybc/types"
)

func GetUtxoByOutputId(store *Store, outputId types.Hash) (*types.UTXO, error) {
	id := outputId.Bytes()
	utxo, err := (*store).GetUtxo(id)
	if utxo == nil || err != nil {
		return nil, errors.Wrapf(err, "utxo(outputId:%x) does not exists", id)
	}
	return utxo, nil
}
