package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review represents a review record in the MongoDB collection
type Review struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BookingID  string             `bson:"booking_id" json:"booking_id"`
	UserID     string             `bson:"user_id" json:"user_id"`
	ProviderID string             `bson:"provider_id" json:"provider_id"`
	Rating     float32            `bson:"rating" json:"rating"`
	Comment    string             `bson:"comment" json:"comment"`
	CreatedAt  string             `bson:"created_at" json:"created_at"`
	UpdatedAt  string             `bson:"updated_at" json:"updated_at"`
}
