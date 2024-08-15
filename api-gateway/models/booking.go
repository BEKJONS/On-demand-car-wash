package models

type Booking struct {
	ProviderID    string   `bson:"provider_id" json:"provider_id"`
	ServiceID     string   `bson:"service_id" json:"service_id"`
	ScheduledTime string   `bson:"scheduled_time" json:"scheduled_time"`
	Location      GeoPoint `bson:"location" json:"location"`
}
type GeoPoint struct {
	Longitude float64 `bson:"longitude" json:"longitude"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
}

type Review struct {
	Rating  float32 `bson:"rating" json:"rating"`
	Comment string  `bson:"comment" json:"comment"`
}
type ReviewRequest struct {
	BookingID  string  `bson:"booking_id" json:"booking_id"`
	ProviderID string  `bson:"provider_id" json:"provider_id"`
	Rating     float32 `bson:"rating" json:"rating"`
	Comment    string  `bson:"comment" json:"comment"`
}

type Provider struct {
	CompanyName   string    `protobuf:"bytes,2,opt,name=company_name,json=companyName,proto3" json:"company_name,omitempty"`
	Description   string    `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Services      []string  `protobuf:"bytes,4,rep,name=services,proto3" json:"services,omitempty"`
	Availability  string    `protobuf:"bytes,5,rep,name=availability,proto3" json:"availability,omitempty"`
	AverageRating float64   `protobuf:"fixed64,6,opt,name=average_rating,json=averageRating,proto3" json:"average_rating,omitempty"`
	Location      *GeoPoint `protobuf:"bytes,7,opt,name=location,proto3" json:"location,omitempty"`
}
