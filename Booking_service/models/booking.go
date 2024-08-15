package models

type Booking struct {
	UserID        string   `bson:"user_id"`
	ProviderID    string   `bson:"provider_id"`
	ServiceID     string   `bson:"service_id"`
	Status        string   `bson:"status"`
	ScheduledTime string   `bson:"scheduled_time"`
	Location      GeoPoint `bson:"location"`
	TotalPrice    float64  `bson:"total_price"`
	CreatedAt     string   `bson:"created_at"`
	UpdatedAt     string   `bson:"updated_at"`
}

type GeoPoint struct {
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
}
