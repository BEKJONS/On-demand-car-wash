package mongodb

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/models"
	"Booking_Service/storage"
	redisDB "Booking_Service/storage/redis"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type BookingRepo struct {
	col  *mongo.Collection
	col1 *mongo.Collection
	Rdb  *redis.Client
}

func NewBookingRepo(db *mongo.Database, rd *redis.Client) storage.IBookingStorage {
	return &BookingRepo{
		col:  db.Collection("bookings"),
		col1: db.Collection("services"),
		Rdb:  rd,
	}
}

func (b *BookingRepo) Add(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error) {
	var service models.Service
	m := &ManagementRepo{
		collection: b.col1,
	}
	// Convert ServiceId to ObjectID
	serviceID, err := primitive.ObjectIDFromHex(req.ServiceId)
	if err != nil {
		return nil, errors.Wrap(err, "invalid service ID format")
	}

	// Find the service by ID
	err = m.collection.FindOne(ctx, bson.M{"_id": serviceID}).Decode(&service)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("service not found")
		}
		return nil, errors.Wrap(err, "failed to find service")
	}
	fmt.Println(service.Price, serviceID)
	a := models.GeoPoint{req.Location.Latitude, req.Location.Longitude}
	// Create a new booking instance
	booking := &models.Booking{
		UserID:        req.UserId,
		ProviderID:    req.ProviderId,
		ServiceID:     req.ServiceId,
		Status:        "pending",
		ScheduledTime: req.ScheduledTime,
		Location:      a,
		TotalPrice:    service.Price + 5,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}
	// Insert the booking into MongoDB
	res, err := b.col.InsertOne(ctx, booking)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "failed to insert booking")
	}

	// Get the inserted ID and convert it to a hex string
	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert inserted ID to ObjectID")
	}

	redisDB.IncrementServiceOrderCount(b.Rdb, req.ServiceId)

	// Return the response
	return &pb.BookingResponse{
		Id:            objID.Hex(),
		UserId:        req.UserId,
		ProviderId:    req.ProviderId,
		ServiceId:     req.ServiceId,
		Status:        booking.Status,
		ScheduledTime: booking.ScheduledTime,
		Location:      req.Location,
		TotalPrice:    booking.TotalPrice,
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
	}, nil
}

func (b *BookingRepo) Read(ctx context.Context, req *pb.IdRequest) (*pb.BookingResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errors.Wrap(err, "invalid booking ID")
	}

	var booking pb.BookingResponse
	err = b.col.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("booking not found")
		}
		return nil, errors.Wrap(err, "failed to decode booking")
	}

	return &booking, nil
}

func (b *BookingRepo) Update(ctx context.Context, req *pb.UpdateBookingRequest) (*pb.BookingResponse, error) {
	m := &ManagementRepo{
		collection: b.col1,
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errors.Wrap(err, "invalid booking ID")
	}

	service := &models.Service{}
	serviceID, err := primitive.ObjectIDFromHex(req.ServiceId)
	if err != nil {
		return nil, errors.Wrap(err, "invalid service ID format")
	}

	// Find the service by ID
	err = m.collection.FindOne(ctx, bson.M{"_id": serviceID}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("service not found")
		}
		return nil, errors.Wrap(err, "failed to find service")
	}
	updateData := bson.M{
		"$set": bson.M{
			"service_id":  req.ServiceId,
			"total_price": service.Price + 5,
			"updated_at":  time.Now().Format(time.RFC3339),
		},
	}

	_, err = b.col.UpdateOne(ctx, bson.M{"_id": id}, updateData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update booking")
	}

	var updatedBooking *pb.BookingResponse
	err = b.col.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedBooking)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve updated booking")
	}

	return updatedBooking, nil
}
func (b *BookingRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "invalid booking ID")
	}

	_, err = b.col.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return errors.Wrap(err, "failed to update status")
	}
	return nil
}

func (b *BookingRepo) Cancel(ctx context.Context, id *pb.IdRequest) (*pb.BookingResponse, error) {
	ID, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return nil, errors.Wrap(err, "invalid booking ID")
	}
	_, err = b.col.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": bson.M{"status": "canceled"}})
	var updatedBooking *pb.BookingResponse
	err = b.col.FindOne(ctx, bson.M{"_id": ID}).Decode(&updatedBooking)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve updated booking")
	}

	return updatedBooking, nil
}

func (b *BookingRepo) List(ctx context.Context, req *pb.ListBookingsRequest) (*pb.ListBookingsResponse, error) {
	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}

	opts := options.Find()
	if req.Page > 0 && req.Limit > 0 {
		opts.SetSkip(int64(req.Page * req.Limit))
		opts.SetLimit(int64(req.Limit))
	}

	cur, err := b.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list bookings")
	}
	defer cur.Close(ctx)

	var bookings []*pb.BookingResponse
	for cur.Next(ctx) {
		var booking pb.BookingResponse
		err := cur.Decode(&booking)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode booking")
		}
		bookings = append(bookings, &booking)
	}

	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor error")
	}

	if len(bookings) == 0 {
		fmt.Println("No bookings found for the given filter and pagination")
	}

	return &pb.ListBookingsResponse{
		Bookings: bookings,
	}, nil
}
