package main

import (
	"ai-service/internal/config"
	"ai-service/internal/model"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func PublishTestMessages() {
	r := config.NewRabbitMQ() // RabbitMQ ulanish
	ch := r.Channel
	defer r.Conn.Close()
	defer ch.Close()

	queueName := "arena_queue" // Queue nomi
	exchange := "direct_exchange"

	for i := 1; i <= 3; i++ {
		msg := model.TestRequest{
			TelegramID: 3555345,
			Prompt:     fmt.Sprintf("Test prompt %d: Present Simple haqida", i),
			Count:      1,
		}

		body, err := json.Marshal(msg)
		if err != nil {
			log.Println("Xabar marshalingda xatolik:", err)
			continue
		}

		// routing key = queue bilan bind qilingan key
		err = ch.Publish(
			exchange,
			queueName, // endi routing key = queueName
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Println("Xabar yuborishda xatolik:", err)
			continue
		}

		fmt.Println("Xabar yuborildi:", string(body))
	}
}
