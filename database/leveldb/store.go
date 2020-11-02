package leveldb

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tybc/core/types"
)

var (
	utxoPrefix = "UTXO:"
)

type Store struct {
	db dbm.DB
}

func NewStore(db dbm.DB) *Store {
	return &Store{
		db: db,
	}
}

func getKey(id *[]byte) []byte {
	return []byte(utxoPrefix + string(*id))
}

func (store *Store) GetUtxo(id []byte) (*types.TxOutput, error) {
	var utxo types.TxOutput
	data := store.db.Get(getKey(&id))
	if data == nil {
		return nil, errors.New("utxo does not exists")
	}

	if err := proto.Unmarshal(data, &utxo); err != nil {
		return nil, err
	}

	return &utxo, nil
}

func (store *Store) SaveUtxo(utxo *types.TxOutput) {
	b, err := proto.Marshal(utxo)
	fmt.Printf("data %x \n", b)

	if err != nil {
		return
	}

	k := getKey(&utxo.Id)
	fmt.Printf("Key %x \n", k)

	store.db.Set(getKey(&utxo.Id), b)
}
