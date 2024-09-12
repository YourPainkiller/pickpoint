package cmd

import (
	"fmt"
	"homework1/storage"
	"time"

	"github.com/spf13/cobra"
)

// Команда для выдачи заказов клиенту. Обязательный флаг ListOfIds --- список Id заказов для выдачи
var GiveCmd = &cobra.Command{
	Use:     "give",
	Short:   "Give list of orders to client",
	Long:    `With this command you can give orders to client.`,
	Example: "give --loi=1,2,3",
	Run: func(cmd *cobra.Command, args []string) {
		// Парсинг и обработка флагов
		orderIds, err := cmd.Flags().GetIntSlice("loi")
		if err != nil {
			fmt.Println(err)
			return
		}

		// Для быстрой проверки того, что нам нужно выдать заказ именно с таким Id
		orderIdsSearch := make(map[int]bool)
		for _, id := range orderIds {
			if _, ok := orderIdsSearch[id]; !ok {
				orderIdsSearch[id] = true
			}
		}
		// Чтение нашей базы данных
		database, err := storage.GetData()
		if err != nil {
			fmt.Println(err)
			return
		}

		var userId int
		var cnt int // Счетчик выданных заказов

		for k, order := range database.Orders {
			if _, ok := orderIdsSearch[order.Id]; ok { // Проверка на то, что заказ из базы должен быть выдан
				curTime := time.Now()
				orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
				if order.State != "accepted" { // Проверка на то, что заказ с таким OrderId был принят
					fmt.Printf("Error: order with id %d can't be taken, because it has been already taken or still didn't come\n", order.Id)
					return
				}
				if curTime.After(orderTime) { // Проверка на время
					fmt.Printf("Error: order with id %d can't be taken, because time left\n", order.Id)
					return
				} else if userId == 0 { // Проверка на то, что все заказы принадлежат одному UserId
					userId = order.UserId
					tmpTime := time.Now().Add(time.Hour * 72).Format(TIMELAYOUT)
					database.Orders[k].State = "gived"     // Помечаем как выданный
					database.Orders[k].ValidTime = tmpTime // Ставим время до которого заказ можно вернуть
					cnt++
				} else if userId != order.UserId {
					fmt.Println("Error: one or more orders are not yours")
					return
				} else {
					tmpTime := time.Now().Add(time.Hour * 72).Format(TIMELAYOUT)
					database.Orders[k].State = "gived"
					database.Orders[k].ValidTime = tmpTime
					cnt++
				}
			}
		}
		if cnt != len(orderIds) {
			fmt.Println("Error: one or more orders are not in our pick-point")
			return
		}
		storage.SendData(database)
		fmt.Println("All of your orders given")
	},
}

func init() {
	RootCmd.AddCommand(GiveCmd)
}
