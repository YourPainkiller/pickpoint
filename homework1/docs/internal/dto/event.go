package dto

type EventDto struct {
	OrderId int    `json:"orderId"`
	Method  string `json:"method"`
	Time    string `json:"time"`
}
