package config

import (
	"fmt"
	"log"
	"os"

	
	"github.com/streadway/amqp"
)


type RabbitMQConnection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}



func NewRabbitMQ() *RabbitMQConnection {
	user := os.Getenv("RABBITMQ_USER")
	
	pass := os.Getenv("RABBITMQ_PASS")

	host := os.Getenv("RABBITMQ_HOST")
	url := fmt.Sprintf("amqp://%s:%s@%s/%%2F", user, pass, host)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf(" RabbitMQ ga ulanib bo‘lmadi: %v", err)
	}
	fmt.Println("RabbitMQ ga ulanish muvaffaqiyatli")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf(" Channel ochib bo‘lmadi: %v", err)
	}
	fmt.Println("Channel ochildi")

	r := &RabbitMQConnection{
		Conn:    conn,
		Channel: ch,
	}


	if err := r.setupQueues(); err != nil {
		log.Fatalf(" Queue setupda xatolik: %v", err)
	}
	fmt.Println(" RabbitMQ setup tugadi")

	return r
}


func (r *RabbitMQConnection) setupQueues() error {
	exchange := "direct_exchange"
	queue := "arena_queue"
	routingKey := "queue_key"
	dlq := "arena_queue_dlq"
	retry := "arena_queue_retry"



	if err := r.Channel.ExchangeDeclare(
		exchange, // nomi
		"direct", // turi
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // noWait
		nil,      // args
	); err != nil {
		return fmt.Errorf("exchange yaratishda xatolik: %w", err)
	}

	// 2️⃣ DLQ (Dead Letter Queue)
	if _, err := r.Channel.QueueDeclare(
		dlq,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("DLQ yaratishda xatolik: %w", err)
	}

	// 3️⃣ Retry Queue
	retryArgs := amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": routingKey,
		"x-message-ttl":             int32(10000), // 10s delay
	}
	if _, err := r.Channel.QueueDeclare(
		retry,
		true,
		false,
		false,
		false,
		retryArgs,
	); err != nil {
		return fmt.Errorf("retry queue yaratishda xatolik: %w", err)
	}


	// 5️⃣ Asosiy Queue (priority va DLQ bilan)
	args := amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": dlq,
		"x-max-priority":            int32(10),
	}
	if _, err := r.Channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		args,
	); err != nil {
		return fmt.Errorf("asosiy queue yaratishda xatolik: %w", err)
	}

	// 6️⃣ Binding (queue ↔ exchange)
	if err := r.Channel.QueueBind(
		queue,      // queue
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	); err != nil {
		return fmt.Errorf("queue bindda xatolik: %w", err)
	}

	return nil
}