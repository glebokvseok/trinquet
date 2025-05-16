package sender

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/move-mates/trinquet/library/common/pkg/extensions"
	"github.com/move-mates/trinquet/library/common/pkg/request"
	"github.com/move-mates/trinquet/library/kafka/pkg/event"
	"github.com/sirupsen/logrus"
)

func listener(producer *kafka.Producer, logger *logrus.Logger, senderName string) {
	logger.Infof("%s sender event listener started", senderName)

	for ok := true; ok; {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("panic occured in %s sender event listener: %+v", senderName, r)
				}
			}()

			var evRaw kafka.Event
			evRaw, ok = <-producer.Events()

			switch ev := evRaw.(type) {
			case *kafka.Message:
				var wrapper event.Wrapper
				err := json.Unmarshal(ev.Value, &wrapper)
				if err != nil {
					logger.Warnf("unknown event format at %s sender event listener: \n%s", senderName, string(ev.Value))

					break
				}

				ctx := request.CreateContext(context.Background(), wrapper.RequestID)

				if ev.TopicPartition.Error != nil {
					logger.
						WithContext(ctx).
						Errorf("event delivery failed: %+v", ev.TopicPartition.Error)
				} else {
					logger.
						WithContext(ctx).
						Infof("event successfully delivered to %s topic [partition: %d, offset: %d]",
							extensions.SafeUnwrap(ev.TopicPartition.Topic),
							ev.TopicPartition.Partition,
							ev.TopicPartition.Offset,
						)
				}
			case *kafka.Error:
				logger.Errorf("error occured in %s sender event listener: %+v", senderName, ev)
			}
		}()
	}

	logger.Warnf("%s sender event listener stopped", senderName)
}
