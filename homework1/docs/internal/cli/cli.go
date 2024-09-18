package cli

import (
	"bufio"
	"fmt"
	"homework1/internal/usecase"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var Reset = "\033[0m"
var Cyan = "\033[36m"

func flagsReset(acceptCmd, acceptReturn, giveCmd, returnCmd, userOrdersCmd, userReturnsCmd *cobra.Command) {
	RootCmd.DisableSuggestions = true
	RootCmd.CompletionOptions.HiddenDefaultCmd = true
	RootCmd.CompletionOptions.DisableDescriptions = true
	RootCmd.ResetFlags()
	acceptCmd.ResetFlags()
	returnCmd.ResetFlags()
	giveCmd.ResetFlags()
	userOrdersCmd.ResetFlags()
	acceptReturn.ResetFlags()
	userReturnsCmd.ResetFlags()
	userReturnsCmd.ResetFlags()

	acceptCmd.Flags().Int("oi", -1, "Order ID, required")
	acceptCmd.Flags().Int("ui", -1, "User ID, required")
	acceptCmd.Flags().String("vt", "2006-01-02", "Valid time in format YYYY-MM-DD, required")
	acceptCmd.Flags().String("package", "none", "Package type (stretch, box or bag), required")
	acceptCmd.Flags().Int("price", -1, "Price of order, required")
	acceptCmd.Flags().Int("weight", -1, "Weight of order, required")
	acceptCmd.Flags().Bool("addstr", false, "Additional stretch, optional")
	returnCmd.Flags().Int("oi", -1, "Order ID, required")
	giveCmd.Flags().IntSlice("loi", []int{}, "Slice of orders ids")
	userOrdersCmd.Flags().Int("ui", -1, "User Id, required")
	userOrdersCmd.Flags().Int("last", -1, "Recieve list of N last orders")
	acceptReturn.Flags().Int("oi", -1, "Order Id, required")
	acceptReturn.Flags().Int("ui", -1, "User Id, required")
	userReturnsCmd.Flags().Int("size", 5, "page size, optional")
	userReturnsCmd.Flags().Int("page", 1, "page, required")

	acceptCmd.MarkFlagRequired("oi")
	acceptCmd.MarkFlagRequired("ui")
	acceptCmd.MarkFlagRequired("vt")
	returnCmd.MarkFlagRequired("oi")
	acceptCmd.MarkFlagRequired("package")
	acceptCmd.MarkFlagRequired("price")
	acceptCmd.MarkFlagRequired("weight")
	giveCmd.MarkFlagRequired("loi")
	userOrdersCmd.MarkFlagRequired("ui")
	acceptReturn.MarkFlagRequired("oi")
	acceptReturn.MarkFlagRequired("ui")
	userReturnsCmd.MarkFlagRequired("page")
}

func Run(orderUseCase usecase.OrderUseCase) {
	acceptCmd := initAcceptCmd(orderUseCase)
	acceptReturn := initAcceptReturnCmd(orderUseCase)
	exitCmd := initExitCmd()
	giveCmd := initGiveCmd(orderUseCase)
	returnCmd := initReturnCmd(orderUseCase)
	userOrdersCmd := initUserOrdersCmd(orderUseCase)
	userReturnsCmd := initUserReturnsCmd(orderUseCase)
	RootCmd = InitRootCmd(acceptCmd, acceptReturn, exitCmd, giveCmd, returnCmd, userOrdersCmd, userReturnsCmd)

	fmt.Println("Welcome to CLI. help for more info")
	flagsReset(acceptCmd, acceptReturn, giveCmd, returnCmd, userOrdersCmd, userReturnsCmd)
	Execute()
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
		flagsReset(acceptCmd, acceptReturn, giveCmd, returnCmd, userOrdersCmd, userReturnsCmd) // Сброс и востановление флагов всех команд
		RootCmd.SetArgs(cmdArgs)
		Execute()

	}
}
