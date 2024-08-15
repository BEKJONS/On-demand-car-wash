package handler

import (
	"Api_Gateway/config"
	pbb "Api_Gateway/genproto/booking"
	pb "Api_Gateway/genproto/users"
	"Api_Gateway/pkg"
	"Api_Gateway/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"

	broker "Api_Gateway/producer"
)

type str string

type Handler struct {
	User           pb.AuthServiceClient
	Admin          pb.AdminClient
	Booking        pbb.BookingServiceClient
	Service        pbb.ServiceManagementServiceClient
	Payment        pbb.PaymentServiceClient
	Review         pbb.ReviewServiceClient
	Provider       pbb.ProviderServiceClient
	Search         pbb.SearchingServiceClient
	Notification   pbb.NotificationsClient
	Log            *slog.Logger
	ContextTimeout time.Duration
	UserIDKey      str
	Broker         *broker.MsgBroker
}

func NewHandler(cfg *config.Config, conn *amqp.Channel) *Handler {
	return &Handler{
		User:           pkg.NewUserClient(cfg),
		Admin:          pkg.NewAdminClient(cfg),
		Booking:        pkg.NewBookingClient(cfg),
		Service:        pkg.NewServiceClient(cfg),
		Payment:        pkg.NewPaymentClient(cfg),
		Review:         pkg.NewReviewClient(cfg),
		Provider:       pkg.NewProviderClient(cfg),
		Notification:   pkg.NewNotificationClient(cfg),
		Search:         pkg.NewSearchingClient(cfg),
		Log:            logger.NewLogger(),
		ContextTimeout: 5 * time.Second,
		UserIDKey:      str("user_id"),
		Broker:         broker.NewMsgBroker(conn, logger.NewLogger()),
	}
}
