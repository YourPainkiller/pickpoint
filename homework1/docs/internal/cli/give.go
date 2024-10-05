package cli

import (
	"context"
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/usecase"
	"homework1/internal/workerPool"

	"github.com/spf13/cobra"
)

func initGiveCmd(orderUseCase usecase.OrderUseCase, pool *workerPool.Pool) *cobra.Command {
	// Команда для выдачи заказов клиенту. Обязательный флаг ListOfIds --- список Id заказов для выдачи
	var giveCmd = &cobra.Command{
		Use:     "give",
		Short:   "Give list of orders to client",
		Long:    `With this command you can give orders to client.`,
		Example: "give --loi=1,2,3",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Парсинг и обработка флагов
			oIds, err := cmd.Flags().GetIntSlice("loi")
			if err != nil {
				return err
			}
			orderIds := []dto.OrderId{}
			for _, id := range oIds {
				orderIds = append(orderIds, dto.OrderId{Id: id})
			}

			slow, err := cmd.Flags().GetBool("slow")
			if err != nil {
				return err
			}

			request := &dto.GiveOrderRequest{
				OrderIds: orderIds,
			}

			ctx := context.Background()
			pool.SubmitTask(func() {
				err := orderUseCase.Give(ctx, request)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("All of your orders given")
				}
			})
			if slow {
				pool.GetTasksWg().Wait()
			}
			return nil
		},
	}
	return giveCmd
}
