package cli

import (
	"context"
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/usecase"
	"homework1/internal/workerPool"

	"github.com/spf13/cobra"
)

func initAcceptReturnCmd(orderUseCase usecase.OrderUseCase, pool *workerPool.Pool) *cobra.Command {
	// Команда для прием товара на возврат от клиента. Обязательные флаги UserId, OrderId
	var acceptReturnCmd = &cobra.Command{
		Use:     "acceptReturn",
		Short:   "Accept return from client",
		Long:    `With this command you can try to recieve return from client.`,
		Example: "acceptReturn --ui=1 --oi=15",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Парсинг и обработка флагов
			orderId, err := cmd.Flags().GetInt("oi")
			if err != nil {
				return err
			}

			userId, err := cmd.Flags().GetInt("ui")
			if err != nil {
				return err
			}

			slow, err := cmd.Flags().GetBool("slow")
			if err != nil {
				return err
			}

			request := &dto.AcceptReturnOrderRequest{
				Id:     orderId,
				UserId: userId,
			}

			ctx := context.Background()
			pool.SubmitTask(func() {
				err := orderUseCase.AcceptReturn(ctx, request)
				if err != nil {
					err = fmt.Errorf("error in accpeting return from userid=%d, orderid=%d: %v", request.UserId, request.Id, err)
					fmt.Println(err)
				} else {
					fmt.Println("Return accepted succesfull")
				}
			})
			if slow {
				pool.GetTasksWg().Wait()
			}
			return nil
		},
	}
	return acceptReturnCmd
}
