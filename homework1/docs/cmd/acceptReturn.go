package cmd

import (
	"fmt"
	"homework1/storage"
	"time"

	"github.com/spf13/cobra"
)

// Команда для прием товара на возврат от клиента. Обязательные флаги UserId, OrderId
var AcceptReturnCmd = &cobra.Command{
	Use:     "acceptReturn",
	Short:   "Accept return from client",
	Long:    `With this command you can try to recieve return from client.`,
	Example: "acceptReturn --ui=1 --oi==15",
	Run: func(cmd *cobra.Command, args []string) {
		// Парсинг и обработка флагов
		orderId, err := cmd.Flags().GetInt("oi")
		if err != nil {
			fmt.Println(err)
			return
		} else if orderId < 1 {
			fmt.Println("Error: bad arguments")
			return
		}
		userId, err := cmd.Flags().GetInt("ui")
		if err != nil {
			fmt.Println(err)
			return
		} else if userId < 1 {
			fmt.Println("Error: bad arguments")
		}

		// Чтение нашей базы данных
		database, err := storage.GetData()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Проходим по нашей базе данных
		var check bool
		for k, order := range database.Orders {
			if order.Id == orderId { // Проверяем на наличие заказа
				check = true
				if order.UserId == userId { // Проверяем на совпадение UserId
					if order.State == "gived" { // Проверяем что заказ был выдан клиенту
						orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
						curTime := time.Now()
						if curTime.Before(orderTime) { // Проверяем что заказ еще можно вернуть по времени
							database.Orders[k].State = "returned"
						} else {
							fmt.Println("Error: no time to return")
							return
						}
					} else {
						fmt.Println("Error: your order already returned or still not gived")
						return
					}
				} else {
					fmt.Println("Error: it's not yours order")
					return
				}
			}
		}
		// Проверяем вернули ли мы заказ, если да, то обновляем базу
		if !check {
			fmt.Println("Error: No such order")
		} else {
			err = storage.SendData(database)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Return accepted succesfull")
		}
	},
}

func init() {
	RootCmd.AddCommand(AcceptReturnCmd)
}
