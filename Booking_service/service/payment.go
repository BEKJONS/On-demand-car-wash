package service

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/storage"
	"context"
	"github.com/pkg/errors"
	"log/slog"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewPaymentService(storage storage.IStorage, logger *slog.Logger) *PaymentService {
	return &PaymentService{
		storage: storage,
		logger:  logger,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.PaymentResponse, error) {
	s.logger.Info("Creating payment", "booking_id", req.BookingId, "amount", req.Amount)

	payment, err := s.storage.Payment().CreatePayment(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create payment", "error", err)
		return nil, errors.Wrap(err, "could not create payment")
	}

	s.logger.Info("Payment created successfully", "payment_id", payment.Id)
	return payment, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, req *pb.IdRequest) (*pb.PaymentResponse, error) {
	s.logger.Info("Fetching payment", "payment_id", req.Id)

	payment, err := s.storage.Payment().GetPayment(ctx, req)
	if err != nil {
		s.logger.Error("Failed to fetch payment", "error", err)
		return nil, errors.Wrap(err, "could not fetch payment")
	}

	s.logger.Info("Payment fetched successfully", "payment_id", payment.Id)
	return payment, nil
}

func (s *PaymentService) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	s.logger.Info("Listing payments", "user_id", req.UserId, "page", req.Page, "limit", req.Limit)

	payments, err := s.storage.Payment().ListPayments(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list payments", "error", err)
		return nil, errors.Wrap(err, "could not list payments")
	}

	s.logger.Info("Payments listed successfully", "count", len(payments.Payments))
	return payments, nil
}
