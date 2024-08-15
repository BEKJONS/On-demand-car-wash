package mongodb

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/models"
	"Booking_Service/storage"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type PaymentRepo struct {
	collection *mongo.Collection
}

func NewPaymentRepo(db *mongo.Database) storage.IPaymentStorage {
	return &PaymentRepo{
		collection: db.Collection("Payments"),
	}
}

func (r *PaymentRepo) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.PaymentResponse, error) {
	// Create a new payment instance
	payment := &models.Payment{
		UserID:        req.UserId,
		BookingID:     req.BookingId,
		Amount:        req.Amount,
		Status:        "pending", // Default status
		PaymentMethod: req.PaymentMethod,
		TransactionID: req.TransactionId,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	// Insert the payment into MongoDB
	res, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}

	// Convert the inserted ID to a hex string
	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert inserted ID to ObjectID")
	}

	// Return the payment response
	return &pb.PaymentResponse{
		Id:            objID.Hex(),
		BookingId:     payment.BookingID,
		Amount:        payment.Amount,
		Status:        payment.Status,
		PaymentMethod: payment.PaymentMethod,
		TransactionId: payment.TransactionID,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}, nil
}

func (r *PaymentRepo) GetPayment(ctx context.Context, id *pb.IdRequest) (*pb.PaymentResponse, error) {
	// Convert the hex string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return nil, errors.Wrap(err, "invalid payment ID")
	}

	// Find the payment by ID
	var payment models.Payment
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("payment not found")
		}
		return nil, errors.Wrap(err, "failed to find payment")
	}

	// Return the payment response
	return &pb.PaymentResponse{
		Id:            objID.Hex(),
		BookingId:     payment.BookingID,
		Amount:        payment.Amount,
		Status:        payment.Status,
		PaymentMethod: payment.PaymentMethod,
		TransactionId: payment.TransactionID,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}, nil
}

func (r *PaymentRepo) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	opts := options.Find().
		SetSkip(int64((req.Page - 1) * req.Limit)).
		SetLimit(int64(req.Limit))

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": req.UserId}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []*pb.PaymentResponse
	for cursor.Next(ctx) {
		var payment pb.PaymentResponse
		if err := cursor.Decode(&payment); err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.ListPaymentsResponse{Payments: payments}, nil
}
