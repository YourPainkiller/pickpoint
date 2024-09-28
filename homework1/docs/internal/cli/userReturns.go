package cli

import (
	"context"
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/usecase"
	"strings"

	"github.com/spf13/cobra"
)

func initUserReturnsCmd(orderUseCase usecase.OrderUseCase) *cobra.Command {
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

			ctx := context.Background()
			response, err := orderUseCase.UserReturns(ctx, request)
			if err != nil {
				return err
			}

			output := []string{}
			for _, order := range response.Orders {
				output = append(output, fmt.Sprintf("Order Id: %d, User Id: %d", order.Id, order.UserId))
			}

			fmt.Println(strings.Join(output, "\n"))

			return nil

		},
	}
	return userReturnsCmd
}
