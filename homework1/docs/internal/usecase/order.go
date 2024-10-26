package usecase

import (
	"context"
	"errors"
	"fmt"
	"homework1/internal/domain"
	"homework1/internal/domain/strategy"
	"homework1/internal/dto"
	"homework1/internal/infra/kafka/producer"
	"homework1/internal/metrics"
	"homework1/internal/repository"
	"homework1/internal/repository/postgres"
	cliserver "homework1/pkg/cli/v1"
	"log"
	"time"

	"github.com/IBM/sarama"
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
	kafkaProducer sarama.SyncProducer
}

func NewOrderUseCase(repo orderRepository, psqlRepoFacade repository.Facade, kafkaProducer sarama.SyncProducer) *OrderUseCase {
	return &OrderUseCase{repo: repo, psqlRepoFacade: psqlRepoFacade, kafkaProducer: kafkaProducer}
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

func (oc *OrderUseCase) AcceptOrderGrpc(ctx context.Context, req *cliserver.AcceptOrderRequest) (*cliserver.AcceptOrderResponse, error) {
	if err := req.Validate(); err != nil {
		metrics.IncBadRespByHandler("acceptOrder", 3)
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
		metrics.IncBadRespByHandler("acceptOrder", 3)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unknown box type: %s", req.PackageType))
	}

	newOrder, err := domain.NewOrder(int(req.GetId()), int(req.GetUserId()), int(req.GetPrice()), int(req.GetWeight()), req.GetValidTime(), "accepted", req.GetPackageType(), opackageStrategy, req.GetAdditionalStretch())
	if err != nil {
		metrics.IncBadRespByHandler("acceptOrder", 3)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = oc.psqlRepoFacade.AddOrder(ctx, newOrder.ToDTO())
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrAlreadyInBase):
			metrics.IncBadRespByHandler("acceptOrder", 3)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			metrics.IncBadRespByHandler("acceptOrder", 13)
			log.Printf("UNKNOWN ERROR IN ACCEPTING ORDER: %s\n", err.Error())
			return nil, status.Error(codes.Internal, "Unkown Error")
		}
	}

	msg := producer.CreateMessage(int(req.GetId()), "AcceptOrder")
	p, o, err := producer.SendMessage(oc.kafkaProducer, int(req.GetUserId()), msg, "pvz.events-log")
	if err != nil {
		log.Println("[SEND MESSAGE TO KAFKA]: ", err)
	}
	fmt.Println(p, o)

	metrics.AddAcceptedOrder()
	metrics.IncOkRespByHandler("acceptOrder")

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

