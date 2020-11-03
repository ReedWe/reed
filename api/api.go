package api

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	bc "github.com/tybc/blockchain"
	"github.com/tybc/blockchain/tx"
	"github.com/tybc/core"
	"github.com/tybc/database/leveldb"
	"github.com/tybc/log"
	"io/ioutil"
	"net"
	"net/http"
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

	fmt.Println("new api...")
	leveldbStore := leveldb.NewStore(dbm.NewDB("core", dbm.LevelDBBackend, "/Users/jan/go/src/github.com/tybc/database/file/"))

	api := &API{}

	//init api server
	fmt.Println("init api server")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Welcome to Tiny chain!")
	})
	mux.HandleFunc("/sumbit-transaction", api.SubmitTxHandler)

	fmt.Println("handlefunc complete")

	httpServer := &http.Server{
		Addr:    mainURL,
		Handler: mux,
	}
	api.Chain = bc.Chain{Store: leveldbStore, Txpool: &core.Txpool{}}
	api.Server = httpServer

	return api
}

func (a *API) StartApiServer() {
	fmt.Println("start api server")
	listen, err := net.Listen("tcp", "0.0.0.0:9888")
	if err != nil {
		common.Exit(common.Fmt("failed to start api server %v", err))
	}

	go func() {
		if err := a.Server.Serve(listen); err != nil {
			fmt.Println("Rpc server error")
		}
	}()

	fmt.Println("start api server complete")

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
		m := &tx.SubmitTxRequest{}
		err := json.Unmarshal(data, m)
		if err != nil {
			PrintErrorRes(writer, err)
			return
		}
		log.Logger.Infof("submit tx handler success")

		txResponse, err := tx.SubmitTx(&a.Chain, m)
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
