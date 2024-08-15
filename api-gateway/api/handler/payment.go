package handler

import (
	pb "Api_Gateway/genproto/booking"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// CreatePayment godoc
// @Summary Creates a new payment
// @Description Processes a payment and returns the payment details
// @Tags payments
// @Security ApiKeyAuth
// @Param payment body booking.CreatePaymentRequest true "Payment details"
// @Success 200 {object} booking.PaymentResponse
// @Failure 400 {object} string "Invalid payment data"
// @Failure 500 {object} string "Server error while processing payment"
// @Router /customer/payments [post]
func (h *Handler) CreatePayment(c *gin.Context) {
	h.Log.Info("CreatePayment handler is starting")
	Id, ok := c.Get("user_id")
	if !ok {
		er := errors.New("user id not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	var req *pb.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid payment data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	req.UserId = Id.(string)
	//ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	//defer cancel()

	body, err := json.Marshal(req)
	if err != nil {
		er := errors.Wrap(err, "invalid booking data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	err = h.Broker.Payment(body)

	//resp, err := h.Payment.CreatePayment(ctx, req)
	//if err != nil {
	//	er := errors.Wrap(err, "error creating payment").Error()
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
	//	h.Log.Error(er)
	//	return
	//}

	h.Log.Info("CreatePayment handler is completed")
	c.JSON(http.StatusOK, gin.H{"payment": "PAYMENT BY RABBITMQ"})
}

// GetPayment godoc
// @Summary Retrieves a payment by ID
// @Description Retrieves payment details by payment ID
// @Tags payments
// @Security ApiKeyAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} booking.PaymentResponse
// @Failure 400 {object} string "Invalid payment ID"
// @Failure 404 {object} string "Payment not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/payments/{id} [get]
func (h *Handler) GetPayment(c *gin.Context) {
	h.Log.Info("GetPayment handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("payment ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Payment.GetPayment(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting payment").Error()
		if err.Error() == "payment not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetPayment handler is completed")
	c.JSON(http.StatusOK, gin.H{"payment": resp})
}

// GetPayment1 godoc
// @Summary Retrieves a payment by ID
// @Description Retrieves payment details by payment ID
// @Tags payments
// @Security ApiKeyAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} booking.PaymentResponse
// @Failure 400 {object} string "Invalid payment ID"
// @Failure 404 {object} string "Payment not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /admin/payments/{id} [get]
func (h *Handler) GetPayment1(c *gin.Context) {
	h.Log.Info("GetPayment handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("payment ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Payment.GetPayment(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting payment").Error()
		if err.Error() == "payment not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetPayment handler is completed")
	c.JSON(http.StatusOK, gin.H{"payment": resp})
}

// ListPaymentse
// godoc
// @Summary Lists all payments with pagination
// @Description Retrieves a list of payments with pagination
// @Tags payments
// @Security ApiKeyAuth
// @Param page query int true "Page number"
// @Param limit query int true "Page limit"
// @Success 200 {object} booking.ListPaymentsResponse
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/payments [get]
func (h *Handler) ListPaymentse(c *gin.Context) {
	h.Log.Info("ListPayments handler is starting")

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

	resp, err := h.Payment.ListPayments(ctx, &pb.ListPaymentsRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		UserId: c.MustGet("user_id").(string),
	})
	if err != nil {
		er := errors.Wrap(err, "error listing payments").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("ListPayments handler is completed")
	c.JSON(http.StatusOK, gin.H{"payments": resp})
}
