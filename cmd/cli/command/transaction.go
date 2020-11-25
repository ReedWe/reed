package command

import (
	"github.com/spf13/cobra"
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

func AddTransaction(cmd *cobra.Command, args []string) {
	if len(txName) == 0 {
		cmd.Help()
		return
	}
}
