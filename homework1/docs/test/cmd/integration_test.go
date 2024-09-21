package test

import (
	"homework1/internal/dto"
	"homework1/internal/repository"
	"homework1/internal/usecase"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepoAdd(t *testing.T) {
	testRepo, err := repository.NewOrderRepository("./testdb.json")
	require.NoError(t, err)

	oc := usecase.NewOrderUseCase(testRepo)
	err = oc.Accept(&dto.AcceptOrderRequest{
		Id:                rand.Int(),
		UserId:            777,
		ValidTime:         "2025-09-21",
		Price:             100,
		Weight:            100,
		PackageType:       "box",
		AdditionalStretch: false,
	})
	require.NoError(t, err)
}

func TestRepoUserOrders(t *testing.T) {
	testRepo, err := repository.NewOrderRepository("./testdb.json")
	require.NoError(t, err)

	oc := usecase.NewOrderUseCase(testRepo)

	_, err = oc.UserOrders(&dto.UserOrdersRequest{UserId: 777})
	require.NoError(t, err)
}
