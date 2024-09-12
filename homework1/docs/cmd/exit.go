package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Команда для выхода
var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Quit CLI",
	Long:  `With this command u can quit CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Goodbye!")
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(exitCmd)
}
