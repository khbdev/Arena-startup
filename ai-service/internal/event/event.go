package event

import (
	"ai-service/internal/model"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)


func PublishNotification(ch *amqp.Channel, telegramID int64, testID string) error {
	
	event := model.NotificationEvent{
		TelegramID: telegramID,
		TestID:     testID,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("JSON marshal xatolik: %v", err)
	}

	queueName := "notifications_queue"

	
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key = queue nomi
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("RabbitMQ publish xatolik: %v", err)
	}

	log.Println(" Event notifications_queue ga yuborildi:", string(body))
	return nil
}