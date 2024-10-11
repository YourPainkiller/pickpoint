package cli

import (
	"bufio"
	"fmt"
	"homework1/internal/usecase"
	"homework1/internal/workerPool"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var Reset = "\033[0m"
var Cyan = "\033[36m"

const MAXWORKERS = 10000
const MAXTASKS = 1000

func flagsReset(acceptCmd, acceptReturn, giveCmd, returnCmd, userOrdersCmd, userReturnsCmd, changeCmd *cobra.Command) {
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
	changeCmd.ResetFlags()

	acceptCmd.Flags().Int("oi", -1, "Order ID, required")
	acceptCmd.Flags().Bool("slow", false, "Run command not cuncurrnetly")
	acceptCmd.Flags().Int("ui", -1, "User ID, required")
	acceptCmd.Flags().String("vt", "2006-01-02", "Valid time in format YYYY-MM-DD, required")
	acceptCmd.Flags().String("package", "none", "Package type (stretch, box or bag), required")
	acceptCmd.Flags().Int("price", -1, "Price of order, required")
	acceptCmd.Flags().Int("weight", -1, "Weight of order, required")
	acceptCmd.Flags().Bool("addstr", false, "Additional stretch, optional")
	returnCmd.Flags().Int("oi", -1, "Order ID, required")
	returnCmd.Flags().Bool("slow", false, "Run command not cuncurrnetly")
	giveCmd.Flags().IntSlice("loi", []int{}, "Slice of orders ids")
	giveCmd.Flags().Bool("slow", false, "Run command not cuncurrnetly")
	userOrdersCmd.Flags().Int("ui", -1, "User Id, required")
	userOrdersCmd.Flags().Int("last", -1, "Recieve list of N last orders")
	acceptReturn.Flags().Int("oi", -1, "Order Id, required")
	acceptReturn.Flags().Int("ui", -1, "User Id, required")
	acceptReturn.Flags().Bool("slow", false, "Run command not cuncurrnetly")
	userReturnsCmd.Flags().Int("size", 5, "page size, optional")
	userReturnsCmd.Flags().Int("page", 1, "page, required")
	changeCmd.Flags().Int("delta", 0, "number of workers to add or to delete")

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
	changeCmd.MarkFlagRequired("delta")
}

func Run(orderUseCase usecase.OrderUseCase) {
	//Создаем пул и добавляем в него 5 воркеров
	pool := workerPool.NewPool(MAXWORKERS, MAXTASKS)
	for i := 0; i < 5; i++ {
		go pool.CreateWorker()
	}

	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)

	acceptCmd := initAcceptCmd(orderUseCase, pool)
	acceptReturn := initAcceptReturnCmd(orderUseCase, pool)
	exitCmd := initExitCmd()
	giveCmd := initGiveCmd(orderUseCase, pool)
	returnCmd := initReturnCmd(orderUseCase, pool)
	userOrdersCmd := initUserOrdersCmd(orderUseCase)
	userReturnsCmd := initUserReturnsCmd(orderUseCase)
	changeCmd := initChangeCmd(pool)
	RootCmd = InitRootCmd(acceptCmd, acceptReturn, exitCmd, giveCmd, returnCmd, userOrdersCmd, userReturnsCmd, changeCmd)

	fmt.Println("Welcome to CLI. help for more info")
	flagsReset(acceptCmd, acceptReturn, giveCmd, returnCmd, userOrdersCmd, userReturnsCmd, changeCmd)
	Execute()
	reader := bufio.NewReader(os.Stdin)

	go func() { // Запуск бесконечного чтения нашего терминала при помощи кобры
		for {
			time.Sleep(50 * time.Millisecond)
			fmt.Print(Cyan + "homework1 " + Reset)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}

			cmdArgs := strings.Split(input, " ")
			flagsReset(acceptCmd, acceptReturn, giveCmd, returnCmd, userOrdersCmd, userReturnsCmd, changeCmd) // Сброс и востановление флагов всех команд
			RootCmd.SetArgs(cmdArgs)
			Execute()

		}
	}()
	<-chSig
	fmt.Println("waiting for tasks")
	pool.GetTasksWg().Wait()
	//можно добавить ожидание еще чего-то
	fmt.Println("gracefullshutdown")
}
