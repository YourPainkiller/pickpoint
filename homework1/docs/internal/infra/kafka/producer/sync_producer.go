package producer

import (
	"fmt"
	"homework1/internal/infra/kafka"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

func NewSyncProducer(conf kafka.Config, opts ...Option) (sarama.SyncProducer, error) {
	config := PrepareConfig(opts...)

	syncProducer, err := sarama.NewSyncProducer(conf.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	return syncProducer, nil
}

func SendMessage(prod sarama.SyncProducer, key int, message []byte, topic string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(int64(key), 10)), // comment for random
		Value: sarama.ByteEncoder(message),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("app-name"),
				Value: []byte("route256-sync-prod"),
			},
		},
		// Partition: 0,
		Timestamp: time.Now(),
	}
	partition, offset, err = prod.SendMessage(msg)
	return partition, offset, err
}
