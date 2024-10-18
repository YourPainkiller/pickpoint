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

func initAcceptReturnCmd(cliclient cliserver.CliClient, pool *workerPool.Pool) *cobra.Command {
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

			byteReq, err := json.Marshal(request)
			if err != nil {
				return err
			}

			grpcReq := &cliserver.AcceptReturnRequest{}
			if err := protojson.Unmarshal(byteReq, grpcReq); err != nil {
				return err
			}

			pool.SubmitTask(func() {
				ctx := context.Background()
				_, respErr := cliclient.AcceptReturnGrpc(ctx, grpcReq)
				if respErr != nil {
					err = fmt.Errorf("error in accpeting return from userid=%d, orderid=%d: %v", request.UserId, request.Id, err)
					fmt.Println(err)
					return
				}

				fmt.Println("Return accepted succesfull")
			})
			if slow {
				pool.GetTasksWg().Wait()
			}
			return nil
		},
	}
	return acceptReturnCmd
}
