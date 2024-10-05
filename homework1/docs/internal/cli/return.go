package cli

import (
	"context"
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/usecase"
	"homework1/internal/workerPool"

	"github.com/spf13/cobra"
)

func initReturnCmd(orderUseCase usecase.OrderUseCase, pool *workerPool.Pool) *cobra.Command {
	// Команда для возрвата заказа обратно курьеру. Обязательный флаг OrderId
	var ReturnCmd = &cobra.Command{
		Use:     "return",
		Short:   "Return order to courier",
		Long:    `With this command you can return order to courier.`,
		Example: "return --oi=1",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Парсинг и обработка флагов
			orderId, err := cmd.Flags().GetInt("oi")
			if err != nil {
				return err
			}

			slow, err := cmd.Flags().GetBool("slow")
			if err != nil {
				return err
			}

			request := &dto.ReturnOrderRequest{
				Id: orderId,
			}

			ctx := context.Background()
			pool.SubmitTask(func() {
				err := orderUseCase.Return(ctx, request)
				if err != nil {
					err = fmt.Errorf("error in returning order with orderid=%d: %v", request.Id, err)
					fmt.Println(err)
				} else {
					fmt.Println("Order succesfull returned to courier")
				}
			})
			if slow {
				pool.GetTasksWg().Wait()
			}

			return nil
		},
	}
	return ReturnCmd
}
