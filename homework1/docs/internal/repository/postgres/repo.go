package postgres

import (
	"context"
	"errors"
	"fmt"
	"homework1/internal/dto"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const sameIdErrorCode = "23505"

type PgRepository struct {
	txManager TransactionManager
}

func NewPgRepository(txManager TransactionManager) *PgRepository {
	return &PgRepository{txManager: txManager}
}

func (r *PgRepository) AddOrder(ctx context.Context, req dto.OrderDto) error {
	tx := r.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
	insert into orders(order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch) values($1, $2, $3, $4, $5, $6, $7, $8)
	`, req.Id, req.UserId, req.ValidTime, req.State, req.Price, req.Weight, req.PackageType, req.AdditionalStretch)

	if unwrapPgCode(err) == sameIdErrorCode {
		return errors.New("this order already in base")
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *PgRepository) GetOrderById(ctx context.Context, id int) (dto.OrderDto, error) {
	tx := r.txManager.GetQueryEngine(ctx)
	var order dto.OrderDto

	err := tx.QueryRow(ctx, `
	select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where order_id = $1
	`, id).Scan(&order.Id, &order.UserId, &order.ValidTime, &order.State, &order.Price, &order.Weight, &order.PackageType, &order.AdditionalStretch)

	if errors.Is(err, pgx.ErrNoRows) {

		return dto.OrderDto{}, fmt.Errorf("no such order with id=%d", id)
	}
	if err != nil {
		return dto.OrderDto{}, err
	}

	return order, nil
}

func (r *PgRepository) UpdateOrderInfo(ctx context.Context, req dto.OrderDto) error {
	tx := r.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
	update orders set valid_time = $1, order_state = $2 where order_id = $3
	`, req.ValidTime, req.State, req.Id)

	if err != nil {
		return err
	}

	return nil
}

func (r *PgRepository) GetOrdersByUserId(ctx context.Context, userId int) (dto.UserOrdersResponse, error) {
	var orders []dto.OrderDto
	tx := r.txManager.GetQueryEngine(ctx)

	err := pgxscan.Select(ctx, tx, &orders, `
	select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where user_id = $1 and order_state != 'deleted'
	`, userId)
	if err != nil {
		return dto.UserOrdersResponse{}, err
	}
	return dto.UserOrdersResponse{ListOrdersDto: dto.ListOrdersDto{Orders: orders}}, nil
}

func (r *PgRepository) GetUserReturns(ctx context.Context) (dto.UserReturnsResponse, error) {
	var orders []dto.OrderDto
	tx := r.txManager.GetQueryEngine(ctx)

	err := pgxscan.Select(ctx, tx, &orders, `
	select order_id, user_id, valid_time, order_state, price, weight, package, additional_stretch from orders where order_state = 'returned'
	`)
	if err != nil {
		return dto.UserReturnsResponse{}, err
	}
	return dto.UserReturnsResponse{ListOrdersDto: dto.ListOrdersDto{Orders: orders}}, nil
}

func (r *PgRepository) DropTable(ctx context.Context) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `
	truncate table orders`)

	if err != nil {
		return err
	}
	return nil
}

func unwrapPgCode(err error) string {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgErr.Code
		}
	}
	return ""
}
