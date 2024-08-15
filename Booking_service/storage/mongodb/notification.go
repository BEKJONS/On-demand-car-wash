package mongodb

import (
	"Booking_Service/models"
	"Booking_Service/storage"
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepo struct {
	col *mongo.Collection
}

func NewNotificationRepo(db *mongo.Database) storage.INotificationStorage {
	return &NotificationRepo{col: db.Collection("notifications")}
}

func (r *NotificationRepo) Create(ctx context.Context, req *models.NewNotification) (string, error) {
	res, err := r.col.InsertOne(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "query execution failed")
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("id conversion failed")
	}
	return id.Hex(), nil
}

func (r *NotificationRepo) Get(ctx context.Context, id string) (*models.Notification, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "invalid id")
	}

	res := r.col.FindOne(ctx, bson.M{"_id": objId})
	if res.Err() != nil {
		return nil, errors.Wrap(res.Err(), "query execution failed")
	}

	var n models.Notification
	if err := res.Decode(&n); err != nil {
		return nil, errors.Wrap(err, "decoding failed")
	}

	return &n, nil
}
