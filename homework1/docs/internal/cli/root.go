package cli

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "homework1",
	Short: "CLI app for pick-up point",
	Long:  `Simple CLI app for pick-up point using Golang Cobra`,
}

func InitRootCmd(commands ...*cobra.Command) *cobra.Command {
	RootCmd.AddCommand(commands...)
	return RootCmd
}

func Execute() {
	RootCmd.Execute()
}
