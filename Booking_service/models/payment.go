package models

// Payment represents a payment record in the MongoDB collection
type Payment struct {
	ID            string  `bson:"_id,omitempty" json:"id,omitempty"`
	UserID        string  `bson:"userId,omitempty" json:"userId,omitempty"`
	BookingID     string  `bson:"booking_id" json:"booking_id"`
	Amount        float64 `bson:"amount" json:"amount"`
	Status        string  `bson:"status" json:"status"`
	PaymentMethod string  `bson:"payment_method" json:"payment_method"`
	TransactionID string  `bson:"transaction_id" json:"transaction_id"`
	CreatedAt     string  `bson:"created_at" json:"created_at"`
	UpdatedAt     string  `bson:"updated_at" json:"updated_at"`
}
