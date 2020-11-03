package command

import (
	"fmt"
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
	fmt.Println("transaction init...")
	txCmd.Flags().StringVarP(&txName, "name", "n", txName, "a transaction name")
}

func AddTransaction(cmd *cobra.Command, args []string) {
	fmt.Println("AddTransaction...")
	if len(txName) == 0 {
		cmd.Help()
		return
	}
}
