package cmd

import (
	"fmt"
	"homework1/storage"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Команда для вывода возвратов клиента с постраничной пагинацией. Опциаональный флаг size --- размер одной страницы
var UserReturnsCmd = &cobra.Command{
	Use:   "userReturns",
	Short: "Recieve list of all returns",
	Long:  `With this command you can recieve list of all return in pick-point`,
	Run: func(cmd *cobra.Command, args []string) {
		// Парсинг и обработка флагов
		pageSize, err := cmd.Flags().GetInt("size")
		if err != nil {
			fmt.Println(err)
			return
		}
		if pageSize < 1 {
			pageSize = 1
		}

		//Чтение базы данных
		database, err := storage.GetData()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Сборка всех возвратов для вывода
		var returns []string
		for _, v := range database.Orders {
			if v.State == "returned" {
				text := fmt.Sprintf("Order Id: %d, User Id: %d", v.Id, v.UserId)
				returns = append(returns, text)
			}
		}
		if len(returns) == 0 {
			fmt.Println("empty")
			return
		}

		// Подсчет количества страниц в пагинации
		totalPages := len(returns) / pageSize
		if len(returns)%pageSize != 0 {
			totalPages += 1
		}

		numPages := make([]string, totalPages)
		for i := 0; i < totalPages; i++ {
			numPages[i] = strconv.Itoa(i + 1)
		}

		// Запрос выбора страницы
		prompt := promptui.Select{
			Label: "Select Page:",
			Items: numPages,
		}
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Вывод возвратов на странице selectedPage
		selectedPage, _ := strconv.Atoi(result)
		if selectedPage*pageSize >= len(returns) {
			fmt.Println(strings.Join(returns[(selectedPage-1)*pageSize:], "\n"))
		} else {
			fmt.Println(strings.Join(returns[(selectedPage-1)*pageSize:selectedPage*pageSize], "\n"))
		}

	},
}

func init() {
	RootCmd.AddCommand(UserReturnsCmd)
}
