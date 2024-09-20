package main

import (
	"context"
	"os"
	"os/signal"
	"jade-mes/app"
	"jade-mes/app/infrastructure/trace"
	"jade-mes/config"

	"jade-mes/app/infrastructure/persistence/database"

	"jade-mes/app/infrastructure/log"
)

func main() {
	logger, accessLogger := log.GetLoggers()
	log.RedirectStdLog(logger)
	log.InitGlobalLogger(logger)
	log.InitAccessLogger(accessLogger)

	defer database.GetRedis().Close()
	defer database.GetDB().Close()

	//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//	defer cancel()
	//	defer database.GetMongoDB().Client().Disconnect(ctx)

	// rabbitMQProducer := mq.GetRabbitMQProducer()
	// rabbitMQConsumers := mq.GetRabbitMQConsumers()
	// defer rabbitMQProducer.Close()
	// for _, consumer := range rabbitMQConsumers {
	// 	defer consumer.Close()
	// }

	//app.Start()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	tp, cleanup, err := trace.NewTracerProvider(config.Config.Tracer.Endpoint, config.Config.Tracer.Name)
	if err != nil {
		log.Printf("tracerProvider init error: %s", err)
		return
	}
	defer cleanup(context.Background(), tp)

	go app.Start()
	select {
	case <-sigCh:
		log.Printf("goodbye")
		return
	}
}
