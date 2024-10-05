package cli

import (
	"context"
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/usecase"
	"homework1/internal/workerPool"

	"github.com/spf13/cobra"
)

const TIMELAYOUT = "2006-01-02"

func initAcceptCmd(orderUseCase usecase.OrderUseCase, pool *workerPool.Pool) *cobra.Command {
	// Команда для приема заказов от курьера. Обязательные флаги OrderId, UserId и ValidTime
	var AcceptCmd = &cobra.Command{
		Use:     "accept",
		Short:   "Accept order from courier",
		Long:    `With this command you can accept order to pick-up point.`,
		Example: "accept --oi=1 --ui=1 --vt=2006-01-02 --package=box --price=777 --weight=600",

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

			date, err := cmd.Flags().GetString("vt")
			if err != nil {
				return err
			}

			price, err := cmd.Flags().GetInt("price")
			if err != nil {
				return err
			}

			weight, err := cmd.Flags().GetInt("weight")
			if err != nil {
				return err
			}

			pack, err := cmd.Flags().GetString("package")
			if err != nil {
				return err
			}

			additionalStretch, err := cmd.Flags().GetBool("addstr")
			if err != nil {
				return err
			}

			slow, err := cmd.Flags().GetBool("slow")
			if err != nil {
				return err
			}

			request := &dto.AcceptOrderRequest{
				Id:                orderId,
				UserId:            userId,
				ValidTime:         date,
				Price:             price,
				Weight:            weight,
				PackageType:       pack,
				AdditionalStretch: additionalStretch,
			}

			ctx := context.Background()
			pool.SubmitTask(func() {
				err := orderUseCase.Accept(ctx, request)
				//time.Sleep(5 * time.Second) testing gracefull shutdown
				if err != nil {
					err = fmt.Errorf("error in accpeting order with id=%d: %v", request.Id, err)
					fmt.Println(err)
				} else {
					fmt.Printf("Order with id=%d accepted\n", request.Id)
				}
			})
			if slow {
				pool.GetTasksWg().Wait()
			}

			return nil
		},
	}

	return AcceptCmd
}

// func init() {
// 	RootCmd.AddCommand(AcceptCmd)
// }
