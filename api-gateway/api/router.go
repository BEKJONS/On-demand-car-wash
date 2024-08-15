package api

import (
	"Api_Gateway/api/handler"
	"Api_Gateway/api/middleware"
	"github.com/casbin/casbin/v2"
	amqp "github.com/rabbitmq/amqp091-go"

	"Api_Gateway/config"

	_ "Api_Gateway/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title On-Demand Car Wash Service
// @version 1.0
// @description API Gateway of On-Demand Car Wash Service
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(cfg *config.Config, casbin *casbin.Enforcer, conn *amqp.Channel) *gin.Engine {
	h := handler.NewHandler(cfg, conn)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	api.Use(middleware.PermissionMiddleware(casbin))

	// Admin routes
	admin := api.Group("/admin")
	{
		u := admin.Group("/user")
		{
			u.GET("/:id", h.GetUser)
			u.PUT("/", h.UpdateUser)
			u.DELETE("/:id", h.DeleteUser)
		}
		admin.GET("/users", h.FetchUsers)
		admin.GET("/bookings/all/:id", h.ListBookings)
		admin.GET("/providers/:id", h.GetProvider1)
		admin.GET("/services/:id", h.GetService1)
		admin.GET("/reviews/:id", h.GetReview1)
		admin.GET("/payments/:id", h.GetPayment1)
	}

	// Customer routes
	customer := api.Group("/customer")
	{
		p := customer.Group("/profile")
		{
			p.GET("/", h.GetProfile)
			p.PUT("/", h.UpdateProfile)
		}

		b := customer.Group("/bookings")
		{
			b.POST("", h.CreateBooking)
			b.GET("/:id", h.GetBooking)
			b.PUT("/:id", h.UpdateBooking)
			b.DELETE("/:id", h.CancelBooking)
			b.GET("/all", h.ListBookings1)
		}
		tr := customer.Group("/payments")
		{
			tr.POST("/", h.CreatePayment)
			tr.GET("/:id", h.GetPayment)
			tr.GET("/", h.ListPaymentse)
		}
		r := customer.Group("/reviews")
		{
			r.POST("/", h.CreateReview)
			r.GET("/:id", h.GetReview)
			r.PUT("/:id", h.UpdateReview)
			r.GET("/", h.ListReviews)
			r.DELETE("/:id", h.DeleteReview)
		}
		s := customer.Group("/services")
		{
			s.GET("/", h.ListServices)
			s.GET("/:id", h.GetService)
			s.GET("popular", h.GetPopular)
		}
		sr := customer.Group("/search")
		{
			sr.GET("/provider", h.SearchProviders)
			sr.GET("/service", h.SearchServices)
		}
	}

	// Provider routes
	provider := api.Group("/provider")
	{
		provider.POST("/", h.RegisterProvider)
		provider.GET("/", h.ListProviders1)
		provider.GET("/:id", h.GetProvider2)
		provider.PUT("/:id", h.UpdateProvider)
		provider.DELETE("/:id", h.DeleteProvider)
		provider.GET("/services", h.ListServices1)
		provider.PUT("/services/:id", h.UpdateService)
		provider.POST("/services", h.CreateService)
		provider.POST("/notifications", h.CreateNotification)
		provider.GET("/notifications/:id", h.GetNotification)
	}

	return router
}
