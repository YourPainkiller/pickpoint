package cli

import (
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/usecase"

	"github.com/spf13/cobra"
)

func initAcceptReturnCmd(orderUseCase usecase.OrderUseCase) *cobra.Command {
	// Команда для прием товара на возврат от клиента. Обязательные флаги UserId, OrderId
	var acceptReturnCmd = &cobra.Command{
		Use:     "acceptReturn",
		Short:   "Accept return from client",
		Long:    `With this command you can try to recieve return from client.`,
		Example: "acceptReturn --ui=1 --oi==15",
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

			request := &dto.AcceptReturnOrderRequest{
				Id:     orderId,
				UserId: userId,
			}

			err = orderUseCase.AcceptReturn(request)
			if err != nil {
				return err
			}
			fmt.Println("Return accepted succesfull")
			return nil
		},
	}
	return acceptReturnCmd
}
