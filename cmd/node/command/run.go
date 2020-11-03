package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "node",
		Short: "run the tinychain node",
		RunE:  RunNode,
	}

	name = ""
)

func Execute() {
	fmt.Println("Node Execute...")
	runCmd.Execute()
}

func init() {
	fmt.Println("Node init...")
	runCmd.Flags().StringVarP(&name, "name", "n", name, "set node name")

}

func RunNode(cmd *cobra.Command, args []string) error {
	fmt.Println("Node run....")

	n := NewNode()

	if err := n.Start(); err != nil {
		return fmt.Errorf("Failed to start node: %v", err)
	} else {
		fmt.Println("Start Node")
	}
	//if len(name) == 0 {
	//	cmd.Help()
	//	return
	//}
	n.RunFover()

	//var coreDB = dbm.NewDB("core", dbm.LevelDBBackend, "/Users/jan/go/src/github.com/tybc/database/file/")
	//store := leveldb.NewStore(coreDB)
	//fmt.Println("store=", store)

	return nil
}
