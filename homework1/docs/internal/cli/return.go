package cli

import (
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/usecase"

	"github.com/spf13/cobra"
)

func initReturnCmd(orderUseCase usecase.OrderUseCase) *cobra.Command {
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

			request := &dto.ReturnOrderRequest{
				Id: orderId,
			}

			err = orderUseCase.Return(request)
			if err != nil {
				return err
			}
			fmt.Println("Order succesfull returned to courier")
			return nil
		},
	}
	return ReturnCmd
}