package cmd

import (
	"fmt"
	"homework1/storage"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Команда для вывода заказов от клиента с пагинацией скроллом. Обязательный флаг UserId. Необязательный last --- количество последних по получению заказов для вывода
var UserOrdersCmd = &cobra.Command{
	Use:     "userOrders",
	Short:   "Recieve list of all user orders",
	Long:    `With this command you can recieve list of all user orders.`,
	Example: "userOrders --ui=1",
	Run: func(cmd *cobra.Command, args []string) {
		// Парсинг и обработка флагов
		userId, err := cmd.Flags().GetInt("ui")
		if err != nil {
			fmt.Println(err)
			return
		} else if userId < 1 {
			fmt.Println("Error: Bad arguments")
			return
		}
		last, err := cmd.Flags().GetInt("last")
		if err != nil {
			fmt.Println(err)
			return
		}

		// Чтение базы данных
		database, err := storage.GetData()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Поиск и сборка всех заказов от клиента с userId
		var userOrders []string
		for _, order := range database.Orders {
			if order.UserId == userId && order.State != "gived" {
				text := fmt.Sprintf("Order Id: %d, Valid untill: %s, State: %s", order.Id, order.ValidTime, order.State)
				userOrders = append(userOrders, text)
			}
		}
		if len(userOrders) == 0 {
			fmt.Println("empty")
			return
		}

		// Проверка на наличие флага last
		if last == -1 {
			// Составление промпта для выбора страницы. Используется доп пакет promptui
			prompt := promptui.Select{
				Label: "Select Order:",
				Items: userOrders,
			}
			_, result, err := prompt.Run()
			if err != nil {
				fmt.Println(err)
				return
			}
			// TODO: можно добавить доп информацию о выбранном заказе
			fmt.Printf("You choose %s.\nPlace to additional info about order\n", result)
		} else {
			if last > len(userOrders) {
				last = len(userOrders)
			}

			prompt := promptui.Select{
				Label: "Select Order:",
				Items: userOrders[:len(userOrders)-last],
			}
			_, result, err := prompt.Run()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("You choose %s.\nPlace to additional info about order\n", result)
		}
	},
}

func init() {
	RootCmd.AddCommand(UserOrdersCmd)
}
