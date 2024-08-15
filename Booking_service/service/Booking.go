package service

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/storage"
	"context"
	"log/slog"
)

type BookingService struct {
	pb.UnimplementedBookingServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewBookingService(s storage.IStorage, logger *slog.Logger) *BookingService {
	return &BookingService{
		storage: s,
		logger:  logger,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error) {
	s.logger.Info("CreateBooking called", "request", req)

	resp, err := s.storage.Booking().Add(ctx, req)

	if err != nil {
		s.logger.Error("failed to create booking", "error", err)
		return nil, err
	}

	s.logger.Info("Booking created successfully", "response", resp)
	return resp, nil
}

func (s *BookingService) GetBooking(ctx context.Context, req *pb.IdRequest) (*pb.BookingResponse, error) {
	s.logger.Info("GetBooking called", "request", req)

	resp, err := s.storage.Booking().Read(ctx, req)
	if err != nil {
		s.logger.Error("failed to get booking", "error", err)
		return nil, err
	}

	s.logger.Info("Booking retrieved successfully", "response", resp)
	return resp, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, req *pb.UpdateBookingRequest) (*pb.BookingResponse, error) {
	s.logger.Info("UpdateBooking called", "request", req)

	resp, err := s.storage.Booking().Update(ctx, req)
	if err != nil {
		s.logger.Error("failed to update booking", "error", err)
		return nil, err
	}

	s.logger.Info("Booking updated successfully", "response", resp)
	return resp, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, req *pb.IdRequest) (*pb.BookingResponse, error) {
	s.logger.Info("CancelBooking called", "request", req)

	resp, err := s.storage.Booking().Cancel(ctx, req)
	if err != nil {
		s.logger.Error("failed to cancel booking", "error", err)
		return nil, err
	}

	s.logger.Info("Booking canceled successfully", "response", resp)
	return resp, nil
}

func (s *BookingService) ListBookings(ctx context.Context, req *pb.ListBookingsRequest) (*pb.ListBookingsResponse, error) {
	s.logger.Info("ListBookings called", "request", req)

	resp, err := s.storage.Booking().List(ctx, req)
	if err != nil {
		s.logger.Error("failed to list bookings", "error", err)
		return nil, err
	}

	s.logger.Info("Bookings listed successfully", "response", resp)
	return resp, nil
}
