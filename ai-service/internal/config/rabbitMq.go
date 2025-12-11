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
		log.Fatalf("RabbitMQ ga ulanib bo‘lmadi: %v", err)
	}
	fmt.Println("RabbitMQ ga ulanish muvaffaqiyatli")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Channel ochib bo‘lmadi: %v", err)
	}
	fmt.Println("Channel ochildi")

	r := &RabbitMQConnection{
		Conn:    conn,
		Channel: ch,
	}

	if err := r.setupQueues(); err != nil {
		log.Fatalf("Queue setupda xatolik: %v", err)
	}
	fmt.Println("RabbitMQ setup tugadi")

	return r
}

func (r *RabbitMQConnection) setupQueues() error {
	exchange := "direct_exchange"

	// Arena Queue params
	arenaQueue := "arena_queue"
	arenaRoutingKey := "queue_key"
	arenaDLQ := "arena_queue_dlq"
	arenaRetry := "arena_queue_retry"

	// Notifications Queue params
	notificationsQueue := "notifications_queue"
	notificationsRetry := "notifications_queue_retry"
	notificationsRoutingKey := "notif_key"

	// --- Exchange ---
	if err := r.Channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange yaratishda xatolik: %w", err)
	}

	// --- Arena DLQ ---
	if _, err := r.Channel.QueueDeclare(
		arenaDLQ,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("arena DLQ yaratishda xatolik: %w", err)
	}

	// --- Arena Retry Queue (10s) ---
	if _, err := r.Channel.QueueDeclare(
		arenaRetry,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    exchange,
			"x-dead-letter-routing-key": arenaRoutingKey,
			"x-message-ttl":             int32(10000),
		},
	); err != nil {
		return fmt.Errorf("arena retry queue yaratishda xatolik: %w", err)
	}

	// --- Arena Main Queue (priority + DLQ) ---
	if _, err := r.Channel.QueueDeclare(
		arenaQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    exchange,
			"x-dead-letter-routing-key": arenaDLQ,
			"x-max-priority":            int32(10),
		},
	); err != nil {
		return fmt.Errorf("arena queue yaratishda xatolik: %w", err)
	}

	if err := r.Channel.QueueBind(arenaQueue, arenaRoutingKey, exchange, false, nil); err != nil {
		return fmt.Errorf("arena queue bindda xatolik: %w", err)
	}

	// --- Notifications Retry Queue (1 soat TTL) ---
	if _, err := r.Channel.QueueDeclare(
		notificationsRetry,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    exchange,
			"x-dead-letter-routing-key": notificationsRoutingKey,
			"x-message-ttl":             int32(3600 * 1000),
		},
	); err != nil {
		return fmt.Errorf("notifications retry queue yaratishda xatolik: %w", err)
	}

	// --- Notifications Main Queue ---
	if _, err := r.Channel.QueueDeclare(
		notificationsQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    exchange,
			"x-dead-letter-routing-key": notificationsRetry,
		},
	); err != nil {
		return fmt.Errorf("notifications queue yaratishda xatolik: %w", err)
	}

	if err := r.Channel.QueueBind(notificationsQueue, notificationsRoutingKey, exchange, false, nil); err != nil {
		return fmt.Errorf("notifications queue bindda xatolik: %w", err)
	}

	return nil
}
