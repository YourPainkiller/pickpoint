package dto

type OrderDto struct {
	Id                int    `json:"orderId" db:"order_id" fake:"skip"`
	UserId            int    `json:"userId" db:"user_id" fake:"{number:1,10000}"`
	ValidTime         string `json:"validTime" db:"valid_time" fake:"skip"`
	State             string `json:"state" db:"order_state" fake:"{randomstring:[accepted, gived, returned, deleted]}"`
	Price             int    `json:"price" db:"price" fake:"{number:1,10000}"`
	Weight            int    `json:"weight" db:"weight" fake:"{number:1,10000}"`
	PackageType       string `json:"packageType" db:"package" fake:"{randomstring:[box, bag, stretch]}"`
	AdditionalStretch bool   `json:"additionalStretch" db:"additional_stretch" fake:"skip"`
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
