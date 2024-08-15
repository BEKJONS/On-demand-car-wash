package service

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/storage"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

type ProviderService struct {
	pb.UnimplementedProviderServiceServer
	storage storage.IProviderStorage
	logger  *slog.Logger
}

// NewProviderService creates a new ProviderService.
func NewProviderService(s storage.IStorage, logger *slog.Logger) *ProviderService {
	return &ProviderService{
		storage: s.Provider(),
		logger:  logger,
	}
}

// RegisterProvider handles the registration of a new provider.
func (s *ProviderService) RegisterProvider(ctx context.Context, req *pb.RegisterProviderRequest) (*pb.ProviderResponse, error) {
	s.logger.Info("Registering provider", "user_id", req.UserId)
	res, err := s.storage.RegisterProvider(ctx, req)
	if err != nil {
		s.logger.Error("Failed to register provider", "error", err)
		return nil, err
	}
	s.logger.Info("Provider registered successfully", "company_name", req.CompanyName)
	return res, nil
}

// GetProvider retrieves a provider by its ID.
func (s *ProviderService) GetProvider(ctx context.Context, req *pb.IdRequest) (*pb.ProviderResponse, error) {
	s.logger.Info("Fetching provider by ID", "id", req.Id)
	provider, err := s.storage.GetProvider(ctx, req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Warn("Provider not found", "id", req.Id)
			return nil, errors.New("provider not found")
		}
		s.logger.Error("Failed to fetch provider", "error", err)
		return nil, err
	}
	s.logger.Info("Provider fetched successfully", "id", req.Id)
	return provider, nil
}

// UpdateProvider updates the details of an existing provider.
func (s *ProviderService) UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest) (*pb.ProviderResponse, error) {
	s.logger.Info("Updating provider", "id", req.Id)
	res, err := s.storage.UpdateProvider(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update provider", "error", err)
		return nil, err
	}
	s.logger.Info("Provider updated successfully", "id", req.Id)
	return res, nil
}

// DeleteProvider deletes a provider by its ID.
func (s *ProviderService) DeleteProvider(ctx context.Context, req *pb.IdRequest) (*pb.Void, error) {
	s.logger.Info("Deleting provider", "id", req.Id)
	err := s.storage.DeleteProvider(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete provider", "error", err)
		return nil, err
	}
	s.logger.Info("Provider deleted successfully", "id", req.Id)
	return &pb.Void{}, nil
}

// ListProviders lists all providers with pagination.
func (s *ProviderService) ListProviders(ctx context.Context, req *pb.ListProvidersRequest) (*pb.ListProvidersResponse, error) {
	s.logger.Info("Listing providers", "page", req.Page, "limit", req.Limit)
	providers, err := s.storage.ListProviders(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list providers", "error", err)
		return nil, err
	}
	s.logger.Info("Providers listed successfully", "count", len(providers.Providers))
	return providers, nil
}
