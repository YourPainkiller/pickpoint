package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"homework1/internal/dto"
	cliserver "homework1/pkg/cli/v1"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func initUserReturnsCmd(cliclient cliserver.CliClient) *cobra.Command {
	// Команда для вывода возвратов клиента с постраничной пагинацией. Обязательный флаг page --- номер страницы для вывода
	// Опциаональный флаг size --- размер одной страницы
	var userReturnsCmd = &cobra.Command{
		Use:   "userReturns",
		Short: "Recieve list of all returns",
		Long:  `With this command you can recieve list of all return in pick-point`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Парсинг и обработка флагов
			page, err := cmd.Flags().GetInt("page")
			if err != nil {
				return err
			}
			if page < 1 {
				page = 1
			}

			pageSize, err := cmd.Flags().GetInt("size")
			if err != nil {
				return err
			}
			if pageSize < 1 {
				pageSize = 1
			}

			request := &dto.UserReturnsRequest{
				Page: page,
				Size: pageSize,
			}
			byteReq, err := json.Marshal(request)
			if err != nil {
				return err
			}

			grpcReq := &cliserver.UserReturnsRequest{}
			if err := protojson.Unmarshal(byteReq, grpcReq); err != nil {
				return err
			}

			ctx := context.Background()
			resp, respErr := cliclient.UserReturnsGrpc(ctx, grpcReq)
			if respErr != nil {
				return err
			}

			if resp == nil {
				fmt.Println("empty")
				return nil
			}

			output := []string{}
			for _, order := range resp.OrderDtos {
				output = append(output, fmt.Sprintf("Order Id: %d, User Id: %d", order.Id, order.UserId))
			}

			fmt.Println(strings.Join(output, "\n"))

			return nil

		},
	}
	return userReturnsCmd
}
