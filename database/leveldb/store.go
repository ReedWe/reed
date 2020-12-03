package leveldb

import (
	"encoding/json"
	"github.com/reed/errors"
	"github.com/reed/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	storeTxErr    = errors.New("transaction leveldb store error")
	storeUtxoErr  = errors.New("utxo leveldb store error")
	storeBlockErr = errors.New("block leveldb store error")
)

const (
	txpoolPrefix = "TP:"
	utxoPrefix   = "U:"
	blockPrefix  = "B:"
	highestBlock = "HIGHESTBLOCK"
)

type Store struct {
	db dbm.DB
}

func NewStore(db dbm.DB) *Store {
	return &Store{
		db: db,
	}
}
//
//func (store *Store) SaveTx(tx *types.Tx) error {
//	b, err := json.Marshal(tx)
//	if err != nil {
//		return errors.Wrapf(storeTxErr, "SaveTx json marshal error (id=%x)", tx.ID)
//	}
//
//	store.db.Set(getKey(txPrefix, tx.ID.Bytes()), b)
//	return nil
//}
//
//func (store *Store) GetTx(id []byte) (*types.Tx, error) {
//	b := store.db.Get(getKey(txPrefix, id))
//	if b == nil {
//		return nil, nil
//	}
//	tx := &types.Tx{}
//
//	if err := json.Unmarshal(b, tx); err != nil {
//		return nil, errors.Wrapf(storeTxErr, "transaction unmarshal failed (id=%x)", id)
//	}
//	return tx, nil
//}

func (store *Store) GetUtxo(id []byte) (*types.UTXO, error) {
	data := store.db.Get(getKey(utxoPrefix, id))
	if data == nil {
		return nil, errors.Wrapf(storeUtxoErr, "utxo does not exists (id=%x)", id)
	}

	utxo := &types.UTXO{}
	if err := json.Unmarshal(data, &utxo); err != nil {
		return nil, err
	}

	return utxo, nil
}

func (store *Store) SaveUtxos(expiredUtxoIds []*types.Hash, utxos []*types.UTXO) error {
	batch := store.db.NewBatch()
	for _, e := range expiredUtxoIds {
		batch.Delete(e.Bytes())
	}

	for _, utxo := range utxos {
		b, err := json.Marshal(utxo)
		if err != nil {
			return errors.Wrapf(storeUtxoErr, "SaveUtxos json marshal error (ID=%x)", utxo.ID)
		}
		batch.Set(getKey(utxoPrefix, utxo.ID.Bytes()), b)
	}
	batch.Write()
	return nil
}

func (store *Store) GetHighestBlock() (*types.Block, error) {
	data := store.db.Get([]byte(highestBlock))
	if data == nil {
		return nil, nil
	}

	block := &types.Block{}
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, errors.Wrapf(err, "GetHighestBlock json umarshal error")
	}
	return block, nil
}

func (store *Store) GetBlock(hash []byte) (*types.Block, error) {
	data := store.db.Get(getKey(blockPrefix, hash))
	if data == nil {
		return nil, errors.Wrapf(storeBlockErr, "Block does not exists (hash=%x)", hash)
	}
	block := &types.Block{}
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, errors.Wrapf(err, "GetBlock json umarshal error (hash=%x)", hash)
	}
	return block, nil
}

func (store *Store) SaveBlock(block *types.Block) error {
	b, err := json.Marshal(block)
	if err != nil {
		return errors.Wrapf(storeBlockErr, "SaveBlock json marshal error (hash=%x)", block.GetHash())
	}

	batch := store.db.NewBatch()
	batch.Set(getKey(blockPrefix, block.GetHash().Bytes()), b)
	batch.Set([]byte(highestBlock), b)
	batch.Write()
	return nil
}

func getKey(prefix string, id []byte) []byte {
	return []byte(prefix + string(id))
}
