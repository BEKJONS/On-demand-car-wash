package main

import (
	"Api_Gateway/api"
	"Api_Gateway/config"
	"github.com/casbin/casbin/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	cfg := config.Load()

	CasbinEnforcer, err := casbin.NewEnforcer("./config/model.conf", "./config/policy.csv")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	time.Sleep(10 * time.Second)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	router := api.NewRouter(cfg, CasbinEnforcer, ch)

	router.Run(cfg.HTTP_PORT)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
