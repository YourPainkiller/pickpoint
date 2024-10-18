package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/workerPool"
	cliserver "homework1/pkg/cli/v1"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func initReturnCmd(cliclient cliserver.CliClient, pool *workerPool.Pool) *cobra.Command {
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

			byteReq, err := json.Marshal(request)
			if err != nil {
				return err
			}

			grpcReq := &cliserver.ReturnOrderRequest{}
			if err := protojson.Unmarshal(byteReq, grpcReq); err != nil {
				return err
			}

			pool.SubmitTask(func() {
				ctx := context.Background()
				_, respErr := cliclient.ReturnOrderGrpc(ctx, grpcReq)
				if respErr != nil {
					fmt.Printf("error in returning order with orderid=%d: %v", request.Id, err)
					return
				}

				fmt.Println("Order succesfull returned to courier")

			})
			if slow {
				pool.GetTasksWg().Wait()
			}

			return nil
		},
	}
	return ReturnCmd
}
