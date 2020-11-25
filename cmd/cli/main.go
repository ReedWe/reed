package main

import (
	"encoding/json"
	"github.com/tybc/cmd/cli/command"
	"github.com/tybc/log"
	"github.com/tybc/types"
)

func main() {

	log.Init()

	var data = `{"tx_inputs":[{"spend_output_id":"b19645016b9dc0dfcd272f718281568d7de4a5bc8e6acaea25722e29d1cd6e8d"}],"tx_outputs":[{"address":"d1cd6e8da1ba6fe9e9388c10f2f30ec5329911fd043b3b49d4266b24fb8f5e25","amount":120}]}`

	m := &types.SubmitTxRequest{}
	err := json.Unmarshal([]byte(data), m)
	if err != nil {
		log.Logger.Errorf("%v", err)
	}
	command.Call("/submit-transaction", m)

}
