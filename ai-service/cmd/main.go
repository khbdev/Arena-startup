package main

import (
	"ai-service/internal/config"
	// "ai-service/test"
	"ai-service/internal/rabbitmq"
)


func main(){
config.Init()
config.InitRedis()
rabbitmq.StartConsumer(config.NewRabbitMQ().Channel)
// test.PublishTestMessages()
}


