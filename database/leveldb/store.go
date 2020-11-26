package leveldb

import (
	"encoding/json"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/reed/errors"
	"github.com/reed/types"
)

var (
	txPrefix   = "TX:"
	utxoPrefix = "UTXO:"
	storeErr   = errors.New("leveldb error")
)

type Store struct {
	db dbm.DB
}

func NewStore(db dbm.DB) *Store {
	return &Store{
		db: db,
	}
}

func getTxKey(id []byte) []byte {
	return []byte(txPrefix + string(id))
}

func getUtxoKey(id []byte) []byte {
	return []byte(utxoPrefix + string(id))
}

func (store *Store) AddTx(tx *types.Tx) error {
	b, err := json.Marshal(tx)
	if err != nil {
		return errors.Wrapf(err, "AddTx json marshal error")
	}

	store.db.Set(getTxKey(tx.ID.Bytes()), b)
	return nil
}

func (store *Store) GetTx(id []byte) (*types.Tx, error) {
	b := store.db.Get(getTxKey(id))
	if b == nil {
		return nil, nil
	}
	tx := &types.Tx{}

	if err := json.Unmarshal(b, tx); err != nil {
		return nil, errors.Wrapf(storeErr, "tx(id=%x) unmarshal failed", id)
	}
	return tx, nil
}

func (store *Store) GetUtxo(id []byte) (*types.UTXO, error) {
	data := store.db.Get(getUtxoKey(id))
	if data == nil {
		return nil, errors.Wrapf(storeErr, "utxo(id=%x) does not exists", id)
	}

	utxo := &types.UTXO{}
	if err := json.Unmarshal(data, &utxo); err != nil {
		return nil, err
	}

	return utxo, nil
}

func (store *Store) SaveUtxos(expiredUtxoIds []types.Hash, utxos *[]types.UTXO) error {

	batch := store.db.NewBatch()
	for _, e := range expiredUtxoIds {
		batch.Delete(e.Bytes())
	}

	for _, utxo := range *utxos {
		b, err := json.Marshal(utxo)
		if err != nil {
			return errors.Wrapf(err, "SaveUtxos json marshal error")
		}
		batch.Set(getUtxoKey(utxo.OutputId.Bytes()), b)
	}
	batch.Write()
	return nil
}
