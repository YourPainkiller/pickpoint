package repository

import (
	"context"
	"homework1/internal/dto"
	"homework1/internal/repository/postgres"
)

type Facade interface {
	AddOrder(ctx context.Context, req dto.OrderDto) error
	GetOrderById(ctx context.Context, id int) (dto.OrderDto, error)
	UpdateOrderInfo(ctx context.Context, req dto.OrderDto) error
	GetOrdersByUserId(ctx context.Context, userId int) (dto.UserOrdersResponse, error)
	GetUserReturns(ctx context.Context) (dto.UserReturnsResponse, error)
	DropTable(ctx context.Context) error
}

type storageFacade struct {
	txManager    postgres.TransactionManager
	pgRepository postgres.PgRepository
}

func NewStorageFacade(pgRepository postgres.PgRepository, txManager postgres.TransactionManager) *storageFacade {
	return &storageFacade{pgRepository: pgRepository, txManager: txManager}
}

func (s *storageFacade) AddOrder(ctx context.Context, req dto.OrderDto) error {
	return s.txManager.RunReadWriteCommited(ctx, func(ctxTx context.Context) error {
		err := s.pgRepository.AddOrder(ctxTx, req)
		if err != nil {
			return err
		}
		return nil
	})

}

func (s *storageFacade) GetOrderById(ctx context.Context, id int) (dto.OrderDto, error) {
	order, err := s.pgRepository.GetOrderById(ctx, id)
	if err != nil {
		return dto.OrderDto{}, err
	}
	return order, nil
}

func (s *storageFacade) UpdateOrderInfo(ctx context.Context, req dto.OrderDto) error {
	err := s.pgRepository.UpdateOrderInfo(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *storageFacade) GetOrdersByUserId(ctx context.Context, userId int) (dto.UserOrdersResponse, error) {
	orders, err := s.pgRepository.GetOrdersByUserId(ctx, userId)
	if err != nil {
		return dto.UserOrdersResponse{}, err
	}

	return orders, nil
}

func (s *storageFacade) GetUserReturns(ctx context.Context) (dto.UserReturnsResponse, error) {
	orders, err := s.pgRepository.GetUserReturns(ctx)
	if err != nil {
		return dto.UserReturnsResponse{}, err
	}

	return orders, nil
}

func (s *storageFacade) DropTable(ctx context.Context) error {
	err := s.pgRepository.DropTable(ctx)
	if err != nil {
		return err
	}
	return nil
}
