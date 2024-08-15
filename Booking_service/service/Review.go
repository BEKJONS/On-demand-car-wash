package service

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/storage"
	"context"
	"github.com/pkg/errors"
	"log/slog"
)

type ReviewService struct {
	pb.UnimplementedReviewServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewReviewService(storage storage.IStorage, logger *slog.Logger) *ReviewService {
	return &ReviewService{
		storage: storage,
		logger:  logger,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.ReviewResponse, error) {
	s.logger.Info("Creating review", "booking_id", req.BookingId, "user_id", req.UserId, "provider_id", req.ProviderId)

	review, err := s.storage.Review().Create(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create review", "error", err)
		return nil, errors.Wrap(err, "could not create review")
	}

	s.logger.Info("Review created successfully", "review_id", review.Id)
	return review, nil
}

func (s *ReviewService) GetReviewById(ctx context.Context, req *pb.IdRequest) (*pb.ReviewResponse, error) {
	s.logger.Info("Fetching review", "review_id", req.Id)

	review, err := s.storage.Review().GetByID(req)
	if err != nil {
		s.logger.Error("Failed to fetch review", "error", err)
		return nil, errors.Wrap(err, "could not fetch review")
	}

	s.logger.Info("Review fetched successfully", "review_id", review.Id)
	return review, nil
}

func (s *ReviewService) ListReviews(ctx context.Context, req *pb.ListReviewsRequest) (*pb.ListReviewsResponse, error) {
	s.logger.Info("Listing reviews", "provider_id", req.ProviderId, "page", req.Page, "limit", req.Limit)

	reviews, err := s.storage.Review().GetAll(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list reviews", "error", err)
		return nil, errors.Wrap(err, "could not list reviews")
	}

	s.logger.Info("Reviews listed successfully", "count", len(reviews.Reviews))
	return reviews, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.ReviewResponse, error) {
	s.logger.Info("Updating review", "review_id", req.Id)

	err := s.storage.Review().Update(req)
	if err != nil {
		s.logger.Error("Failed to update review", "error", err)
		return nil, errors.Wrap(err, "could not update review")
	}

	s.logger.Info("Review updated successfully", "review_id", req.Id)

	// Retrieve the updated review to return
	updatedReview, err := s.storage.Review().GetByID(&pb.IdRequest{Id: req.Id})
	if err != nil {
		s.logger.Error("Failed to fetch updated review", "error", err)
		return nil, errors.Wrap(err, "could not fetch updated review")
	}

	return updatedReview, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.IdRequest) (*pb.Void, error) {
	s.logger.Info("Deleting review", "review_id", req.Id)

	err := s.storage.Review().Delete(req)
	if err != nil {
		s.logger.Error("Failed to delete review", "error", err)
		return nil, errors.Wrap(err, "could not delete review")
	}

	s.logger.Info("Review deleted successfully", "review_id", req.Id)
	return nil, nil
}
