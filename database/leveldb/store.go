package leveldb

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tybc/errors"
	"github.com/tybc/types"
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

func getTxKey(id *[]byte) []byte {
	return []byte(txPrefix + string(*id))
}

func getUtxoKey(id *[]byte) []byte {
	return []byte(utxoPrefix + string(*id))
}

func (store *Store) GetTx(id []byte) (*types.Tx, error) {
	b := store.db.Get(getTxKey(&id))
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
	var utxo types.UTXO
	data := store.db.Get(getUtxoKey(&id))
	if data == nil {
		return nil, errors.Wrapf(storeErr, "utxo(id=%x) does not exists", id)
	}

	if err := proto.Unmarshal(data, &utxo); err != nil {
		return nil, err
	}

	return &utxo, nil
}

func (store *Store) SaveUtxo(utxo *types.TxOutput) {
	//b, err := proto.Marshal(utxo)
	//fmt.Printf("data %x \n", b)
	//
	//if err != nil {
	//	return
	//}
	//
	//k := getUtxoKey(&utxo.Id)
	//fmt.Printf("Key %x \n", k)
	//
	//store.db.Set(getUtxoKey(&utxo.Id), b)
}
