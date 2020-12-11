package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	bc "github.com/reed/blockchain"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/types"
	"github.com/reed/wallet"
	cmn "github.com/tendermint/tmlibs/common"
	"io/ioutil"
	"net"
	"net/http"
)

var (
	mainURL = "http://127.0.0.1:9888"
)

type API struct {
	Chain  *bc.Chain
	Server *http.Server
}

type Res struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func NewApi(c *bc.Chain) *API {

	api := &API{}

	//init api server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Welcome to Reed chain!")
	})
	mux.HandleFunc("/send", api.SubmitTxHandler)

	httpServer := &http.Server{
		Addr:    mainURL,
		Handler: mux,
	}

	api.Chain = c
	api.Server = httpServer

	return api
}

func (a *API) StartApiServer() {
	listen, err := net.Listen("tcp", ":9888")
	if err != nil {
		cmn.Exit(cmn.Fmt("failed to start api server:%v", err))
	}

	go func() {
		if err := a.Server.Serve(listen); err != nil {
			log.Logger.WithField("error", errors.Wrap(err,
				"server"))
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
		if bytes.Equal(data, []byte("block")) {
			PrintSuccessRes(writer, "args is 'block'")
			ch := a.Chain.GetWriteReceptionChan()

			s := a.Chain.Store
			pre, _ := (*s).GetHighestBlock()

			newBlock := &types.Block{
				BlockHeader:  *pre.Copy(),
				Transactions: []*types.Tx{},
			}

			w, _ := wallet.My("123")
			cbTx, _ := types.NewCoinbaseTx(newBlock.Height, w.Pub, bc.CalcCoinbaseAmt(newBlock.Height))
			var txs []*types.Tx
			txs = append(txs, cbTx)
			newBlock.Transactions = txs

			fmt.Println("send new block...")
			ch <- &types.RecvWrap{Block: newBlock, SendBreakWork: true}
		} else if bytes.Equal(data, []byte("h")) {
			s := a.Chain.Store
			highest, _ := (*s).GetHighestBlock()
			PrintSuccessRes(writer, highest.Height)
		}else {
			PrintSuccessRes(writer, "nothing")

		}

		//m := &types.SubmitTxRequest{}
		//err := json.Unmarshal(data, m)
		//if err != nil {
		//	PrintErrorRes(writer, err)
		//	return
		//}
		//txResponse, err := txbuilder.SubmitTx(a.Chain, m)
		//if err != nil {
		//	log.Logger.Error(err.Error())
		//	PrintErrorRes(writer, err.Error())
		//	return
		//}
		//PrintSuccessRes(writer, txResponse)
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
