package producer

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	MsgBroker struct {
		channel *amqp.Channel
		logger  *slog.Logger
	}
)

func NewMsgBroker(channel *amqp.Channel, logger *slog.Logger) *MsgBroker {
	return &MsgBroker{
		channel: channel,
		logger:  logger,
	}
}

func (b *MsgBroker) CreateBooking(body []byte) error {
	fmt.Println("------------------------------------------CREATE BOOOK -----------------------------------")
	return b.publishMessage("create_booking", body)
}

func (b *MsgBroker) CancelBooking(body []byte) error {
	fmt.Println("------------------------------------------Cancel BOOOK -----------------------------------")

	return b.publishMessage("booking_cancelled", body)
}

func (b *MsgBroker) Payment(body []byte) error {
	fmt.Println("------------------------------------------PAYMENT ------------------------------")
	return b.publishMessage("payment_processed", body)
}

func (b *MsgBroker) Review(body []byte) error {
	fmt.Println("------------------------------------------REVIEW ------------------------------")
	return b.publishMessage("review_submitted", body)
}

func (b *MsgBroker) CreateNotification(body []byte) error {
	return b.publishMessage("create_notification", body)
}

func (b *MsgBroker) publishMessage(queueName string, body []byte) error {
	err := b.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		b.logger.Error("Failed to publish message", "queue", queueName, "error", err.Error())
		return err
	}

	b.logger.Info("Message published", "queue", queueName)
	return nil
}
