package sender

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/move-mates/trinquet/library/kafka/pkg/config"
	"github.com/move-mates/trinquet/library/kafka/pkg/event"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EventSender interface {
	SendEvent(event event.Wrapper) error
}

type eventSender struct {
	producer *kafka.Producer
	topic    string
}

func NewEventSender(
	config config.EventSenderConfig,
	logger *logrus.Logger,
	senderName string,
) (EventSender, error) {
	producer, err := kafka.NewProducer(config.ToKafkaConfigMap())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	go listener(producer, logger, senderName)

	return &eventSender{
		producer: producer,
		topic:    config.Topic,
	}, nil
}

func (sender *eventSender) SendEvent(event event.Wrapper) error {
	value, err := json.Marshal(event)
	if err != nil {
		return errors.WithStack(err)
	}

	message := &kafka.Message{
		Key:   []byte(event.UserID.String()),
		Value: value,
		TopicPartition: kafka.TopicPartition{
			Topic:     &sender.topic,
			Partition: kafka.PartitionAny,
		},
		Timestamp: event.Timestamp,
	}

	err = sender.producer.Produce(message, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
