package mongodb

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/models"
	"Booking_Service/storage"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ProviderRepo struct {
	collection *mongo.Collection
}

func NewProviderRepo(db *mongo.Database) storage.IProviderStorage {
	return &ProviderRepo{
		collection: db.Collection("providers"),
	}
}

func (r *ProviderRepo) RegisterProvider(ctx context.Context, req *pb.RegisterProviderRequest) (*pb.ProviderResponse, error) {
	provider := &models.Provider{
		UserID:        req.UserId,
		CompanyName:   req.CompanyName,
		Description:   req.Description,
		Services:      req.Services,
		Availability:  req.Availability,
		AverageRating: req.AverageRating,
		Location:      &models.GeoPoint{Latitude: req.Location.Latitude, Longitude: req.Location.Longitude},
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}
	res, err := r.collection.InsertOne(ctx, provider)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(primitive.ObjectID)

	return &pb.ProviderResponse{
		Id:            id.Hex(),
		UserId:        provider.UserID,
		CompanyName:   provider.CompanyName,
		Description:   provider.Description,
		Services:      provider.Services,
		AverageRating: provider.AverageRating,
		Location:      req.Location,
		CreatedAt:     provider.CreatedAt,
		UpdatedAt:     provider.UpdatedAt,
	}, nil
}

func (r *ProviderRepo) GetProvider(ctx context.Context, id *pb.IdRequest) (*pb.ProviderResponse, error) {
	var provider models.Provider1
	ID, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&provider)
	if err != nil {
		return nil, err
	}

	return &pb.ProviderResponse{
		Id:            provider.ID.Hex(),
		UserId:        provider.UserID,
		CompanyName:   provider.CompanyName,
		Description:   provider.Description,
		Services:      provider.Services,
		AverageRating: provider.AverageRating,
		CreatedAt:     provider.CreatedAt,
		UpdatedAt:     provider.UpdatedAt,
	}, nil
}

func (r *ProviderRepo) UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest) (*pb.ProviderResponse, error) {
	ID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"company_name":   req.CompanyName,
			"description":    req.Description,
			"services":       req.Services,
			"availability":   req.Availability,
			"average_rating": req.AverageRating,
			"location":       req.Location,
			"updated_at":     time.Now().Format(time.RFC3339),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": ID}, update)
	if err != nil {
		return nil, err
	}

	return r.GetProvider(ctx, &pb.IdRequest{Id: req.Id})
}

func (r *ProviderRepo) DeleteProvider(ctx context.Context, id *pb.IdRequest) error {
	ID, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": ID})
	return err
}

func (r *ProviderRepo) ListProviders(ctx context.Context, req *pb.ListProvidersRequest) (*pb.ListProvidersResponse, error) {
	var providers []*models.Provider
	opts := options.Find().
		SetSkip(int64((req.Page - 1) * req.Limit)).
		SetLimit(int64(req.Limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &providers); err != nil {
		return nil, err
	}

	var providerResponses []*pb.ProviderResponse
	for _, provider := range providers {
		providerResponses = append(providerResponses, &pb.ProviderResponse{
			Id:            provider.ID.Hex(),
			UserId:        provider.UserID,
			CompanyName:   provider.CompanyName,
			Description:   provider.Description,
			Services:      provider.Services,
			AverageRating: provider.AverageRating,
			Location:      &pb.GeoPoint{Latitude: provider.Location.Latitude, Longitude: provider.Location.Longitude},
			CreatedAt:     provider.CreatedAt,
			UpdatedAt:     provider.UpdatedAt,
		})
	}

	return &pb.ListProvidersResponse{
		Providers: providerResponses,
	}, nil
}
