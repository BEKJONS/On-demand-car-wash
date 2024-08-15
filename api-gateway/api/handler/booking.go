package handler

import (
	pb "Api_Gateway/genproto/booking"
	"Api_Gateway/models"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// CreateBooking godoc
// @Summary Creates a new booking
// @Description Creates a new booking and returns the booking details
// @Tags bookings
// @Security ApiKeyAuth
// @Param booking body models.Booking true "Booking details"
// @Success 200 {object} booking.BookingResponse
// @Failure 400 {object} string "Invalid booking data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/bookings [post]
func (h *Handler) CreateBooking(c *gin.Context) {
	h.Log.Info("CreateBooking handler is starting")

	var req *models.Booking
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid booking data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	id, ok := c.Get("user_id")
	if !ok {
		er := errors.New("user id not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	req1 := &pb.CreateBookingRequest{
		UserId:        id.(string),
		ProviderId:    req.ProviderID,
		ServiceId:     req.ServiceID,
		ScheduledTime: req.ScheduledTime,
		Location:      &pb.GeoPoint{Longitude: req.Location.Longitude, Latitude: req.Location.Latitude},
	}

	body, err := json.Marshal(req1)
	if err != nil {
		er := errors.Wrap(err, "invalid booking data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	err = h.Broker.CreateBooking(body)
	if err != nil {
		er := errors.Wrap(err, "invalid booking data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	//ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	//defer cancel()

	//resp, err := h.Booking.CreateBooking(ctx, &pb.CreateBookingRequest{
	//	UserId:        id.(string),
	//	ProviderId:    req.ProviderID,
	//	ServiceId:     req.ServiceID,
	//	ScheduledTime: req.ScheduledTime,
	//	Location:      &pb.GeoPoint{Longitude: req.Location.Longitude, Latitude: req.Location.Latitude},
	//})
	//if err != nil {
	//	er := errors.Wrap(err, "error creating booking").Error()
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
	//	h.Log.Error(er)
	//	return
	//}

	h.Log.Info("CreateBooking handler is completed")
	c.JSON(http.StatusOK, gin.H{"booking": "CREATED BY RABBITMQ"})
}

// GetBooking godoc
// @Summary Retrieves a booking by ID
// @Description Retrieves booking details by booking ID
// @Tags bookings
// @Security ApiKeyAuth
// @Param id path string true "Booking ID"
// @Success 200 {object} booking.BookingResponse
// @Failure 400 {object} string "Invalid booking ID"
// @Failure 404 {object} string "Booking not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/bookings/{id} [get]
func (h *Handler) GetBooking(c *gin.Context) {
	h.Log.Info("GetBooking handler is starting")

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Booking.GetBooking(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting booking").Error()
		if err.Error() == "booking not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetBooking handler is completed")
	c.JSON(http.StatusOK, gin.H{"booking": resp})
}

// UpdateBooking godoc
// @Summary Updates a booking by ID
// @Description Updates booking details by booking ID
// @Tags bookings
// @Security ApiKeyAuth
// @Param id path string true "Booking ID"
// @Param booking body booking.UpdateBookingRequest true "Updated booking details"
// @Success 200 {object} booking.BookingResponse
// @Failure 400 {object} string "Invalid booking data or ID"
// @Failure 404 {object} string "Booking not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/bookings/{id} [put]
func (h *Handler) UpdateBooking(c *gin.Context) {
	h.Log.Info("UpdateBooking handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("booking ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	var req pb.UpdateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid booking data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	req.Id = id

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Booking.UpdateBooking(ctx, &req)
	if err != nil {
		er := errors.Wrap(err, "error updating booking").Error()
		if err.Error() == "booking not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("UpdateBooking handler is completed")
	c.JSON(http.StatusOK, gin.H{"booking": resp})
}

// CancelBooking godoc
// @Summary Cancels a booking by ID
// @Description Cancels booking by booking ID
// @Tags bookings
// @Security ApiKeyAuth
// @Param id path string true "Booking ID"
// @Success 200 {object} booking.BookingResponse
// @Failure 400 {object} string "Invalid booking ID"
// @Failure 404 {object} string "Booking not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/bookings/{id} [delete]
func (h *Handler) CancelBooking(c *gin.Context) {
	h.Log.Info("CancelBooking handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("booking ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	req1 := &pb.IdRequest{Id: id}
	body, err := json.Marshal(req1)
	if err != nil {
		er := errors.Wrap(err, "invalid booking data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	err = h.Broker.CancelBooking(body)

	//ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	//defer cancel()
	//
	//resp, err := h.Booking.CancelBooking(ctx, &pb.IdRequest{Id: id})
	//if err != nil {
	//	er := errors.Wrap(err, "error canceling booking").Error()
	//	if err.Error() == "booking not found" {
	//		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
	//	} else {
	//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
	//	}
	//	h.Log.Error(er)
	//	return
	//}

	h.Log.Info("CancelBooking handler is completed")
	c.JSON(http.StatusOK, gin.H{"booking": "CANCELED BY RABBITMQ"})
}

// ListBookings godoc
// @Summary Lists bookings with pagination
// @Description Retrieves a list of bookings with pagination
// @Tags bookings
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param page query int true "Page number"
// @Param limit query int true "Page limit"
// @Success 200 {object} booking.ListBookingsResponse
// @Failure 500 {object} string "Server error while processing request"
// @Router /admin/bookings/all/:id [get]
func (h *Handler) ListBookings(c *gin.Context) {
	h.Log.Info("ListBookings handler is starting")

	id := c.Param("user_id")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		er := "invalid page number"
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		er := "invalid limit value"
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()
	resp, err := h.Booking.ListBookings(ctx, &pb.ListBookingsRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		UserId: id,
	})
	if err != nil {
		er := errors.Wrap(err, "error listing bookings").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("ListBookings handler is completed")
	c.JSON(http.StatusOK, gin.H{"bookings": resp})
}

// ListBookings1 godoc
// @Summary Lists bookings with pagination
// @Description Retrieves a list of bookings with pagination
// @Tags bookings
// @Security ApiKeyAuth
// @Param page query int true "Page number"
// @Param limit query int true "Page limit"
// @Success 200 {object} booking.ListBookingsResponse
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/bookings/all [get]
func (h *Handler) ListBookings1(c *gin.Context) {
	h.Log.Info("ListBookings handler is starting")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		er := "invalid page number"
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		er := "invalid limit value"
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()
	resp, err := h.Booking.ListBookings(ctx, &pb.ListBookingsRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		UserId: c.MustGet("user_id").(string),
	})
	if err != nil {
		er := errors.Wrap(err, "error listing bookings").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("ListBookings handler is completed")
	c.JSON(http.StatusOK, gin.H{"bookings": resp})
}
