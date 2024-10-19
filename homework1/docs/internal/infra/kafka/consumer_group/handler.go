package consumer_group

import (
	"encoding/json"
	"homework1/internal/dto"
	"log"

	"github.com/IBM/sarama"
)

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

type ConsumerGroupHandler struct {
	ready chan bool
}

func NewConsumerGroupHandler(
// map[topic]TopicHandler
) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{}
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор, пока сессия не завершилась
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msg := convertMsg(message)
			data, _ := json.Marshal(msg)
			log.Printf("Message claimed: key=%s, offset=%d, part=%d, topic=%s\n", msg.Key, msg.Offset, msg.Partition, msg.Topic)
			var event dto.EventDto
			err := json.Unmarshal([]byte(msg.Payload), &event)
			if err != nil {
				log.Printf("error: %s in decoding data payload: %s\n", err.Error(), data)
			} else {
				log.Printf("Parsed payload of income message: order id= %d, method=%s, time=%s\n", event.OrderId, event.Method, event.Time)
			}

			// mark message as successfully handled and ready to commit offset
			// autocommit may commit message offset sometime
			session.MarkMessage(message, "")

			// commit offset manually right now
			// works when autocommit disabled
			session.Commit()

			// autocommit not to work if cg not gracefully shut downed
			// but manual commit does
			// panic("emulate not commit")
		case <-session.Context().Done():
			return nil
		}
	}
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
