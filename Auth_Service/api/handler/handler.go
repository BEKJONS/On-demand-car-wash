package handler

import (
	"Auth/pkg/logger"
	"Auth/service"
	"Auth/storage"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	User  *service.UserService
	Log   *slog.Logger
	Admin *service.AdminService
	Token storage.ITokenStorage
}

func NewHandler(s storage.IStorage) *Handler {
	return &Handler{
		User:  service.NewUserService(s),
		Admin: service.NewAdminService(s),
		Token: s.Token(),
		Log:   logger.NewLogger(),
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
