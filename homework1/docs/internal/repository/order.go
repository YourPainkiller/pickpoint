package repository

import (
	"encoding/json"
	"errors"
	"homework1/internal/dto"
	"io"
	"os"
)

//const PATH = "./database.json"

type OrderRepository struct {
	path string
}

var ErrUnableToOpen = errors.New("unable to open database")

func NewOrderRepository(path string) (*OrderRepository, error) {
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return &OrderRepository{path: path}, nil
}

func (r *OrderRepository) InsertOrders(data *dto.ListOrdersDto) error {
	if err := os.Truncate(r.path, 0); err != nil {
		return err
	}
	jsonFile, err := os.OpenFile(r.path, os.O_CREATE, 0666)
	if err != nil {
		return ErrUnableToOpen
	}
	defer jsonFile.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.New("unable to marshal data")
	}
	if err := os.WriteFile(r.path, jsonData, 0666); err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetOrders() (*dto.ListOrdersDto, error) {
	jsonFile, err := os.OpenFile(r.path, os.O_CREATE, 0666)
	if err != nil {
		return &dto.ListOrdersDto{}, ErrUnableToOpen
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var jsonData *dto.ListOrdersDto
	err = json.Unmarshal(byteValue, &jsonData) // Если json по какой-то причине на распарсился (в т.ч. файл пустой)
	if err != nil {                            // Возвращаем пустую базу данных
		return &dto.ListOrdersDto{}, nil
	}

	return jsonData, nil
}
