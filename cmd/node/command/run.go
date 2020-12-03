package command

import (
	"github.com/reed/errors"
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
	runCmd.Execute()
}

func init() {
	runCmd.Flags().StringVarP(&name, "name", "n", name, "set node name")
}

func RunNode(cmd *cobra.Command, args []string) error {
	n := NewNode()

	if err := n.Start(); err != nil {
		return errors.Wrapf(err, "Failed to start node")
	}

	n.RunFover()
	return nil
}
