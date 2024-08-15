package main

import (
	"Booking_Service/config"
	pb "Booking_Service/genproto/booking"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net"
	"sync"

	//"product-service/pkg"
	messagebroker "Booking_Service/messagebroker"
	"Booking_Service/pkg/logger"
	"Booking_Service/service"
	mongodb "Booking_Service/storage/mongodb"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	db, err := mongodb.ConnectDB()
	if err != nil {
		log.Fatalf("error while connecting to database: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", cfg.BOOKING_SERVICE_PORT)
	if err != nil {
		log.Fatalf("error while listening: %v", err)
	}
	defer lis.Close()
	loggers := logger.NewLogger()
	b := service.NewBookingService(db, loggers)
	s := service.NewManagementService(db, loggers)
	p := service.NewPaymentService(db, loggers)
	r := service.NewReviewService(db, loggers)
	pr := service.NewProviderService(db, loggers)
	sr := service.NewSearchService(db, loggers)
	n := service.NewNotificationService(db)
	server := grpc.NewServer()

	pb.RegisterBookingServiceServer(server, b)
	pb.RegisterServiceManagementServiceServer(server, s)
	pb.RegisterPaymentServiceServer(server, p)
	pb.RegisterReviewServiceServer(server, r)
	pb.RegisterProviderServiceServer(server, pr)
	pb.RegisterSearchingServiceServer(server, sr)
	pb.RegisterNotificationsServer(server, n)

	conn, err := amqp.Dial(cfg.RABBIT)
	if err != nil {
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------", err, "-------------------------------------------------------------------------------------")
		panic(err)
	}
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
	}
	defer ch.Close()

	servis := grpc.NewServer()

	bService := service.NewBookingService(db, loggers)
	pb.RegisterBookingServiceServer(servis, bService)

	pService := service.NewPaymentService(db, loggers)
	pb.RegisterPaymentServiceServer(servis, pService)

	rService := service.NewReviewService(db, loggers)
	pb.RegisterReviewServiceServer(servis, rService)

	nService := service.NewNotificationService(db)
	pb.RegisterNotificationsServer(servis, nService)

	notQueue, err := getQueue(ch, "create_notification")
	if err != nil {
		log.Println(err)
	}

	crtQueue, err := getQueue(ch, "create_booking")
	if err != nil {
		log.Println(err)
	}

	delQueue, err := getQueue(ch, "booking_cancelled")
	if err != nil {
		log.Println(err)
	}

	payQueue, err := getQueue(ch, "payment_processed")
	if err != nil {
		log.Println(err)
	}

	revQueue, err := getQueue(ch, "review_submitted")
	if err != nil {
		log.Println(err)
	}

	notMsg, err := getMessageQueue(ch, notQueue)
	if err != nil {
		log.Println(err)
	}

	crtMsg, err := getMessageQueue(ch, crtQueue)
	if err != nil {
		log.Println(err)
	}

	delMsg, err := getMessageQueue(ch, delQueue)
	if err != nil {
		log.Println(err)
	}

	payMsg, err := getMessageQueue(ch, payQueue)
	if err != nil {
		log.Println(err)
	}

	revMsg, err := getMessageQueue(ch, revQueue)
	if err != nil {
		log.Println(err)
	}

	res := messagebroker.New(bService, pService, rService, nService, ch, loggers, crtMsg, delMsg, payMsg, revMsg, notMsg, &sync.WaitGroup{}, 6)
	go res.StartToConsume(context.Background())

	log.Printf("Service is listening on port %s...\n", cfg.BOOKING_SERVICE_PORT)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("error while serving product service: %s", err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
}

func getMessageQueue(ch *amqp.Channel, q amqp.Queue) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		q.Name,
		"",
		false, false, false, false, nil,
	)
}
