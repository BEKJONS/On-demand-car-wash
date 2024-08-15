package mongodb

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/storage"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchingRepo struct {
	collection *mongo.Collection
}

func NewSearchingRepo(db *mongo.Database) storage.ISearchStorage {
	return &SearchingRepo{
		collection: db.Collection("providers"), // Default collection to "providers"
	}
}

func (s *SearchingRepo) SearchProviders(ctx context.Context, req *pb.Filter) (*pb.ListProvidersResponses, error) {

	filter := bson.M{}
	if req.Location != "" {
		filter["location"] = req.Location
	}
	if req.Rating > 0 {
		filter["average_rating"] = bson.M{"$gte": req.Rating}
	}
	if req.CompanyName != "" {
		filter["company_name"] = req.CompanyName
	}

	opts := options.Find()
	opts.SetSkip(int64(req.Page * req.Limit))
	opts.SetLimit(int64(req.Limit))

	if req.ByRating {
		opts.SetSort(bson.D{{Key: "average_rating", Value: -1}})
	}

	if req.NumberOfComments {
		opts.SetSort(bson.D{{Key: "number_of_comments", Value: -1}})
	}

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var providers []*pb.Providers
	for cursor.Next(ctx) {
		var provider pb.Providers
		if err := cursor.Decode(&provider); err != nil {
			fmt.Println(err)
			return nil, err
		}
		providers = append(providers, &provider)
	}

	return &pb.ListProvidersResponses{Provideres: providers}, nil
}

func (s *SearchingRepo) SearchServices(ctx context.Context, req *pb.Filter) (*pb.ListServicesResponses, error) {
	collection := s.collection.Database().Collection("services")

	filter := bson.M{}
	if req.Price > 0 {
		filter["price"] = bson.M{"$lte": req.Price}
	}
	if req.Location != "" {
		filter["location"] = req.Location
	}
	if req.Duration > 0 {
		filter["duration"] = bson.M{"$gte": req.Duration}
	}

	opts := options.Find()
	opts.SetSkip(int64(req.Page * req.Limit))
	opts.SetLimit(int64(req.Limit))

	if req.ByPrice {
		opts.SetSort(bson.D{{Key: "price", Value: 1}})
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []*pb.Service
	for cursor.Next(ctx) {
		var service pb.Service
		if err := cursor.Decode(&service); err != nil {
			return nil, err
		}
		services = append(services, &service)
	}

	return &pb.ListServicesResponses{Services: services}, nil
}
