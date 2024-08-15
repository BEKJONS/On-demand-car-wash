package pkg

import (
	"Api_Gateway/config"
	pbb "Api_Gateway/genproto/booking"
	pbu "Api_Gateway/genproto/users"
	"log"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(cfg *config.Config) pbu.AuthServiceClient {
	conn, err := grpc.NewClient(cfg.AUTH_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbu.NewAuthServiceClient(conn)
}

func NewAdminClient(cfg *config.Config) pbu.AdminClient {
	conn, err := grpc.NewClient(cfg.AUTH_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbu.NewAdminClient(conn)
}

func NewBookingClient(cfg *config.Config) pbb.BookingServiceClient {
	conn, err := grpc.NewClient(cfg.BOOKING_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbb.NewBookingServiceClient(conn)
}
func NewServiceClient(cfg *config.Config) pbb.ServiceManagementServiceClient {
	conn, err := grpc.NewClient(cfg.BOOKING_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}
	return pbb.NewServiceManagementServiceClient(conn)
}
func NewPaymentClient(cfg *config.Config) pbb.PaymentServiceClient {
	conn, err := grpc.NewClient(cfg.BOOKING_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}
	return pbb.NewPaymentServiceClient(conn)
}
func NewReviewClient(cfg *config.Config) pbb.ReviewServiceClient {
	conn, err := grpc.NewClient(cfg.BOOKING_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}
	return pbb.NewReviewServiceClient(conn)
}
func NewProviderClient(cfg *config.Config) pbb.ProviderServiceClient {
	conn, err := grpc.NewClient(cfg.BOOKING_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}
	return pbb.NewProviderServiceClient(conn)
}

func NewSearchingClient(cfg *config.Config) pbb.SearchingServiceClient {
	conn, err := grpc.NewClient(cfg.BOOKING_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}
	return pbb.NewSearchingServiceClient(conn)
}
func NewNotificationClient(cfg *config.Config) pbb.NotificationsClient {
	conn, err := grpc.NewClient(cfg.BOOKING_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}
	return pbb.NewNotificationsClient(conn)
}
