package command

import (
	"github.com/reed/blockchain/store"
	"github.com/reed/database/leveldb"
	"github.com/spf13/cobra"
	dbm "github.com/tendermint/tmlibs/db"
	"os"
)

var (
	txCmd = &cobra.Command{
		Use:   "transaction",
		Short: "add a transaction",
		Run:   AddTransaction,
	}

	txName = ""
)

func init() {
	txCmd.Flags().StringVarP(&txName, "name", "n", txName, "a transaction name")
}

func Execute() {
	txCmd.Execute()
}

func AddTransaction(cmd *cobra.Command, args []string) {

	if len(txName) == 0 {
		cmd.Help()
		return
	}
}

func getStore() store.Store {
	return leveldb.NewStore(dbm.NewDB("core", dbm.LevelDBBackend, os.Getenv("GOPATH")+"/src/github.com/reed/database/file/"))
}
