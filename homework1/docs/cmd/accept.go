package cmd

import (
	"fmt"
	"homework1/storage"
	"time"

	"github.com/spf13/cobra"
)

const TIMELAYOUT = "2006-01-02"

// Команда для приема заказов от курьера. Обязательные флаги OrderId, UserId и ValidTime
var AcceptCmd = &cobra.Command{
	Use:     "accept",
	Short:   "Accept order from courier",
	Long:    `With this command you can accept order to pick-up point.`,
	Example: "accept --oi=1 --ui=1 --vt=2006-01-02",

	Run: func(cmd *cobra.Command, args []string) {
		// Парсинг и обработка флагов
		orderId, err := cmd.Flags().GetInt("oi")
		if err != nil {
			fmt.Println(err)
			return
		} else if orderId < 0 {
			fmt.Println("Error: bad value of argument")
			return
		}
		userId, err := cmd.Flags().GetInt("ui")
		if err != nil {
			fmt.Println(err)
			return
		} else if userId < 0 {
			fmt.Println("Error: bad value of argument")
			return
		}
		date, err := cmd.Flags().GetString("vt")
		if err != nil {
			fmt.Println(err)
			return
		}
		validDate, err := time.Parse(TIMELAYOUT, date)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Проверка на то, что время истечения срока хранения больше текущей даты
		validDate = validDate.Add(24 * time.Hour)
		curTime := time.Now()
		if validDate.Before(curTime) {
			fmt.Println("Error: invalid time")
			return
		}

		// Чтение нашей json базы данных и проверка на наличие заказа с OrderId уже в ней
		database, err := storage.GetData()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, order := range database.Orders {
			if order.Id == orderId {
				fmt.Println("Error: this order already in base")
				return
			}
		}

		// Конвертация времени в строку. До этого добавили к дате 24 часа дабы затем проще сравнивать (обычно срок хранения указывается включительно)
		dt := validDate.Format(TIMELAYOUT)

		// Добавление полученного заказа в нашу базу данных
		database.Orders = append(database.Orders, storage.Order{Id: orderId, UserId: userId, ValidTime: dt, State: "accepted"})
		err = storage.SendData(database)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("OrderId %d, UserId %d, Valid time %s: succesfull accepted\n", orderId, userId, dt)
	},
}

func init() {
	RootCmd.AddCommand(AcceptCmd)
}
