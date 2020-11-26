package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/reed/log"
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
	runCmd.Execute()
}

func init() {
	runCmd.Flags().StringVarP(&name, "name", "n", name, "set node name")

}

func RunNode(cmd *cobra.Command, args []string) error {
	n := NewNode()

	if err := n.Start(); err != nil {
		return fmt.Errorf("Failed to start node: %v", err)
	}
	log.Logger.Info("node start...")

	n.RunFover()
	return nil
}
