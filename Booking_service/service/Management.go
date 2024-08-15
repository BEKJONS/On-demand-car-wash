package service

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/storage"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

type ManagementService struct {
	pb.UnimplementedServiceManagementServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewManagementService(s storage.IStorage, logger *slog.Logger) *ManagementService {
	return &ManagementService{
		storage: s,
		logger:  logger,
	}
}

func (s *ManagementService) CreateService(ctx context.Context, req *pb.CreateServiceRequest) (*pb.ServiceResponse, error) {
	s.logger.Info("Creating service", "name", req.Name)
	res, err := s.storage.Management().CreateService(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create service", "error", err)
		return nil, err
	}
	s.logger.Info("Service created successfully", "name", req.Name)
	return res, nil
}

func (s *ManagementService) GetServiceByID(ctx context.Context, req *pb.IdRequest) (*pb.ServiceResponse, error) {
	s.logger.Info("Fetching service by ID", "id", req.Id)
	service, err := s.storage.Management().GetServiceByID(ctx, req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Warn("Service not found", "id", req.Id)
			return nil, errors.New("service not found")
		}
		s.logger.Error("Failed to fetch service", "error", err)
		return nil, err
	}
	s.logger.Info("Service fetched successfully", "id", req.Id)
	return service, nil
}

func (s *ManagementService) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest) (*pb.ServiceResponse, error) {
	s.logger.Info("Updating service", "id", req.Id)
	err := s.storage.Management().UpdateService(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update service", "error", err)
		return nil, err
	}
	s.logger.Info("Service updated successfully", "id", req.Id)
	return &pb.ServiceResponse{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Duration:    req.Duration,
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

func (s *ManagementService) DeleteService(ctx context.Context, req *pb.IdRequest) (*pb.Void, error) {
	s.logger.Info("Deleting service", "id", req.Id)
	err := s.storage.Management().DeleteService(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete service", "error", err)
		return nil, err
	}
	s.logger.Info("Service deleted successfully", "id", req.Id)
	return nil, nil
}

func (s *ManagementService) ListServices(ctx context.Context, req *pb.ListServicesRequest) (*pb.ListServicesResponse, error) {
	s.logger.Info("Listing services", "page", req.Page, "limit", req.Limit)
	services, err := s.storage.Management().ListServices(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list services", "error", err)
		return nil, err
	}
	s.logger.Info("Services listed successfully", "count", len(services.Services))
	return services, nil
}

func (s *ManagementService) PopularServices(ctx context.Context, req *pb.Void) (*pb.ListServicesResponse, error) {
	s.logger.Info("Popular services")
	services, err := s.storage.Management().PopularServices(ctx, req)
	if err != nil {
		s.logger.Error("Failed to popular services", "error", err)
		return nil, err
	}
	s.logger.Info("Services populated successfully", "count", len(services.Services))
	return services, nil
}
