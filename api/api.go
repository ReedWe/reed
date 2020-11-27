package api

import (
	"encoding/json"
	"fmt"
	bc "github.com/reed/blockchain"
	"github.com/reed/blockchain/tx/txbuilder"
	"github.com/reed/blockchain/txpool"
	"github.com/reed/database/leveldb"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/types"
	dbm "github.com/tendermint/tmlibs/db"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

var (
	mainURL = "http://127.0.0.1:9888"
)

type API struct {
	Chain  bc.Chain
	Server *http.Server
}

type Res struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func NewApi() *API {

	leveldbStore := leveldb.NewStore(dbm.NewDB("core", dbm.LevelDBBackend, os.Getenv("GOPATH")+"/src/github.com/reed/database/file/"))

	api := &API{}

	//init api server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Welcome to Tiny chain!")
	})
	mux.HandleFunc("/sumbit-transaction", api.SubmitTxHandler)

	httpServer := &http.Server{
		Addr:    mainURL,
		Handler: mux,
	}

	tp := txpool.NewTxpool(leveldbStore)
	api.Chain = bc.Chain{Store: leveldbStore, Txpool: tp}
	api.Server = httpServer

	return api
}

func (a *API) StartApiServer() {
	listen, err := net.Listen("tcp", "0.0.0.0:9888")
	if err != nil {
		log.Logger.Fatalf("failed to start api server %v", err)
	}

	go func() {
		if err := a.Server.Serve(listen); err != nil {
			log.Logger.WithField("error", errors.Wrap(err, "server"))
		}
	}()

}

func PrintSuccessRes(writer http.ResponseWriter, data interface{}) {
	r := &Res{Success: true, Data: data}
	b, err := json.Marshal(r)
	if err != nil {
		log.Logger.Errorf("json marshal %v", err)
	}
	fmt.Fprint(writer, string(b))
}

func PrintErrorRes(writer http.ResponseWriter, data interface{}) {
	r := &Res{Success: false, Data: data}
	b, err := json.Marshal(r)
	if err != nil {
		log.Logger.Errorf("json marshal %v", err)
	}
	fmt.Fprint(writer, string(b))
}

func (a *API) SubmitTxHandler(writer http.ResponseWriter, request *http.Request) {
	if data, err := ioutil.ReadAll(request.Body); err != nil {
		log.Logger.Errorf("sumbit transaction handler %v", err)
		PrintErrorRes(writer, err)
		return
	} else {
		m := &types.SubmitTxRequest{}
		err := json.Unmarshal(data, m)
		if err != nil {
			PrintErrorRes(writer, err)
			return
		}
		txResponse, err := txbuilder.SubmitTx(&a.Chain, m)
		if err != nil {
			log.Logger.Error(err.Error())
			PrintErrorRes(writer, err.Error())
			return
		}
		PrintSuccessRes(writer, txResponse)
	}
}

//
//func jsonHandler(f interface{}) http.Handler {
//	h, err := httpjson.Handler(f, errorFormatter.Write)
//	if err != nil {
//		panic(err)
//	}
//	return h
//}
