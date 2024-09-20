package usecase

import (
	"errors"
	"fmt"
	"homework1/internal/domain"
	"homework1/internal/domain/strategy"
	"homework1/internal/dto"
	"time"
)

const TIMELAYOUT = "2006-01-02"

type (
	orderRepository interface {
		InsertOrders(data *dto.ListOrdersDto) error
		GetOrders() (*dto.ListOrdersDto, error)
	}
)

type OrderUseCase struct {
	repo orderRepository
}

func NewOrderUseCase(repo orderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (oc *OrderUseCase) Accept(req *dto.AcceptOrderRequest) error {
	ord, err := oc.repo.GetOrders()
	if err != nil {
		return err
	}

	for _, order := range ord.Orders {
		if order.Id == req.Id {
			return errors.New("this order already in base")
		}
	}

	var opackageStrategy strategy.OrderPackageStrategy
	switch req.PackageType {
	case domain.TypeBox:
		opackageStrategy = strategy.BoxPackageStrategy{}
	case domain.TypeBag:
		opackageStrategy = strategy.BagPackageStrategy{}
	case domain.TypeStretch:
		opackageStrategy = strategy.StretchPackageStrategy{}
	default:
		return fmt.Errorf("unknown box type: %s", req.PackageType)
	}

	newOrder, err := domain.NewOrder(req.Id, req.UserId, req.Price, req.Weight, req.ValidTime, "accepted", req.PackageType, opackageStrategy, req.AdditionalStretch)
	if err != nil {
		return err
	}
	ord.Orders = append(ord.Orders, newOrder.ToDTO())
	err = oc.repo.InsertOrders(ord)
	if err != nil {
		return err
	}

	return nil
}

func (oc *OrderUseCase) AcceptReturn(req *dto.AcceptReturnOrderRequest) error {
	ord, err := oc.repo.GetOrders()
	if err != nil {
		return err
	}

	// Проходим по нашей базе данных
	var check bool
	for k, order := range ord.Orders {
		if order.Id == req.Id { // Проверяем на наличие заказа
			check = true
			if order.UserId == req.UserId { // Проверяем на совпадение UserId
				if order.State == "gived" { // Проверяем что заказ был выдан клиенту
					orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
					curTime := time.Now()
					if curTime.Before(orderTime) { // Проверяем что заказ еще можно вернуть по времени
						ord.Orders[k].State = "returned"
					} else {
						return errors.New("no time to return")
					}
				} else {
					return errors.New("your order already returned or still not gived")
				}
			} else {
				return errors.New("it's not yours order")
			}
		}
	}
	// Проверяем вернули ли мы заказ, если да, то обновляем базу
	if !check {
		return errors.New("no such order")
	} else {
		err = oc.repo.InsertOrders(ord)
		if err != nil {
			return err
		}
		return nil
	}
}

func (oc *OrderUseCase) Give(req *dto.GiveOrderRequest) error {
	// Для быстрой проверки того, что нам нужно выдать заказ именно с таким Id
	orderIdsSearch := make(map[int]bool)
	for _, id := range req.OrderIds {
		if _, ok := orderIdsSearch[id.Id]; !ok {
			orderIdsSearch[id.Id] = true
		}
	}

	ord, err := oc.repo.GetOrders()
	if err != nil {
		return err
	}

	var userId int
	var cnt int // Счетчик выданных заказов
	for k, order := range ord.Orders {
		if _, ok := orderIdsSearch[order.Id]; ok { // Проверка на то, что заказ из базы должен быть выдан
			curTime := time.Now()
			orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
			if order.State != "accepted" { // Проверка на то, что заказ с таким OrderId был принят
				//fmt.Printf("Error: order with id %d can't be taken, because it has been already taken or still didn't come\n", order.Id)
				return errors.New("order with id  can't be taken, because it has been already taken or still didn't come")
			}
			if curTime.After(orderTime) { // Проверка на время
				//fmt.Printf("Error: order with id %d can't be taken, because time left\n", order.Id)
				return errors.New("order with id can't be taken, because time left")
			} else if userId == 0 { // Проверка на то, что все заказы принадлежат одному UserId
				userId = order.UserId
				tmpTime := time.Now().Add(time.Hour * 72).Format(TIMELAYOUT)
				ord.Orders[k].State = "gived"     // Помечаем как выданный
				ord.Orders[k].ValidTime = tmpTime // Ставим время до которого заказ можно вернуть
				cnt++
			} else if userId != order.UserId {
				//fmt.Println("Error: one or more orders are not yours")
				return errors.New("one or more orders are not yours")
			} else {
				tmpTime := time.Now().Add(time.Hour * 72).Format(TIMELAYOUT)
				ord.Orders[k].State = "gived"
				ord.Orders[k].ValidTime = tmpTime
				cnt++
			}
		}
	}

	if cnt != len(req.OrderIds) {
		//fmt.Println("Error: one or more orders are not in our pick-point")
		return errors.New("one or more orders are not in our pick-point")
	}
	oc.repo.InsertOrders(ord)
	return nil
}

func (oc *OrderUseCase) Return(req *dto.ReturnOrderRequest) error {
	ord, err := oc.repo.GetOrders()
	if err != nil {
		return err
	}

	// Обход нашей базы данных
	for k, order := range ord.Orders {
		if order.Id == req.Id {
			curTime := time.Now().Add(24 * time.Hour)
			orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
			if curTime.After(orderTime) || order.State == "returned" || order.State == "gived" { // Проверяем что срок хранения истек или заказ был возвращен
				if order.State != "gived" { // Проверяем что заказ не находится у клиента
					ord.Orders = append(ord.Orders[:k], ord.Orders[k+1:]...) //Удаляем из базы
					err := oc.repo.InsertOrders(ord)                         // Обновляем базу
					if err != nil {
						return err
					}
					return nil
				} else {
					return errors.New("this order is with the client")
				}
			} else {
				return errors.New("client still can take it")
			}
		}
	}
	return errors.New("no such order")
}

func (oc *OrderUseCase) UserOrders(req *dto.UserOrdersRequest) (*dto.UserOrdersResponse, error) {
	ord, err := oc.repo.GetOrders()
	if err != nil {
		return nil, err
	}

	// Поиск и сборка всех заказов от клиента с userId
	var userOrders []dto.OrderDto
	for _, order := range ord.Orders {
		if order.UserId == req.UserId && order.State != "gived" {
			//text := fmt.Sprintf("Order Id: %d, Valid untill: %s, State: %s", order.Id, order.ValidTime, order.State)
			userOrders = append(userOrders, order)
		}
	}
	if len(userOrders) == 0 {
		return nil, errors.New("empty")
	}

	if req.Last < 1 {
		return &dto.UserOrdersResponse{ListOrdersDto: dto.ListOrdersDto{Orders: userOrders}}, nil
	} else {
		if req.Last > len(userOrders) {
			req.Last = len(userOrders)
		}
		return &dto.UserOrdersResponse{ListOrdersDto: dto.ListOrdersDto{Orders: userOrders[:len(userOrders)-req.Last]}}, nil
	}
}

func (oc *OrderUseCase) UserReturns(req *dto.UserReturnsRequest) (*dto.UserReturnsResponse, error) {
	ord, err := oc.repo.GetOrders()
	if err != nil {
		return nil, err
	}

	var userReturns []dto.OrderDto
	for _, order := range ord.Orders {
		if order.State == "returned" {
			//text := fmt.Sprintf("Order Id: %d, User Id: %d", v.Id, v.UserId)
			userReturns = append(userReturns, order)
		}
	}
	if len(userReturns) == 0 {
		return nil, errors.New("empty")
	}
	totalPages := len(userReturns) / req.Size
	if len(userReturns)%req.Size != 0 {
		totalPages += 1
	}

	if req.Page > totalPages {
		return nil, errors.New("empty")
	}

	if req.Page*req.Size >= len(userReturns) {
		return &dto.UserReturnsResponse{ListOrdersDto: dto.ListOrdersDto{Orders: userReturns[(req.Page-1)*req.Size:]}}, nil
	} else {
		return &dto.UserReturnsResponse{ListOrdersDto: dto.ListOrdersDto{Orders: userReturns[(req.Page-1)*req.Size : req.Page*req.Size]}}, nil
	}
}
