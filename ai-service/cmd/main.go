package main

import (
	"ai-service/internal/config"
	"ai-service/internal/rabbitmq"
)


func main(){
config.Init()
rabbitmq.StartConsumer(config.NewRabbitMQ().Channel)
}


