package mq

import (
	"jade-mes/config"
)

var rabbitMQProducer *RabbitMQSession
var rabbitMQConsumer *RabbitMQSession

func init() {
	println("initing mq...")
	settings := config.GetConfig()

	rabbitMQProducer = initRabbitMQSession(settings, false, false)
	go rabbitMQProducer.handleReconnect(rabbitMQProducer.addr)

	rabbitMQConsumer = initRabbitMQSession(settings, true, false)
	Consume()
	go rabbitMQConsumer.handleReconnect(rabbitMQConsumer.addr)
}
