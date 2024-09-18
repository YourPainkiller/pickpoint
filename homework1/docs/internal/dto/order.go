package dto

type OrderDto struct {
	Id                int    `json:"orderId"`
	UserId            int    `json:"userId"`
	ValidTime         string `json:"validTime"`
	State             string `json:"state"`
	Price             int    `json:"price"`
	Weight            int    `json:"weight"`
	PackageType       string `json:"packageType"`
	AdditionalStretch bool   `json:"additionalStretch"`
}

type ListOrdersDto struct {
	Orders []OrderDto `json:"orders"`
}

type AcceptOrderRequest struct {
	Id                int    `json:"orderId"`
	UserId            int    `json:"userId"`
	ValidTime         string `json:"validTime"`
	Price             int    `json:"price"`
	Weight            int    `json:"weight"`
	PackageType       string `json:"packageType"`
	AdditionalStretch bool   `json:"additionalStretch"`
}

type AcceptOrderResponse struct {
	OrderDto
}

type AcceptReturnOrderRequest struct {
	Id     int `json:"orderId"`
	UserId int `json:"userId"`
}

type OrderId struct {
	Id int `json:"orderId"`
}

type GiveOrderRequest struct {
	OrderIds []OrderId `json:"orderIds"`
}

type ReturnOrderRequest struct {
	Id int `json:"orderId"`
}

type UserOrdersRequest struct {
	UserId int `json:"userId"`
	Last   int `json:"last"`
}

type UserOrdersResponse struct {
	ListOrdersDto
}

type UserReturnsRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type UserReturnsResponse struct {
	ListOrdersDto
}
