package producer

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMySendMessage_Success(t *testing.T) {
	mockProd := mocks.NewSyncProducer(t, nil)
	key := 123
	message := []byte("test message")
	topic := "test-topic"

	mockProd.ExpectSendMessageAndSucceed()

	partition, offset, err := SendMessage(mockProd, key, message, topic)

	assert.NoError(t, err)
	assert.Equal(t, int32(0), partition)
	assert.Equal(t, int64(1), offset)
}

func TestMySendMessage_Error(t *testing.T) {
	mockProd := mocks.NewSyncProducer(t, nil)
	key := 456
	message := []byte("another message")
	topic := "another-topic"

	mockProd.ExpectSendMessageAndFail(sarama.ErrNotLeaderForPartition)

	_, _, err := SendMessage(mockProd, key, message, topic)

	assert.Error(t, err)
	assert.Equal(t, sarama.ErrNotLeaderForPartition, err)

}
