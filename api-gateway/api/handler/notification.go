package handler

import (
	pbn "Api_Gateway/genproto/booking"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// CreateNotification godoc
// @Summary Creates notification
// @Description Adds a new notification
// @Tags notification
// @Security ApiKeyAuth
// @Param data body booking.NewNotification true "Receiver ID, Title and Message"
// @Success 201 {object} string "Notification created"
// @Failure 400 {object} string "Invalid data format"
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider/notifications [post]
func (h *Handler) CreateNotification(c *gin.Context) {
	h.Log.Info("CreateNotification handler is starting")

	var req pbn.NewNotification
	if err := c.ShouldBind(&req); err != nil {
		er := errors.Wrap(err, "invalid data format").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	message, err := json.Marshal(&req)
	if err != nil {
		er := errors.Wrap(err, "error serializing notification").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	err = h.Broker.CreateNotification(message)
	if err != nil {
		er := errors.Wrap(err, "error creating notification").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("CreateNotification handler is completed")
	c.JSON(http.StatusCreated, "Notification created")
}

// GetNotification godoc
// @Summary Gets notification
// @Description Gets notification by ID
// @Tags notification
// @Security ApiKeyAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} booking.Notification
// @Failure 400 {object} string "Invalid data format"
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider/notifications/{id} [get]
func (h *Handler) GetNotification(c *gin.Context) {
	h.Log.Info("GetNotification handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("invalid data format").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Notification.GetNotification(ctx, &pbn.ID{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error finding notification").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetNotification handler is completed")
	c.JSON(http.StatusOK, resp)
}
