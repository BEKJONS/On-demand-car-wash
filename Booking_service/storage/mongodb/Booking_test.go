package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	pb "Booking_Service/genproto/booking"
	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) (*mongo.Client, *mongo.Collection, func()) {
	// Initialize mtest
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	var client *mongo.Client
	var db *mongo.Database
	var col *mongo.Collection

	mt.Run("SetupTestDB", func(mt *mtest.T) {
		client = mt.Client
		db = client.Database("booking_test")
		col = db.Collection("bookings")
	})

	// Cleanup function to drop the database
	cleanup := func() {
		_ = db.Drop(context.Background()) // Drop the test database
	}

	return client, col, cleanup
}

func TestAddBooking(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	// Define the mock responses for the services collection
	mt.Run("SetupTestDB", func(mt *mtest.T) {
		client := mt.Client
		//db := client.Database("booking_test")
		//col := db.Collection("services")

		// Mock response for FindOne
		serviceID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "booking_test.services", mtest.FirstBatch, bson.D{
				{"_id", serviceID},
				{"name", "Sample Service"},
				{"description", "A sample service"},
				{"price", 50.0},
				{"duration", 60},
			}),
		)

		// Mock response for InsertOne
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(), // Successful InsertOne response
		)

		repo := NewBookingRepo(client.Database("booking_test"))
		req := &pb.CreateBookingRequest{
			UserId:        "user123",
			ProviderId:    "provider123",
			ServiceId:     serviceID.Hex(),
			ScheduledTime: time.Now().Format(time.RFC3339),
			Location:      &pb.GeoPoint{Latitude: 40.7128, Longitude: -74.0060},
		}

		resp, err := repo.Add(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, req.UserId, resp.UserId)
	})

}

func TestReadBooking(t *testing.T) {
	client, _, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewBookingRepo(client.Database("booking_test"))
	booking := &pb.CreateBookingRequest{
		UserId:        "user123",
		ProviderId:    "provider123",
		ServiceId:     primitive.NewObjectID().Hex(),
		ScheduledTime: time.Now().Format(time.RFC3339),
		Location:      &pb.GeoPoint{Latitude: 40.7128, Longitude: -74.0060},
	}

	resp, err := repo.Add(context.Background(), booking)
	require.NoError(t, err)
	require.NotNil(t, resp)

	readResp, err := repo.Read(context.Background(), &pb.IdRequest{Id: resp.Id})
	require.NoError(t, err)
	require.NotNil(t, readResp)
	require.Equal(t, resp.UserId, readResp.UserId)
}

func TestUpdateBooking(t *testing.T) {
	client, _, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewBookingRepo(client.Database("booking_test"))
	booking := &pb.CreateBookingRequest{
		UserId:        "user123",
		ProviderId:    "provider123",
		ServiceId:     primitive.NewObjectID().Hex(),
		ScheduledTime: time.Now().Format(time.RFC3339),
		Location:      &pb.GeoPoint{Latitude: 40.7128, Longitude: -74.0060},
	}

	resp, err := repo.Add(context.Background(), booking)
	require.NoError(t, err)
	require.NotNil(t, resp)

	updateReq := &pb.UpdateBookingRequest{
		Id:        resp.Id,
		ServiceId: primitive.NewObjectID().Hex(),
	}

	updateResp, err := repo.Update(context.Background(), updateReq)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, updateReq.ServiceId, updateResp.ServiceId)
}

func TestCancelBooking(t *testing.T) {
	client, _, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewBookingRepo(client.Database("booking_test"))
	booking := &pb.CreateBookingRequest{
		UserId:        "user123",
		ProviderId:    "provider123",
		ServiceId:     primitive.NewObjectID().Hex(),
		ScheduledTime: time.Now().Format(time.RFC3339),
		Location:      &pb.GeoPoint{Latitude: 40.7128, Longitude: -74.0060},
	}

	resp, err := repo.Add(context.Background(), booking)
	require.NoError(t, err)
	require.NotNil(t, resp)

	cancelResp, err := repo.Cancel(context.Background(), &pb.IdRequest{Id: resp.Id})
	require.NoError(t, err)
	require.NotNil(t, cancelResp)
	require.Equal(t, "canceled", cancelResp.Status)
}

func TestListBookings(t *testing.T) {
	client, _, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewBookingRepo(client.Database("booking_test"))

	for i := 0; i < 5; i++ {
		booking := &pb.CreateBookingRequest{
			UserId:        "user" + primitive.NewObjectID().Hex(),
			ProviderId:    "provider" + primitive.NewObjectID().Hex(),
			ServiceId:     primitive.NewObjectID().Hex(),
			ScheduledTime: time.Now().Format(time.RFC3339),
			Location:      &pb.GeoPoint{Latitude: 40.7128, Longitude: -74.0060},
		}

		_, err := repo.Add(context.Background(), booking)
		require.NoError(t, err)
	}

	listReq := &pb.ListBookingsRequest{
		UserId: "user123",
		Limit:  5,
		Page:   0,
	}

	listResp, err := repo.List(context.Background(), listReq)
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Bookings, 5)
}
