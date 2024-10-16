package usecase

import (
	"context"
	"errors"
	"fmt"
	"homework1/internal/domain"
	"homework1/internal/domain/strategy"
	"homework1/internal/dto"
	"homework1/internal/repository"
	"homework1/internal/repository/postgres"
	cliserver "homework1/pkg/cli/v1"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const TIMELAYOUT = "2006-01-02"

type (
	orderRepository interface {
		InsertOrders(data *dto.ListOrdersDto) error
		GetOrders() (*dto.ListOrdersDto, error)
	}
)

type OrderUseCase struct {
	repo           orderRepository
	psqlRepoFacade repository.Facade
	cliserver.UnimplementedCliServer
}

func NewOrderUseCase(repo orderRepository, psqlRepoFacade repository.Facade) *OrderUseCase {
	return &OrderUseCase{repo: repo, psqlRepoFacade: psqlRepoFacade}
}

func (oc *OrderUseCase) Accept(ctx context.Context, req *dto.AcceptOrderRequest) error {
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

	err = oc.psqlRepoFacade.AddOrder(ctx, newOrder.ToDTO())
	if err != nil {
		return err
	}

	return nil
}

func (oc *OrderUseCase) AcceptOrder(ctx context.Context, req *cliserver.AcceptOrderRequest) (*cliserver.AcceptOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown box type: %s", req.PackageType))
	}

	newOrder, err := domain.NewOrder(int(req.GetId()), int(req.GetUserId()), int(req.GetPrice()), int(req.GetWeight()), req.GetValidTime(), "accepted", req.GetPackageType(), opackageStrategy, req.GetAdditionalStretch())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = oc.psqlRepoFacade.AddOrder(ctx, newOrder.ToDTO())
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrAlreadyInBase):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			log.Println(err.Error())
			return nil, status.Error(codes.Internal, "Unkown Error")
		}
	}
	return &cliserver.AcceptOrderResponse{}, nil
}

func (oc *OrderUseCase) AcceptReturn(ctx context.Context, req *dto.AcceptReturnOrderRequest) error {
	order, err := oc.psqlRepoFacade.GetOrderById(ctx, req.Id)
	if err != nil {
		return err
	}

	if order.State != "gived" {
		return errors.New("your order already returned or still not gived")
	}

	orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
	curTime := time.Now()
	if curTime.After(orderTime) {
		return errors.New("no time to return")
	}

	order.State = "returned"
	err = oc.psqlRepoFacade.UpdateOrderInfo(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (oc *OrderUseCase) Give(ctx context.Context, req *dto.GiveOrderRequest) error {
	var AllOrders []dto.OrderDto
	var uniqueUserIds int
	for _, id := range req.OrderIds {
		order, err := oc.psqlRepoFacade.GetOrderById(ctx, id.Id)
		if err != nil {
			return err
		}

		if uniqueUserIds == 0 {
			uniqueUserIds = order.UserId
		} else if uniqueUserIds != order.UserId {
			return errors.New("one of orders is not yours")
		}

		if order.State != "accepted" {
			return errors.New("order with id  can't be taken, because it has been already taken or still didn't come")
		}

		curTime := time.Now()
		orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
		if curTime.After(orderTime) {
			return errors.New("order with id can't be taken, because time left")
		}
		AllOrders = append(AllOrders, order)
	}

	for _, order := range AllOrders {
		tmpTime := time.Now().Add(time.Hour * 72).Format(TIMELAYOUT)
		order.State = "gived"
		order.ValidTime = tmpTime

		err := oc.psqlRepoFacade.UpdateOrderInfo(ctx, order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (oc *OrderUseCase) Return(ctx context.Context, req *dto.ReturnOrderRequest) error {
	order, err := oc.psqlRepoFacade.GetOrderById(ctx, req.Id)
	if err != nil {
		return err
	}

	if order.State == "gived" {
		return errors.New("this order is with the client")
	}

	curTime := time.Now().Add(24 * time.Hour)
	orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
	if curTime.Before(orderTime) && order.State == "accepted" {
		return errors.New("client still can take it")
	}

	order.State = "deleted"
	err = oc.psqlRepoFacade.UpdateOrderInfo(ctx, order)

	if err != nil {
		return err
	}

	return nil
}

func (oc *OrderUseCase) UserOrders(ctx context.Context, req *dto.UserOrdersRequest) (*dto.UserOrdersResponse, error) {
	orders, err := oc.psqlRepoFacade.GetOrdersByUserId(ctx, req.UserId)
	if err != nil {
		return &dto.UserOrdersResponse{}, err
	}

	var userOrders []dto.OrderDto
	for _, order := range orders.Orders {
		if order.State != "gived" {
			//text := fmt.Sprintf("Order Id: %d, Valid untill: %s, State: %s", order.Id, order.ValidTime, order.State)
			userOrders = append(userOrders, order)
		}
	}
	if len(userOrders) == 0 {
		return nil, nil
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

func (oc *OrderUseCase) UserReturns(ctx context.Context, req *dto.UserReturnsRequest) (*dto.UserReturnsResponse, error) {
	orders, err := oc.psqlRepoFacade.GetUserReturns(ctx)
	if err != nil {
		return &dto.UserReturnsResponse{}, err
	}

	if len(orders.Orders) == 0 {
		return nil, errors.New("empty")
	}
	totalPages := len(orders.Orders) / req.Size
	if len(orders.Orders)%req.Size != 0 {
		totalPages += 1
	}

	if req.Page > totalPages {
		return nil, errors.New("empty")
	}

	if req.Page*req.Size >= len(orders.Orders) {
		return &dto.UserReturnsResponse{ListOrdersDto: dto.ListOrdersDto{Orders: orders.Orders[(req.Page-1)*req.Size:]}}, nil
	} else {
		return &dto.UserReturnsResponse{ListOrdersDto: dto.ListOrdersDto{Orders: orders.Orders[(req.Page-1)*req.Size : req.Page*req.Size]}}, nil
	}
}
