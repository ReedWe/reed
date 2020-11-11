package leveldb

import (
	"github.com/golang/protobuf/proto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tybc/core/types"
	"github.com/tybc/errors"
)

var (
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

func getKey(id *[]byte) []byte {
	return []byte(utxoPrefix + string(*id))
}

func (store *Store) GetUtxo(id []byte) (*types.UTXO, error) {
	var utxo types.UTXO
	data := store.db.Get(getKey(&id))
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
	//k := getKey(&utxo.Id)
	//fmt.Printf("Key %x \n", k)
	//
	//store.db.Set(getKey(&utxo.Id), b)
}
