package repository

import (
	"context"
	"homework1/internal/dto"
	"homework1/internal/imcache"
	"homework1/internal/repository/postgres"
	"time"
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
	cache        *imcache.OrdersCache
}

func NewStorageFacade(pgRepository postgres.PgRepository, txManager postgres.TransactionManager, cache *imcache.OrdersCache) *storageFacade {
	return &storageFacade{pgRepository: pgRepository, txManager: txManager, cache: cache}
}

func (s *storageFacade) AddOrder(ctx context.Context, req dto.OrderDto) error {
	if _, ok := s.cache.Get(ctx, req.Id); ok {
		return postgres.ErrAlreadyInBase
	}

	return s.txManager.RunReadWriteCommited(ctx, func(ctxTx context.Context) error {
		err := s.pgRepository.AddOrder(ctxTx, req)
		if err != nil {
			return err
		}
		s.cache.Set(ctxTx, req.Id, &req, time.Now())
		return nil
	})
}

func (s *storageFacade) GetOrderById(ctx context.Context, id int) (dto.OrderDto, error) {
	if order, ok := s.cache.Get(ctx, id); ok {
		return *order, nil
	}

	order, err := s.pgRepository.GetOrderById(ctx, id)
	if err != nil {
		return dto.OrderDto{}, err
	}
	s.cache.Set(ctx, id, &order, time.Now())
	return order, nil
}

func (s *storageFacade) UpdateOrderInfo(ctx context.Context, req dto.OrderDto) error {
	err := s.pgRepository.UpdateOrderInfo(ctx, req)
	if err != nil {
		return err
	}

	s.cache.Set(ctx, req.Id, &req, time.Now())
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
