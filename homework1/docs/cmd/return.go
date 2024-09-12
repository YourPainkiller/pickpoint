package cmd

import (
	"fmt"
	"homework1/storage"
	"time"

	"github.com/spf13/cobra"
)

// Команда для возрвата заказа обратно курьеру. Обязательный флаг OrderId
var ReturnCmd = &cobra.Command{
	Use:     "return",
	Short:   "Return order to courier",
	Long:    `With this command you can return order to courier.`,
	Example: "return --oi=1",
	Run: func(cmd *cobra.Command, args []string) {
		// Парсинг и обработка флагов
		orderId, err := cmd.Flags().GetInt("oi")
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

		// Обход нашей базы данных
		for k, order := range database.Orders {
			if order.Id == orderId {
				curTime := time.Now().Add(24 * time.Hour)
				orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
				if curTime.After(orderTime) || order.State == "returned" { // Проверяем что срок хранения истек или заказ был возвращен
					if order.State != "gived" { // Проверяем что заказ не находится у клиента
						database.Orders = append(database.Orders[:k], database.Orders[k+1:]...) //Удаляем из базы
						err := storage.SendData(database)                                       // Обновляем базу
						if err != nil {
							fmt.Println(err)
							return
						}
						fmt.Println("Order succesfull returned to courier")
						return
					} else {
						fmt.Println("Error: this order is with the client")
						return
					}
				} else {
					fmt.Println("Error: client still can take it")
					return
				}
			}
		}
		fmt.Println("Error: no such order")
	},
}

func init() {
	RootCmd.AddCommand(ReturnCmd)
}
