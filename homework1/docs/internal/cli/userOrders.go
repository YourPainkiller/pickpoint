package cli

import (
	"context"
	"fmt"
	"homework1/internal/dto"
	"homework1/internal/usecase"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func initUserOrdersCmd(orderUseCase usecase.OrderUseCase) *cobra.Command {
	// Команда для вывода заказов от клиента с пагинацией скроллом. Обязательный флаг UserId. Необязательный last --- количество последних по получению заказов для вывода
	var userOrdersCmd = &cobra.Command{
		Use:     "userOrders",
		Short:   "Recieve list of all user orders",
		Long:    `With this command you can recieve list of all user orders.`,
		Example: "userOrders --ui=1",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Парсинг и обработка флагов
			userId, err := cmd.Flags().GetInt("ui")
			if err != nil {
				return err
			}

			last, err := cmd.Flags().GetInt("last")
			if err != nil {
				return err
			}

			request := &dto.UserOrdersRequest{
				UserId: userId,
				Last:   last,
			}

			ctx := context.Background()
			response, err := orderUseCase.UserOrders(ctx, request)
			if err != nil {
				return err
			}

			if response == nil {
				fmt.Println("empty")
				return nil
			}

			text := []string{}
			for _, order := range response.Orders {
				text = append(text, fmt.Sprintf("Order Id: %d, Valid untill: %s, State: %s", order.Id, order.ValidTime, order.State))
			}

			// Составление промпта для выбора страницы. Используется доп пакет promptui
			prompt := promptui.Select{
				Label: "Select Order:",
				Items: text,
			}
			_, result, err := prompt.Run()
			if err != nil {
				return err
			}
			// TODO: можно добавить доп информацию о выбранном заказе
			fmt.Printf("You choose %s.\nPlace to additional info about order\n", result)
			return nil
		},
	}
	return userOrdersCmd
}
