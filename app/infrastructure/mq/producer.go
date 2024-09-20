package mq

import (
	"time"

	"jade-mes/app/infrastructure/log"
	"jade-mes/app/infrastructure/util"
)

// PublishMessage ...
func PublishMessage(data []byte, routingKey string) error {
	logger := log.With(log.String("rawData", string(data)), log.String("routingKey", routingKey))

	logger.Info("publishing message")

	defer util.TimeTrack(time.Now(), logger)

	return rabbitMQProducer.Push(data, routingKey)
}
