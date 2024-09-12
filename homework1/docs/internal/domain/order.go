package domain

import (
	"errors"
	"fmt"
	"homework1/internal/domain/strategy"
	"homework1/internal/dto"
	"time"
)

const TIMELAYOUT = "2006-01-02"

// Пакет для работы с хранилищем: чтением и записью в него json

var (
	ErrInvalidOrderId = errors.New("invalid orderId")
)

const (
	TypeUnknown string = "unknown"
	TypeBag     string = "bag"
	TypeBox     string = "box"
	TypeStretch string = "stretch"
)

type Order struct {
	id                int
	userId            int
	validTime         string
	state             string
	price             int
	weight            int
	orderPackage      string
	additionalStretch bool
}

func NewOrder(id, userid, price, weight int, time, state string, op string, ps strategy.OrderPackageStrategy, astretch bool) (*Order, error) {
	order := Order{
		orderPackage: op,
	}

	err := order.SetId(id)
	if err != nil {
		return nil, fmt.Errorf("NewOrder.SetId: %w", err)
	}

	err = order.SetUserId(userid)
	if err != nil {
		return nil, fmt.Errorf("NewOrder.SetUserId: %w", err)
	}

	err = order.SetTime(time)
	if err != nil {
		return nil, fmt.Errorf("NewOrder.SetTime: %w", err)
	}

	order.SetState(state)

	order.price, order.weight, err = setPackage(price, weight, astretch, ps)
	if err != nil {
		return nil, err
	}
	order.additionalStretch = astretch

	return &order, nil
}

func (o *Order) SetId(id int) error {
	if id < 1 {
		return ErrInvalidOrderId
	}
	o.id = id
	return nil
}

func (o *Order) SetUserId(userid int) error {
	if userid < 1 {
		return ErrInvalidOrderId
	}
	o.userId = userid
	return nil
}

func (o *Order) SetTime(vt string) error {
	validTime, err := time.Parse(TIMELAYOUT, vt)
	if err != nil {
		return err
	}
	validTime = validTime.Add(24 * time.Hour)

	curTime := time.Now()
	if validTime.Before(curTime) {
		return errors.New("invalid time")
	}

	o.validTime = validTime.Format(TIMELAYOUT)
	return nil
}

func (o *Order) SetState(state string) {
	o.state = state
}

func setPackage(price, weight int, additionalStretch bool, ps strategy.OrderPackageStrategy) (int, int, error) {
	return ps.ChoosePackage(price, weight, additionalStretch)
}

func (o *Order) ToDTO() dto.OrderDto {
	return dto.OrderDto{
		Id:                o.id,
		UserId:            o.userId,
		ValidTime:         o.validTime,
		State:             o.state,
		Price:             o.price,
		Weight:            o.weight,
		PackageType:       o.orderPackage,
		AdditionalStretch: o.additionalStretch,
	}
}
