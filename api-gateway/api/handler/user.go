package handler

import (
	pb "Api_Gateway/genproto/users"
	"Api_Gateway/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// GetProfile godoc
// @Summary Gets user
// @Description Retrieves user profile by ID
// @Tags user
// @Security ApiKeyAuth
// @Success 200 {object} users.GetProfileResponse
// @Failure 400 {object} string "Invalid user id"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	h.Log.Info("GetProfile handler is starting")

	id, ok := c.Get("user_id")
	if !ok {
		er := errors.New("user id not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()
	pr, err := h.User.GetProfileU(ctx, &pb.Id{UserId: id.(string)})
	if err != nil {
		er := errors.Wrap(err, "error getting user").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetProfile handler is completed")
	c.JSON(http.StatusOK, gin.H{"user": pr})
}

// UpdateProfile godoc
// @Summary Updates user
// @Description Updates user info by ID
// @Tags user
// @Security ApiKeyAuth
// @Param user body models.UUser true "New user data"
// @Success 200 {object} users.UserResponseU
// @Failure 400 {object} string "Invalid user id or data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	h.Log.Info("UpdateProfile handler is starting")

	id, ok := c.Get("user_id")
	if !ok {
		er := errors.New("user id not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	var data models.UUser
	err := c.ShouldBind(&data)
	if err != nil {
		er := errors.Wrap(err, "invalid user data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.User.UpdateProfile(ctx, &pb.UpdateProfileRequestU{
		UserId:      id.(string),
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		PhoneNumber: data.PhoneNumber,
	})
	if err != nil {
		er := errors.Wrap(err, "error updating user").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	h.Log.Info("UpdateProfile handler is completed")
	c.JSON(http.StatusOK, gin.H{"user": resp})
}
