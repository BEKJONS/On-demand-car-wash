package models

type NewNotification struct {
	UserID    string `bson:"user_id"`
	Title     string `bson:"title"`
	Message   string `bson:"message"`
	CreatedAt string `bson:"created_at"`
}

type Notification struct {
	Id        string `bson:"_id"`
	UserID    string `bson:"user_id"`
	Title     string `bson:"title"`
	Message   string `bson:"message"`
	CreatedAt string `bson:"created_at"`
}
