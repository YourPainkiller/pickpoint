package test

import (
	"context"
	"encoding/json"
	"fmt"
	"homework1/internal/cli"
	"homework1/internal/dto"
	"homework1/internal/imcache"
	"homework1/internal/infra/kafka"
	"homework1/internal/infra/kafka/consumer"
	"homework1/internal/infra/kafka/producer"
	"homework1/internal/repository"
	"homework1/internal/repository/postgres"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
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
	ttlOrdersCache := imcache.NewOrdersCache(60*time.Second, 1000)

	storageFacade := repository.NewStorageFacade(*pgRepo, txManager, ttlOrdersCache)
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

func TestKafka(t *testing.T) {
	kafkaConfig := newConfig()
	prod, err := producer.NewSyncProducer(kafkaConfig,
		producer.WithIdempotent(),
		producer.WithRequiredAcks(sarama.WaitForAll),
		producer.WithMaxOpenRequests(1),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
		// producer.WithProducerPartitioner(sarama.NewManualPartitioner),
		// producer.WithProducerPartitioner(sarama.NewRoundRobinPartitioner),
		// producer.WithProducerPartitioner(sarama.NewRandomPartitioner),
		producer.WithProducerPartitioner(sarama.NewHashPartitioner), // default
	)
	assert.NoError(t, err)

	var (
		wg        = &sync.WaitGroup{}
		conf      = newConfig()
		ctx, cncl = runSignalHandler(context.Background(), wg)
	)

	cons, err := consumer.NewConsumer(conf,
		consumer.WithReturnErrorsEnabled(true),
	)
	assert.NoError(t, err)

	defer cons.Close() // Не забываем освобождать ресурсы :)
	dn := make(chan bool)
	err = cons.ConsumeTopic(ctx, "pvz.events-log", func(msg *sarama.ConsumerMessage) {
		message := convertMsg(msg)
		var event dto.EventDto
		json.Unmarshal([]byte(message.Payload), &event)
		assert.Equal(t, "AcceptOrder", event.Method)
		assert.EqualValues(t, 1, event.OrderId)
		dn <- true
	}, wg)
	msg := producer.CreateMessage(1, "AcceptOrder")
	p, o, err2 := producer.SendMessage(prod, 0, msg, "pvz.events-log")
	<-dn
	cncl()
	cons.Close()
	assert.NoError(t, err2)
	assert.EqualValues(t, 0, p)
	assert.EqualValues(t, 0, o)
	//ctx.Done()
	assert.NoError(t, err)
	wg.Wait()
}

func newConfig() kafka.Config {
	return kafka.Config{
		Brokers: []string{
			"localhost:9092",
		},
	}
}

func runSignalHandler(ctx context.Context, wg *sync.WaitGroup) (context.Context, context.CancelFunc) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	sigCtx, cancel := context.WithCancel(ctx)

	wg.Add(1)
	go func() {
		defer signal.Stop(sigterm)
		defer wg.Done()
		defer cancel()

		for {
			select {
			case sig, ok := <-sigterm:
				if !ok {
					fmt.Printf("[signal] signal chan closed: %s\n", sig.String())
					return
				}

				fmt.Printf("[signal] signal recv: %s\n", sig.String())
				return
			case _, ok := <-sigCtx.Done():
				if !ok {
					fmt.Println("[signal] context closed")
					return
				}

				fmt.Printf("[signal] ctx done: %s\n", ctx.Err().Error())
				return
			}
		}
	}()

	return sigCtx, cancel
}

type Msg struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Key       string `json:"key"`
	Payload   string `json:"payload"`
}

func convertMsg(in *sarama.ConsumerMessage) Msg {
	return Msg{
		Topic:     in.Topic,
		Partition: in.Partition,
		Offset:    in.Offset,
		Key:       string(in.Key),
		Payload:   string(in.Value),
	}
}
