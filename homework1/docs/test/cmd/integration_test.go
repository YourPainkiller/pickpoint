package test

import (
	"errors"
	"homework1/internal/dto"
	"homework1/internal/repository"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepoInsert(t *testing.T) {
	testRepo, err := repository.NewOrderRepository("./testdb.json")
	require.NoError(t, err)

	err = testRepo.InsertOrders(&dto.ListOrdersDto{
		Orders: []dto.OrderDto{
			{
				Id:                1,
				UserId:            1,
				ValidTime:         "2024-09-22",
				Price:             120,
				Weight:            100,
				PackageType:       "box",
				AdditionalStretch: false,
				State:             "accepted",
			},
		},
	})

	require.NoError(t, err)
}

func TestRepoGetEmpty(t *testing.T) {
	testRepo, err := repository.NewOrderRepository("./testdb.json")
	require.NoError(t, err)

	_, err = testRepo.GetOrders()
	require.NoError(t, err)
}

func TestRepoGetData(t *testing.T) {
	testData := &dto.ListOrdersDto{
		Orders: []dto.OrderDto{
			{
				Id:                1,
				UserId:            1,
				ValidTime:         "2024-09-22",
				Price:             120,
				Weight:            100,
				PackageType:       "box",
				AdditionalStretch: false,
				State:             "accepted",
			},
		},
	}

	testRepo, err := repository.NewOrderRepository("./testdb.json")
	require.NoError(t, err)

	err = testRepo.InsertOrders(testData)
	require.NoError(t, err)

	fromRepoData, err := testRepo.GetOrders()
	require.NoError(t, err)

	if !reflect.DeepEqual(testData, fromRepoData) {
		err = errors.New("bad reading")
	}
	require.NoError(t, err)
}
