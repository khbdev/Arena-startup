package rabbitmq

import (
	"ai-service/internal/service"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

func StartConsumer(ch *amqp.Channel) {
	queueName := "arena_queue"

	msgs, err := ch.Consume(
		queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Queue’dan xabar ololmadi: %v", err)
	}

	fmt.Println("Consumer ishga tushdi, xabarlar kutilyapti...")

	var mu sync.Mutex
	var wg sync.WaitGroup

	for d := range msgs {
		wg.Add(1)

		// Ketma-ket ishlashi uchun lock
		mu.Lock()
		service.ProcessMessage(d.Body)
		mu.Unlock()

		wg.Done()
	}

	wg.Wait() // barcha xabarlar tugashini kutadi
}
