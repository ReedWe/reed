package api

import (
	"fmt"
	"github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	bc "github.com/tybc/blockchain"
	"github.com/tybc/database/leveldb"
	"net/http"

	"net"
)

var (
	mainURL = "http://127.0.0.1:9888"
)

type API struct {
	store  bc.Store
	server *http.Server
}

type ApiRes struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func NewApi() *API {

	fmt.Println("new api...")
	leveldbStore := leveldb.NewStore(dbm.NewDB("core", dbm.LevelDBBackend, "/Users/jan/go/src/github.com/tybc/database/file/"))

	//init api server
	fmt.Println("init api server")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Welcome to Tiny chain!")
	})

	fmt.Println("handlefunc complete")

	httpServer := &http.Server{
		Addr:    mainURL,
		Handler: mux,
	}

	api := &API{
		store:  leveldbStore,
		server: httpServer,
	}
	return api
}

func (a *API) StartApiServer() {
	fmt.Println("start api server")
	listen, err := net.Listen("tcp", "0.0.0.0:9888")
	if err != nil {
		common.Exit(common.Fmt("failed to start api server %v", err))
	}

	go func() {
		if err := a.server.Serve(listen); err != nil {
			fmt.Println("Rpc server error")
		}
	}()

	fmt.Println("start api server complete")

}

func NewSuccessRes(data interface{}) ApiRes {
	return ApiRes{Success: true, Data: data}
}

//
//func jsonHandler(f interface{}) http.Handler {
//	h, err := httpjson.Handler(f, errorFormatter.Write)
//	if err != nil {
//		panic(err)
//	}
//	return h
//}
