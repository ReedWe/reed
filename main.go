package main

import (
	"crypto/rand"
	"fmt"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tybc/cmd/node"
	"github.com/tybc/crypto"
	"github.com/tybc/database/leveldb"

	"github.com/tybc/core/types"
	"golang.org/x/crypto/ed25519"
)

func storeWrite() {

	m := types.TxOutput{
		Id:         []byte("123456789"),
		IsCoinBase: false,
		Address:    []byte{},
		Amount:     100,
		ScriptPk:   []byte("a1ba6fe9e9388c10f2f30ec5329911fd043b3b49d4266b24fb8f5e2551859eb8"),
	}

	var coreDB = dbm.NewDB("core", dbm.LevelDBBackend, "/Users/jan/go/src/github.com/tybc/database/file/")
	s := leveldb.NewStore(coreDB)
	s.SaveUtxo(&m)
	fmt.Println("Success！！")

}

func storeRead() {
	var coreDB = dbm.NewDB("core", dbm.LevelDBBackend, "/Users/jan/go/src/github.com/tybc/database/file/")
	s := leveldb.NewStore(coreDB)
	ot, _ := s.GetUtxo([]byte("123456789"))
	fmt.Printf("%v", ot)

}

func main() {

	hash := crypto.Sha256([]byte("123"), []byte("abc"))
	fmt.Println(hash)
	//fmt.Printf("%X \n", hash)

	pk, prk, _ := ed25519.GenerateKey(rand.Reader)

	signMsg := ed25519.Sign(prk, []byte("GO LANG"))
	//
	//fmt.Printf("%x \n", pk)
	//fmt.Printf("%x \n", prk)
	//
	//fmt.Printf("%x \n", signMsg)

	//pk := []byte("a1ba6fe9e9388c10f2f30ec5329911fd043b3b49d4266b24fb8f5e2551859eb8")
	//fmt.Println(len(pk))
	////prk := []byte("b19645016b9dc0dfcd272f718281568d7de4a5bc8e6acaea25722e29d1cd6e8da1ba6fe9e9388c10f2f30ec5329911fd043b3b49d4266b24fb8f5e2551859eb8")
	//
	//signMsg := []byte("d38c8862cb63d4b0e0dd5317371fc54d61f1c6d98ba0f6245769b8f3c6a5221e202c2931f8d33c8f55ce22ca6b40bba7b0a28ed8e1e583b5aeac7b473fc68807")
	pass := ed25519.Verify(pk, []byte("GO LANG"), signMsg)

	fmt.Println(pass)

	node.Execute()
}
