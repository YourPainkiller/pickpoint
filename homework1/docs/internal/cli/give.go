package cli

import (
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/usecase"

	"github.com/spf13/cobra"
)

func initGiveCmd(orderUseCase usecase.OrderUseCase) *cobra.Command {
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

			request := &dto.GiveOrderRequest{
				OrderIds: orderIds,
			}

			err = orderUseCase.Give(request)
			if err != nil {
				return err
			}

			fmt.Println("All of your orders given")
			return nil
		},
	}
	return giveCmd
}
