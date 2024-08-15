package mongodb

import (
	"Booking_Service/config"
	"Booking_Service/storage"
	redisDB "Booking_Service/storage/redis"
	"github.com/redis/go-redis/v9"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type str string

type Storage struct {
	mongo *mongo.Database
	redis *redis.Client
}

func ConnectDB() (storage.IStorage, error) {
	opts := options.Client().ApplyURI(config.Load().MongoURI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	strg := Storage{
		mongo: client.Database(config.Load().MongoDB_NAME),
		redis: redisDB.ConnectRedis(),
	}

	return strg, nil
}

func (s Storage) Close() {
	s.mongo.Client().Disconnect(context.Background())
	s.redis.Close()
}

func (s Storage) Booking() storage.IBookingStorage {
	return NewBookingRepo(s.mongo, s.redis)
}

func (s Storage) Management() storage.IManagementStorage {
	return NewManagementRepo(s.mongo)
}

func (s Storage) Payment() storage.IPaymentStorage {
	return NewPaymentRepo(s.mongo)
}

func (s Storage) Review() storage.IReviewStorage {
	return NewReviewRepo(s.mongo)
}

func (s Storage) Provider() storage.IProviderStorage {
	return NewProviderRepo(s.mongo)
}

func (s Storage) Search() storage.ISearchStorage {
	return NewSearchingRepo(s.mongo)
}

func (s Storage) Notification() storage.INotificationStorage { return NewNotificationRepo(s.mongo) }
