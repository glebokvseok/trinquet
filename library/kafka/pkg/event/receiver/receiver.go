package receiver

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/move-mates/trinquet/library/common/pkg/extensions"
	"github.com/move-mates/trinquet/library/common/pkg/request"
	"github.com/move-mates/trinquet/library/kafka/pkg/config"
	"github.com/move-mates/trinquet/library/kafka/pkg/event"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type EventReceiver interface {
	ReceiveEvent(timeout time.Duration) (event.Wrapper, error)
	Name() string
}

type eventReceiver struct {
	consumer *kafka.Consumer
	logger   *logrus.Logger
	name     string
}

func NewEventReceiver(
	config config.EventReceiverConfig,
	logger *logrus.Logger,
	name string,
) (EventReceiver, error) {
	configMap := config.ToKafkaConfigMap()
	err := configMap.SetKey("group.id", config.ConsumerGroup)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = consumer.Subscribe(config.Topic, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &eventReceiver{
		consumer: consumer,
		logger:   logger,
		name:     name,
	}, nil
}

func (receiver *eventReceiver) ReceiveEvent(timeout time.Duration) (event.Wrapper, error) {
	message, err := receiver.consumer.ReadMessage(timeout)
	if err != nil {
		return event.Wrapper{}, errors.WithStack(err)
	}

	var wrapper event.Wrapper
	err = json.Unmarshal(message.Value, &wrapper)
	if err != nil {
		receiver.logger.Warnf("unknown event format at %s event receiver: \n%s", receiver.name, string(message.Value))

		return event.Wrapper{}, errors.WithStack(err)
	}

	receiver.logger.
		WithContext(request.CreateContext(context.Background(), wrapper.RequestID)).
		Infof("event successfully received from %s topic [partition: %d, offset: %d]",
			extensions.SafeUnwrap(message.TopicPartition.Topic),
			message.TopicPartition.Partition,
			message.TopicPartition.Offset,
		)

	return wrapper, nil
}

func (receiver *eventReceiver) Name() string {
	return receiver.name
}
