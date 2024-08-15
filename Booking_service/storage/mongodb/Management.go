package mongodb

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/models"
	"Booking_Service/storage"
	Redis "Booking_Service/storage/redis"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ManagementRepo struct {
	collection *mongo.Collection
}

func NewManagementRepo(db *mongo.Database) storage.IManagementStorage {
	return &ManagementRepo{
		collection: db.Collection("services"),
	}
}

func (r *ManagementRepo) CreateService(ctx context.Context, service *pb.CreateServiceRequest) (*pb.ServiceResponse, error) {
	req := &models.Service{
		Name:        service.Name,
		Description: service.Description,
		Price:       service.Price,
		Duration:    service.Duration,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	res, err := r.collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(primitive.ObjectID)

	return &pb.ServiceResponse{
		Id:          id.String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Duration:    req.Duration,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}, err
}

func (r *ManagementRepo) GetServiceByID(ctx context.Context, id *pb.IdRequest) (*pb.ServiceResponse, error) {
	var service pb.ServiceResponse
	ID, err := primitive.ObjectIDFromHex(id.Id)
	err = r.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&service)
	if err != nil {
		return nil, err
	}
	return &service, err
}

func (r *ManagementRepo) UpdateService(ctx context.Context, service *pb.UpdateServiceRequest) error {
	update := bson.M{
		"$set": bson.M{
			"name":        service.Name,
			"description": service.Description,
			"duration":    service.Duration,
			"price":       service.Price,
			"updated_at":  time.Now().Format(time.RFC3339),
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": service.Id},
		update,
	)
	return err
}

func (r *ManagementRepo) DeleteService(ctx context.Context, id *pb.IdRequest) error {
	ID, err := primitive.ObjectIDFromHex(id.Id)
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": ID})
	return err
}

func (r *ManagementRepo) ListServices(ctx context.Context, req *pb.ListServicesRequest) (*pb.ListServicesResponse, error) {
	var services []*pb.ServiceResponse
	opts := options.Find().
		SetSkip(int64((req.Page - 1) * req.Limit)).
		SetLimit(int64(req.Limit))
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &services); err != nil {
		return nil, err
	}
	return &pb.ListServicesResponse{Services: services}, nil
}

func (r *ManagementRepo) PopularServices(ctx context.Context, req *pb.Void) (*pb.ListServicesResponse, error) {
	conn := Redis.ConnectRedis()

	popularServices, err := Redis.GetPopularServices(conn, 10)
	if err != nil {
		return nil, err
	}

	response := &pb.ListServicesResponse{}

	for _, service := range popularServices {
		serviceID := service.Member.(string)
		var service pb.ServiceResponse
		ID, err := primitive.ObjectIDFromHex(serviceID)
		if err != nil {
			return nil, err
		}
		err = r.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&service)
		if err != nil {
			return nil, err
		}

		err = r.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&service)
		mockService := &service

		response.Services = append(response.Services, mockService)
	}
	return response, nil
}
