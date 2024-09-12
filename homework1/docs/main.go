package main

import (
	"bufio"
	"fmt"
	"homework1/cmd"
	"os"
	"strings"
)

var Reset = "\033[0m"
var Cyan = "\033[36m"

func main() {
	fmt.Println("Welcome to CLI. help for more info")
	flagsReset()
	cmd.Execute()
	reader := bufio.NewReader(os.Stdin)

	// Запуск бесконечного чтения нашего терминала при помощи кобры
	for {
		fmt.Print(Cyan + "homework1 " + Reset)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cmdArgs := strings.Split(input, " ")
		flagsReset() // Сброс и востановление флагов всех команд
		cmd.RootCmd.SetArgs(cmdArgs)
		cmd.Execute()

	}
}

func flagsReset() {
	cmd.RootCmd.DisableSuggestions = true
	cmd.RootCmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.RootCmd.CompletionOptions.DisableDescriptions = true
	cmd.RootCmd.ResetFlags()
	cmd.AcceptCmd.ResetFlags()
	cmd.ReturnCmd.ResetFlags()
	cmd.GiveCmd.ResetFlags()
	cmd.UserOrdersCmd.ResetFlags()
	cmd.AcceptReturnCmd.ResetFlags()
	cmd.UserReturnsCmd.ResetFlags()

	cmd.AcceptCmd.Flags().Int("oi", -1, "Order ID, required")
	cmd.AcceptCmd.Flags().Int("ui", -1, "User ID, required")
	cmd.AcceptCmd.Flags().String("vt", "2006-01-02", "Valid time in format YYYY-MM-DD, required")
	cmd.ReturnCmd.Flags().Int("oi", -1, "Order ID, required")
	cmd.GiveCmd.Flags().IntSlice("loi", []int{}, "Slice of orders ids")
	cmd.UserOrdersCmd.Flags().Int("ui", -1, "User Id, required")
	cmd.UserOrdersCmd.Flags().Int("last", -1, "Recieve list of N last orders")
	cmd.AcceptReturnCmd.Flags().Int("oi", -1, "Order Id, required")
	cmd.AcceptReturnCmd.Flags().Int("ui", -1, "User Id, required")
	cmd.UserReturnsCmd.Flags().Int("size", 5, "page size, optional")

	cmd.AcceptCmd.MarkFlagRequired("oi")
	cmd.AcceptCmd.MarkFlagRequired("ui")
	cmd.AcceptCmd.MarkFlagRequired("vt")
	cmd.ReturnCmd.MarkFlagRequired("oi")
	cmd.GiveCmd.MarkFlagRequired("loi")
	cmd.UserOrdersCmd.MarkFlagRequired("ui")
	cmd.AcceptReturnCmd.MarkFlagRequired("oi")
	cmd.AcceptReturnCmd.MarkFlagRequired("ui")
}
