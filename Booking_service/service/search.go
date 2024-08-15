package service

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/storage"
	"context"
	"github.com/pkg/errors"
	"log/slog"
)

type SearchService struct {
	pb.UnimplementedSearchingServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewSearchService(storage storage.IStorage, logger *slog.Logger) *SearchService {
	return &SearchService{
		storage: storage,
		logger:  logger,
	}
}

func (s *SearchService) SearchProviders(ctx context.Context, req *pb.Filter) (*pb.ListProvidersResponses, error) {
	s.logger.Info("SearchProviders", "company_name", req.CompanyName)

	provider, err := s.storage.Search().SearchProviders(ctx, req)
	if err != nil {
		s.logger.Error("Failed to search provider", "error", err)
		return nil, errors.Wrap(err, "could not search provider")
	}

	s.logger.Info("Review created successfully")
	return provider, nil
}

func (s *SearchService) SearchServices(ctx context.Context, req *pb.Filter) (*pb.ListServicesResponses, error) {
	s.logger.Info("Fetching review")

	review, err := s.storage.Search().SearchServices(ctx, req)
	if err != nil {
		s.logger.Error("Failed to fetch review", "error", err)
		return nil, errors.Wrap(err, "could not fetch review")
	}

	s.logger.Info("Review fetched successfully")
	return review, nil
}
