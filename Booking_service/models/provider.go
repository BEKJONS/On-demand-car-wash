package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Provider represents a provider entity in the booking service.
type Provider struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        string             `bson:"user_id"`
	CompanyName   string             `bson:"company_name"`
	Description   string             `bson:"description"`
	Services      []string           `bson:"services"`
	Availability  string             `bson:"availability"`
	AverageRating float64            `bson:"average_rating"`
	Location      *GeoPoint          `bson:"location"`
	CreatedAt     string             `bson:"created_at"`
	UpdatedAt     string             `bson:"updated_at"`
}

type Provider1 struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        string             `bson:"user_id"`
	CompanyName   string             `bson:"company_name"`
	Description   string             `bson:"description"`
	Services      []string           `bson:"services"`
	Availability  string             `bson:"availability"`
	AverageRating float64            `bson:"average_rating"`
	Location      string             `bson:"location"`
	CreatedAt     string             `bson:"created_at"`
	UpdatedAt     string             `bson:"updated_at"`
}
