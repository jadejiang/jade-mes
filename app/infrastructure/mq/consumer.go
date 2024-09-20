package mq

import (
	"sync"

	"github.com/IBM/sarama"
	"github.com/streadway/amqp"

	"jade-mes/app/infrastructure/constant"
	"jade-mes/app/infrastructure/log"
	"jade-mes/app/infrastructure/util"
)

var mutex sync.Mutex

var attempts = make(map[string]int, 100)

func dataReportHandler(delivery <-chan amqp.Delivery, err error) {
	queueName := ""
	logger := log.With(log.String("queueName", queueName))

	if err != nil {
		logger.Error("Error in rabbitmq consumer", log.Err(err))
	}

	// limit := make(chan struct{}, 5000)

	// for {
	// 	select {
	// 	// if not ok, it probably means mq is down
	// 	case message, ok := <-delivery:
	// 		if ok {
	// 			limit <- struct{}{}
	// 			go func() {
	// 				_ = consumeWithRetry(&message, handler.HandleDeviceDataReport)
	// 				<-limit
	// 			}()
	// 		}
	// 	}
	// }
}

func consumeHandler(msg *sarama.ConsumerMessage) {
	key := string(msg.Topic)

	spanID, _ := util.NewUUID()
	log.Debug("message incoming", log.String("spanId", spanID), log.Reflect("message_info", msg))

	switch key {
	default:
		log.Warn("message missed", log.String("spanId", spanID), log.String("message_key", key))
	}
}

func consumeWithRetry(message *amqp.Delivery, handler func(rawBody []byte) error) error {
	log.Info("incoming mq message", log.Reflect("message", message))

	if err := handler(message.Body); err != nil {
		mutex.Lock()
		if attempts[message.MessageId] < constant.MaxAttempts {
			attempts[message.MessageId]++
			mutex.Unlock()

			return message.Nack(false, true)
		}

		if _, exists := attempts[message.MessageId]; exists {
			delete(attempts, message.MessageId)
		}
		mutex.Unlock()

		log.Error("max attemps exceeds while handling mq", log.Reflect("message", message), log.Err(err))
		return message.Nack(false, false)
	}

	mutex.Lock()
	if _, exists := attempts[message.MessageId]; exists {
		delete(attempts, message.MessageId)
	}
	mutex.Unlock()

	return message.Ack(false)
}

// Consume ...
func Consume() {
	// dataReportConsumer.RegisterConsumeHandlers(map[string]func(<-chan amqp.Delivery, error){})
}
