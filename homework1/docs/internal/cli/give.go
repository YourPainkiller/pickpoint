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

func initGiveCmd(cliclient cliserver.CliClient, pool *workerPool.Pool) *cobra.Command {
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

			byteReq, err := json.Marshal(request)
			if err != nil {
				return err
			}

			grpcReq := &cliserver.GiveOrderRequest{}
			if err := protojson.Unmarshal(byteReq, grpcReq); err != nil {
				return err
			}

			pool.SubmitTask(func() {
				ctx := context.Background()
				_, respErr := cliclient.GiveOrderGrpc(ctx, grpcReq)
				if respErr != nil {
					fmt.Println("error in giving orders:", respErr)
					return
				}
				fmt.Println("All of your orders given")
			})
			if slow {
				pool.GetTasksWg().Wait()
			}
			return nil
		},
	}
	return giveCmd
}
