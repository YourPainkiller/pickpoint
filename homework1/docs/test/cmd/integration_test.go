package test

import (
	"context"
	"homework1/internal/cli"
	"homework1/internal/dto"
	"homework1/internal/repository"
	"homework1/internal/repository/postgres"
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestDb(t *testing.T) {
	const psqlDSN = "postgres://postgres:qwe@localhost:5433/postgresFake?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	require.NoError(t, err)
	defer pool.Close()

	txManager := postgres.NewTxManager(pool)
	pgRepo := postgres.NewPgRepository(txManager)

	storageFacade := repository.NewStorageFacade(*pgRepo, txManager)
	err = storageFacade.DropTable(ctx)
	require.NoError(t, err)

	t.Run("test inserting in table", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			var order dto.OrderDto
			err := gofakeit.Struct(&order)
			require.NoError(t, err)

			order.Id = i + 1
			if order.PackageType != "stretch" {
				x := rand.Intn(2)
				if x == 0 {
					order.AdditionalStretch = true
				} else {
					order.AdditionalStretch = false
				}
			}
			randHours := rand.Intn(201) - 100
			order.ValidTime = time.Now().Add(time.Hour * time.Duration(randHours)).Format(cli.TIMELAYOUT)
			err = storageFacade.AddOrder(ctx, order)
			require.NoError(t, err)
		}
	})

	t.Run("test get data", func(t *testing.T) {
		ctx := context.Background()
		var order dto.OrderDto
		err := gofakeit.Struct(&order)
		require.NoError(t, err)
		order.Id = 777
		storageFacade.AddOrder(ctx, order)

		orderFromBase, err := storageFacade.GetOrderById(ctx, 777)
		require.NoError(t, err)

		require.EqualValues(t, order.UserId, orderFromBase.UserId)
	})
}
