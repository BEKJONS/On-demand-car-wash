package storage

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/models"
)

import (
	"context"
)

type IStorage interface {
	Booking() IBookingStorage
	Management() IManagementStorage
	Review() IReviewStorage
	Search() ISearchStorage
	Payment() IPaymentStorage
	Provider() IProviderStorage
	Notification() INotificationStorage
	Close()
}

type IBookingStorage interface {
	Add(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error)
	Read(ctx context.Context, req *pb.IdRequest) (*pb.BookingResponse, error)
	Update(ctx context.Context, req *pb.UpdateBookingRequest) (*pb.BookingResponse, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	Cancel(ctx context.Context, id *pb.IdRequest) (*pb.BookingResponse, error)
	List(ctx context.Context, req *pb.ListBookingsRequest) (*pb.ListBookingsResponse, error)
}

type IManagementStorage interface {
	CreateService(ctx context.Context, service *pb.CreateServiceRequest) (*pb.ServiceResponse, error)
	GetServiceByID(ctx context.Context, id *pb.IdRequest) (*pb.ServiceResponse, error)
	UpdateService(ctx context.Context, service *pb.UpdateServiceRequest) error
	DeleteService(ctx context.Context, id *pb.IdRequest) error
	ListServices(ctx context.Context, req *pb.ListServicesRequest) (*pb.ListServicesResponse, error)
	PopularServices(ctx context.Context, req *pb.Void) (*pb.ListServicesResponse, error)
}

type IPaymentStorage interface {
	CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.PaymentResponse, error)
	GetPayment(ctx context.Context, id *pb.IdRequest) (*pb.PaymentResponse, error)
	ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error)
}

type IReviewStorage interface {
	Create(ctx context.Context, req *pb.CreateReviewRequest) (*pb.ReviewResponse, error)
	GetByID(id *pb.IdRequest) (*pb.ReviewResponse, error)
	GetAll(ctx context.Context, req *pb.ListReviewsRequest) (*pb.ListReviewsResponse, error)
	Update(review *pb.UpdateReviewRequest) error
	Delete(id *pb.IdRequest) error
}

type IProviderStorage interface {
	RegisterProvider(ctx context.Context, req *pb.RegisterProviderRequest) (*pb.ProviderResponse, error)
	GetProvider(ctx context.Context, id *pb.IdRequest) (*pb.ProviderResponse, error)
	ListProviders(ctx context.Context, req *pb.ListProvidersRequest) (*pb.ListProvidersResponse, error)
	UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest) (*pb.ProviderResponse, error)
	DeleteProvider(ctx context.Context, req *pb.IdRequest) error
}

type ISearchStorage interface {
	SearchProviders(ctx context.Context, req *pb.Filter) (*pb.ListProvidersResponses, error)
	SearchServices(ctx context.Context, req *pb.Filter) (*pb.ListServicesResponses, error)
}
type INotificationStorage interface {
	Create(ctx context.Context, req *models.NewNotification) (string, error)
	Get(ctx context.Context, id string) (*models.Notification, error)
}
