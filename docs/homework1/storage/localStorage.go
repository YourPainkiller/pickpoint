package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// Пакет для работы с хранилищем: чтением и записью в него json

type Data struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	Id        int    `json:"orderId"`
	UserId    int    `json:"userId"`
	ValidTime string `json:"validTime"`
	State     string `json:"state"`
}

const PATH = "./database.json"

func GetData() (Data, error) {
	jsonFile, err := os.OpenFile(PATH, os.O_CREATE, 0666)
	if err != nil {
		return Data{}, errors.New("uanble to open database")
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var jsonData Data
	err = json.Unmarshal(byteValue, &jsonData) // Если json по какой-то причине на распарсился (в т.ч. файл пустой)
	if err != nil {                            // Возвращаем пустую базу данных
		return Data{}, nil
	}

	return jsonData, nil
}

func SendData(data Data) error {
	if err := os.Truncate(PATH, 0); err != nil {
		return err
	}
	jsonFile, err := os.OpenFile(PATH, os.O_CREATE, 0666)
	if err != nil {
		return errors.New("uanble to open database")
	}
	defer jsonFile.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.New("unable to marshal data")
	}
	if err := os.WriteFile(PATH, jsonData, 0666); err != nil {
		return err
	}
	return nil
}
