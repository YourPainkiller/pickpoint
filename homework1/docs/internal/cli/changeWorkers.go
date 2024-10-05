package cli

import (
	"fmt"
	"homework1/internal/workerPool"

	"github.com/spf13/cobra"
)

func initChangeCmd(pool *workerPool.Pool) *cobra.Command {
	// Команда для приема заказов от курьера. Обязательные флаги OrderId, UserId и ValidTime
	var ChangeCmd = &cobra.Command{
		Use:     "change",
		Short:   "Change number of workers",
		Long:    `With this command you can change number of workers.`,
		Example: "change --delta=-2",

		RunE: func(cmd *cobra.Command, args []string) error {
			// Парсинг и обработка флагов
			delta, err := cmd.Flags().GetInt("delta")
			if err != nil {
				return err
			}

			if delta < 0 {
				delta *= -1
				curWorkers := int(pool.GetCurrentWorkers())
				if delta > curWorkers {
					fmt.Printf("cant delete more workers than we have (%d)\n", curWorkers)
					return nil
				}

				for j := 0; j < delta; j++ {
					pool.StopWorker()
				}
				pool.GetWorkersWg().Wait()
			} else {
				curWorkers := int(pool.GetCurrentWorkers())
				maxWorkers := pool.GetMaxWorkers()
				if curWorkers+delta > maxWorkers {
					fmt.Printf("cant create more workers. Max workers =%d, want create %d + %d=%d\n", maxWorkers, curWorkers, delta, curWorkers+delta)
					return nil
				}
				for j := 0; j < delta; j++ {
					go pool.CreateWorker()
				}
			}
			return nil
		},
	}

	return ChangeCmd
}

// func init() {
// 	RootCmd.AddCommand(AcceptCmd)
// }
