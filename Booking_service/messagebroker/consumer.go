package msgbroker

import (
	genprotos "Booking_Service/genproto/booking"
	"Booking_Service/service"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type (
	MsgBroker struct {
		BService           *service.BookingService
		PService           *service.PaymentService
		RService           *service.ReviewService
		NService           *service.NotificationService
		channel            *amqp.Channel
		createBooking      <-chan amqp.Delivery
		deleteBooking      <-chan amqp.Delivery
		createPayment      <-chan amqp.Delivery
		createReview       <-chan amqp.Delivery
		createNotification <-chan amqp.Delivery
		logger             *slog.Logger
		wg                 *sync.WaitGroup
		numberOfServices   int
	}
)

func New(
	Bservice *service.BookingService,
	Pservice *service.PaymentService,
	Rservice *service.ReviewService,
	Nservice *service.NotificationService,
	channel *amqp.Channel,
	logger *slog.Logger,
	createBooking <-chan amqp.Delivery,
	deleteBooking <-chan amqp.Delivery,
	createPayment <-chan amqp.Delivery,
	createReview <-chan amqp.Delivery,
	createNotification <-chan amqp.Delivery,
	wg *sync.WaitGroup,
	numberOfServices int) *MsgBroker {
	return &MsgBroker{
		BService:           Bservice,
		PService:           Pservice,
		RService:           Rservice,
		NService:           Nservice,
		channel:            channel,
		createBooking:      createBooking,
		deleteBooking:      deleteBooking,
		createPayment:      createPayment,
		createReview:       createReview,
		createNotification: createNotification,
		logger:             logger,
		wg:                 wg,
		numberOfServices:   numberOfServices,
	}
}

func (m *MsgBroker) StartToConsume(ctx context.Context) {
	m.wg.Add(m.numberOfServices)
	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go m.consumeMessages(consumerCtx, m.createBooking, "create_booking")
	go m.consumeMessages(consumerCtx, m.createPayment, "payment_processed")
	go m.consumeMessages(consumerCtx, m.deleteBooking, "booking_cancelled")
	go m.consumeMessages(consumerCtx, m.createReview, "review_submitted")
	go m.consumeMessages(consumerCtx, m.createNotification, "create_notification")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	m.logger.Info("Shutting down, waiting for consumers to finish")
	cancel()
	m.wg.Wait()
	m.logger.Info("All consumers have stopped")
}

func (m *MsgBroker) consumeMessages(ctx context.Context, messages <-chan amqp.Delivery, logPrefix string) {
	defer m.wg.Done()
	for {
		select {
		case val := <-messages:
			var err error

			switch logPrefix {
			case "create_booking":
				var req genprotos.CreateBookingRequest

				if err := json.Unmarshal(val.Body, &req); err != nil {
					val.Nack(false, false)
					continue
				}
				_, err = m.BService.CreateBooking(ctx, &req)
				if err != nil {
					val.Nack(false, false)
					continue
				}
				val.Ack(false)

			case "booking_updated":
				var req genprotos.UpdateBookingRequest
				if err := json.Unmarshal(val.Body, &req); err != nil {
					val.Nack(false, false)
					continue
				}
				_, err = m.BService.UpdateBooking(ctx, &req)

			case "booking_cancelled":
				var req genprotos.IdRequest
				if err := json.Unmarshal(val.Body, &req); err != nil {
					val.Nack(false, false)
					continue
				}
				_, err = m.BService.CancelBooking(ctx, &req)

			case "payment_processed":
				var req genprotos.CreatePaymentRequest
				if err := json.Unmarshal(val.Body, &req); err != nil {
					val.Nack(false, false)
					continue
				}
				_, err = m.PService.CreatePayment(ctx, &req)

			case "review_submitted":
				var req genprotos.CreateReviewRequest
				if err := json.Unmarshal(val.Body, &req); err != nil {
					val.Nack(false, false)
					continue
				}
				_, err = m.RService.CreateReview(ctx, &req)

			case "create_notification":
				var req genprotos.NewNotification
				if err := json.Unmarshal(val.Body, &req); err != nil {
					val.Nack(false, false)
					continue
				}
				_, err = m.NService.CreateNotification(ctx, &req)
			}

			if err != nil {
				val.Nack(false, false)
				continue
			}

			val.Ack(false)

		case <-ctx.Done():
			m.logger.Info("Context done, stopping consumer", "consumer", logPrefix)
			return
		}
	}
}
