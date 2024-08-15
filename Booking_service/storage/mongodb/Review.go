package mongodb

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/models"
	"Booking_Service/storage"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepo struct {
	collection *mongo.Collection
}

func NewReviewRepo(db *mongo.Database) storage.IReviewStorage {
	return &ReviewRepo{
		collection: db.Collection("Reviews"),
	}
}

// Create inserts a new review into the database
func (r *ReviewRepo) Create(ctx context.Context, req *pb.CreateReviewRequest) (*pb.ReviewResponse, error) {
	review := &models.Review{
		BookingID:  req.BookingId,
		UserID:     req.UserId,
		ProviderID: req.ProviderId,
		Rating:     req.Rating,
		Comment:    req.Comment,
		CreatedAt:  time.Now().Format(time.RFC3339),
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}

	// Insert the review into MongoDB
	res, err := r.collection.InsertOne(ctx, review)
	if err != nil {
		return nil, err
	}

	// Get the inserted ID and convert it to a string
	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert inserted ID to ObjectID")
	}

	// Create the response object
	return &pb.ReviewResponse{
		Id:         objID.Hex(),
		BookingId:  review.BookingID,
		UserId:     review.UserID,
		ProviderId: review.ProviderID,
		Rating:     review.Rating,
		Comment:    review.Comment,
		CreatedAt:  review.CreatedAt,
		UpdatedAt:  review.UpdatedAt,
	}, nil
}

// GetByID retrieves a review by its ID
func (r *ReviewRepo) GetByID(id *pb.IdRequest) (*pb.ReviewResponse, error) {
	var review *pb.ReviewResponse
	Id, err := primitive.ObjectIDFromHex(id.Id)
	err = r.collection.FindOne(context.Background(), bson.M{"_id": Id}).Decode(&review)
	if err != nil {
		return nil, err
	}
	return review, nil
}

// GetAll retrieves all reviews from the database
func (r *ReviewRepo) GetAll(ctx context.Context, req *pb.ListReviewsRequest) (*pb.ListReviewsResponse, error) {
	filter := bson.M{}
	if req.ProviderId != "" {
		filter["provider_id"] = req.ProviderId
	}
	opts := options.Find().
		SetSkip(int64((req.Page - 1) * req.Limit)).
		SetLimit(int64(req.Limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list reviews")
	}
	defer cursor.Close(ctx)

	var reviews []*pb.ReviewResponse
	for cursor.Next(ctx) {
		var review models.Review
		if err := cursor.Decode(&review); err != nil {
			return nil, errors.Wrap(err, "failed to decode review")
		}

		reviews = append(reviews, &pb.ReviewResponse{
			Id:         review.ID.Hex(),
			BookingId:  review.BookingID,
			UserId:     review.UserID,
			ProviderId: review.ProviderID,
			Rating:     review.Rating,
			Comment:    review.Comment,
			CreatedAt:  review.CreatedAt,
			UpdatedAt:  review.UpdatedAt,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor error")
	}

	return &pb.ListReviewsResponse{
		Reviews: reviews,
	}, nil
}

// Update modifies an existing review
func (r *ReviewRepo) Update(review *pb.UpdateReviewRequest) error {
	ID, err := primitive.ObjectIDFromHex(review.Id)
	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": ID},
		bson.M{"$set": bson.M{
			"rating":     review.Rating,
			"comment":    review.Comment,
			"updated_at": time.Now().Format(time.RFC3339),
		}},
	)
	return err
}

// Delete removes a review from the database by its ID
func (r *ReviewRepo) Delete(id *pb.IdRequest) error {
	Id, err := primitive.ObjectIDFromHex(id.Id)
	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": Id})
	return err
}