func (oc *OrderUseCase) AcceptReturnGrpc(ctx context.Context, req *cliserver.AcceptReturnRequest) (*cliserver.AcceptReturnResponse, error) {
	if err := req.Validate(); err != nil {
		metrics.IncBadRespByHandler("acceptReturn", 3)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	order, err := oc.psqlRepoFacade.GetOrderById(ctx, int(req.GetId()))
	if err != nil {
		log.Printf("UNKNOWN ERROR IN ACCEPT RETURN %s\n", err.Error())
		metrics.IncBadRespByHandler("acceptReturn", 13)
		return nil, status.Error(codes.Internal, "Unkown Error")
	}

	if order.UserId != int(req.GetUserId()) {
		metrics.IncBadRespByHandler("acceptReturn", 3)
		return nil, status.Error(codes.InvalidArgument, "Not your order")
	}

	if order.State != "gived" {
		metrics.IncBadRespByHandler("acceptReturn", 3)
		return nil, status.Error(codes.InvalidArgument, "Your order already returned or still not gived")
	}

	orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
	curTime := time.Now()
	if curTime.After(orderTime) {
		metrics.IncBadRespByHandler("acceptReturn", 3)
		return nil, status.Error(codes.InvalidArgument, "No time to return")
	}

	order.State = "returned"
	err = oc.psqlRepoFacade.UpdateOrderInfo(ctx, order)
	if err != nil {
		log.Printf("UNKOWN ERROR IN ACCEPT RETURN %s\n", err.Error())
		metrics.IncBadRespByHandler("acceptReturn", 13)
		return nil, status.Error(codes.Internal, "Unkown Error")
	}

	msg := producer.CreateMessage(int(req.GetId()), "AcceptReturn")
	p, o, err := producer.SendMessage(oc.kafkaProducer, int(req.GetUserId()), msg, "pvz.events-log")
	if err != nil {
		log.Println("[SEND MESSAGE TO KAFKA]: ", err)
	}
	fmt.Println(p, o)

	metrics.IncOkRespByHandler("acceptReturn")
	return &cliserver.AcceptReturnResponse{}, nil
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

func (oc *OrderUseCase) GiveOrderGrpc(ctx context.Context, req *cliserver.GiveOrderRequest) (*cliserver.GiveOrderResponse, error) {
	if err := req.Validate(); err != nil {
		metrics.IncBadRespByHandler("giveOrder", 3)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var AllOrders []dto.OrderDto
	var uniqueUserIds int
	for _, id := range req.GetOrderIds() {
		order, err := oc.psqlRepoFacade.GetOrderById(ctx, int(id.GetId()))
		if errors.Is(err, postgres.ErrNoSuchOrderd) {
			metrics.IncBadRespByHandler("giveOrder", 3)
			return nil, status.Error(codes.InvalidArgument, "no such order with")
		}
		if err != nil {
			log.Printf("UNKNOWN ERROR IN GIVING OREDER: %s", err.Error())
			metrics.IncBadRespByHandler("giveOrder", 13)
			return nil, status.Error(codes.Internal, "unkown error")
		}

		if uniqueUserIds == 0 {
			uniqueUserIds = order.UserId
		} else if uniqueUserIds != order.UserId {
			metrics.IncBadRespByHandler("giveOrder", 3)
			return nil, status.Error(codes.InvalidArgument, "one of orders is not yours")
		}

		if order.State != "accepted" {
			metrics.IncBadRespByHandler("giveOrder", 3)
			return nil, status.Error(codes.InvalidArgument, "order with id  can't be taken, because it has been already taken or still didn't come")
		}

		curTime := time.Now()
		orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
		if curTime.After(orderTime) {
			metrics.IncBadRespByHandler("giveOrder", 3)
			return nil, status.Error(codes.InvalidArgument, "order with id can't be taken, because time left")
		}
		AllOrders = append(AllOrders, order)
	}

	for _, order := range AllOrders {
		tmpTime := time.Now().Add(time.Hour * 72).Format(TIMELAYOUT)
		order.State = "gived"
		order.ValidTime = tmpTime

		err := oc.psqlRepoFacade.UpdateOrderInfo(ctx, order)
		if err != nil {
			log.Printf("UNKOWN ERROR IN GIVING OREDER: %s", err.Error())
			metrics.IncBadRespByHandler("giveOrder", 13)
			return nil, status.Error(codes.Internal, "unkown error")
		}
		msg := producer.CreateMessage(int(order.Id), "GiveOrder")
		p, o, err := producer.SendMessage(oc.kafkaProducer, int(order.UserId), msg, "pvz.events-log")
		if err != nil {
			log.Println("[SEND MESSAGE TO KAFKA]: ", err)
		}
		fmt.Println(p, o)
	}
	metrics.IncOkRespByHandler("giveOrder")
	return &cliserver.GiveOrderResponse{}, nil
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

func (oc *OrderUseCase) ReturnOrderGrpc(ctx context.Context, req *cliserver.ReturnOrderRequest) (*cliserver.ReturnOrderResponse, error) {
	if err := req.Validate(); err != nil {
		metrics.IncBadRespByHandler("returnOrder", 3)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	order, err := oc.psqlRepoFacade.GetOrderById(ctx, int(req.GetId()))
	if err != nil {
		log.Printf("UNKNOWN ERROR IN RETURNING OREDER: %s", err.Error())
		metrics.IncBadRespByHandler("returnOrder", 13)
		return nil, status.Error(codes.Internal, "unkown error")
	}

	if order.State == "gived" {
		metrics.IncBadRespByHandler("returnOrder", 3)
		return nil, status.Error(codes.InvalidArgument, "this order is with the client")
	}

	curTime := time.Now().Add(24 * time.Hour)
	orderTime, _ := time.Parse(TIMELAYOUT, order.ValidTime)
	if curTime.Before(orderTime) && order.State == "accepted" {
		metrics.IncBadRespByHandler("returnOrder", 3)
		return nil, status.Error(codes.InvalidArgument, "client still can take it")
	}

	order.State = "deleted"
	err = oc.psqlRepoFacade.UpdateOrderInfo(ctx, order)

	if err != nil {
		log.Printf("UNKNOWN ERROR IN RETURNING OREDER: %s", err.Error())
		metrics.IncBadRespByHandler("returnOrder", 13)
		return nil, status.Error(codes.Internal, "unkown error")
	}

	msg := producer.CreateMessage(int(order.Id), "ReturnOrder")
	p, o, err := producer.SendMessage(oc.kafkaProducer, int(order.UserId), msg, "pvz.events-log")
	if err != nil {
		log.Println("[SEND MESSAGE TO KAFKA]: ", err)
	}
	fmt.Println(p, o)

	metrics.IncOkRespByHandler("returnOrder")
	return &cliserver.ReturnOrderResponse{}, nil
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
		return &dto.UserOrdersResponse{ListOrdersDto: dto.ListOrdersDto{Orders: userOrders[len(userOrders)-req.Last:]}}, nil
	}
}

func (oc *OrderUseCase) UserOrdersGrpc(ctx context.Context, req *cliserver.UserOrdersRequest) (*cliserver.UserOrdersResponse, error) {
	if err := req.Validate(); err != nil {
		metrics.IncBadRespByHandler("userOrders", 3)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orders, err := oc.psqlRepoFacade.GetOrdersByUserId(ctx, int(req.GetUserId()))
	if err != nil {
		log.Printf("UNKNOWN ERROR IN USER ORDERS: %s", err.Error())
		metrics.IncBadRespByHandler("userOrders", 13)
		return nil, status.Error(codes.Internal, "unkown error")
	}
	var userOrders []dto.OrderDto
	for _, order := range orders.Orders {
		if order.State != "gived" {
			//text := fmt.Sprintf("Order Id: %d, Valid untill: %s, State: %s", order.Id, order.ValidTime, order.State)
			userOrders = append(userOrders, order)
		}
	}
	if len(userOrders) == 0 {
		metrics.IncOkRespByHandler("userOrders")
		return nil, nil
	}

	metrics.IncOkRespByHandler("userOrders")
	if req.GetLast() < 1 {
		resp := &cliserver.UserOrdersResponse{
			OrderDtos: make([]*cliserver.OrderDto, 0, len(userOrders)),
		}
		for _, order := range userOrders {
			resp.OrderDtos = append(resp.OrderDtos, &cliserver.OrderDto{
				Id:                int64(order.Id),
				UserId:            int64(order.UserId),
				ValidTime:         order.ValidTime,
				State:             order.State,
				Price:             int64(order.Price),
				Weight:            int64(order.Weight),
				PackageType:       order.PackageType,
				AdditionalStretch: order.AdditionalStretch,
			})
		}
		return resp, nil
	} else {
		lastcounter := int(req.GetLast())
		if int(req.GetLast()) > len(userOrders) {
			lastcounter = len(userOrders)
		}

		resp := &cliserver.UserOrdersResponse{
			OrderDtos: make([]*cliserver.OrderDto, 0, len(userOrders)-lastcounter),
		}
		for i := len(userOrders) - lastcounter; i < len(userOrders); i++ {
			resp.OrderDtos = append(resp.OrderDtos, &cliserver.OrderDto{
				Id:                int64(userOrders[i].Id),
				UserId:            int64(userOrders[i].UserId),
				ValidTime:         userOrders[i].ValidTime,
				State:             userOrders[i].State,
				Price:             int64(userOrders[i].Price),
				Weight:            int64(userOrders[i].Weight),
				PackageType:       userOrders[i].PackageType,
				AdditionalStretch: orders.Orders[i].AdditionalStretch,
			})
		}
		return resp, nil
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

func (oc *OrderUseCase) UserReturnsGrpc(ctx context.Context, req *cliserver.UserReturnsRequest) (*cliserver.UserReturnsResponse, error) {
	if err := req.Validate(); err != nil {
		metrics.IncBadRespByHandler("userReturns", 3)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orders, err := oc.psqlRepoFacade.GetUserReturns(ctx)
	if err != nil {
		log.Printf("UNKNOWN ERROR IN USER RETURNS: %s", err.Error())
		metrics.IncBadRespByHandler("userReturns", 13)
		return nil, status.Error(codes.Internal, "unkown error")
	}

	if len(orders.Orders) == 0 {
		metrics.IncOkRespByHandler("userReturns")
		return nil, nil
	}

	totalPages := len(orders.Orders) / int(req.GetSize())
	if len(orders.Orders)%int(req.GetSize()) != 0 {
		totalPages += 1
	}

	if int(req.GetPage()) > totalPages {
		metrics.IncOkRespByHandler("userReturns")
		return nil, nil
	}

	metrics.IncOkRespByHandler("userReturns")
	pagexsize := int(req.GetPage()) * int(req.GetSize())
	pagem1xsize := (int(req.GetPage()) - 1) * int(req.GetSize())
	if pagexsize >= len(orders.Orders) {
		resp := &cliserver.UserReturnsResponse{
			OrderDtos: make([]*cliserver.OrderDto, 0, len(orders.Orders)-pagem1xsize+1),
		}
		for i := pagem1xsize; i < len(orders.Orders); i++ {
			resp.OrderDtos = append(resp.OrderDtos, &cliserver.OrderDto{
				Id:                int64(orders.Orders[i].Id),
				UserId:            int64(orders.Orders[i].UserId),
				ValidTime:         orders.Orders[i].ValidTime,
				State:             orders.Orders[i].State,
				Price:             int64(orders.Orders[i].Price),
				Weight:            int64(orders.Orders[i].Weight),
				PackageType:       orders.Orders[i].PackageType,
				AdditionalStretch: orders.Orders[i].AdditionalStretch,
			})
		}
		return resp, nil
	} else {
		resp := &cliserver.UserReturnsResponse{
			OrderDtos: make([]*cliserver.OrderDto, 0, pagexsize-pagem1xsize+1),
		}
		for i := pagem1xsize; i < pagexsize; i++ {
			resp.OrderDtos = append(resp.OrderDtos, &cliserver.OrderDto{
				Id:                int64(orders.Orders[i].Id),
				UserId:            int64(orders.Orders[i].UserId),
				ValidTime:         orders.Orders[i].ValidTime,
				State:             orders.Orders[i].State,
				Price:             int64(orders.Orders[i].Price),
				Weight:            int64(orders.Orders[i].Weight),
				PackageType:       orders.Orders[i].PackageType,
				AdditionalStretch: orders.Orders[i].AdditionalStretch,
			})
		}
		return resp, nil
	}
}
