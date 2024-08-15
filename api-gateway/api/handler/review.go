package handler

import (
	pb "Api_Gateway/genproto/booking"
	"Api_Gateway/models"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// CreateReview godoc
// @Summary Creates a new review
// @Description Submits a review for a booking and returns the review details
// @Tags reviews
// @Security ApiKeyAuth
// @Param review body models.ReviewRequest true "Review details"
// @Success 200 {object} booking.ReviewResponse
// @Failure 400 {object} string "Invalid review data"
// @Failure 500 {object} string "Server error while processing review"
// @Router /customer/reviews [post]
func (h *Handler) CreateReview(c *gin.Context) {
	h.Log.Info("CreateReview handler is starting")

	Id, ok := c.Get("user_id")
	if !ok {
		er := errors.New("user id not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	var req *models.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid review data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	req1 := &pb.CreateReviewRequest{
		BookingId:  req.BookingID,
		UserId:     Id.(string),
		ProviderId: req.ProviderID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	body, err := json.Marshal(req1)
	if err != nil {
		er := errors.Wrap(err, "invalid booking data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	err = h.Broker.Review(body)

	//ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	//defer cancel()
	//
	//resp, err := h.Review.CreateReview(ctx, &pb.CreateReviewRequest{
	//  BookingId:  req.BookingID,
	//  UserId:     Id.(string),
	//  ProviderId: req.ProviderID,
	//  Rating:     req.Rating,
	//  Comment:    req.Comment,
	//})
	//if err != nil {
	//  er := errors.Wrap(err, "error creating review").Error()
	//  c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
	//  h.Log.Error(er)
	//  return
	//}

	h.Log.Info("CreateReview handler is completed")
	c.JSON(http.StatusOK, gin.H{"review": req.ProviderID})
}

// ListReviews godoc
// @Summary Lists all reviews with pagination
// @Description Retrieves a list of reviews with pagination
// @Tags reviews
// @Security ApiKeyAuth
// @Param provider_id query string true "Provider ID"
// @Param page query int true "Page number"
// @Param limit query int true "Page limit"
// @Success 200 {object} booking.ListReviewsResponse
// @Failure 400 {object} string "Invalid request parameters"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/reviews [get]
func (h *Handler) ListReviews(c *gin.Context) {
	h.Log.Info("ListReviews handler is starting")

	providerId := c.Query("provider_id")
	if providerId == "" {
		er := "provider_id is required"
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

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

	resp, err := h.Review.ListReviews(ctx, &pb.ListReviewsRequest{
		Page:       int32(page),
		Limit:      int32(limit),
		ProviderId: providerId,
	})
	if err != nil {
		er := errors.Wrap(err, "error listing reviews").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("ListReviews handler is completed")
	c.JSON(http.StatusOK, gin.H{"reviews": resp})
}

// GetReview godoc
// @Summary Retrieves a review by ID
// @Description Retrieves review details by review ID
// @Tags reviews
// @Security ApiKeyAuth
// @Param id path string true "Review ID"
// @Success 200 {object} booking.ReviewResponse
// @Failure 400 {object} string "Invalid review ID"
// @Failure 404 {object} string "Review not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/reviews/{id} [get]
func (h *Handler) GetReview(c *gin.Context) {
	h.Log.Info("GetReview handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("review ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Review.GetReviewById(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting review").Error()
		if err.Error() == "review not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetReview handler is completed")
	c.JSON(http.StatusOK, gin.H{"review": resp})
}

// GetReview1 godoc
// @Summary Retrieves a review by ID
// @Description Retrieves review details by review ID
// @Tags reviews
// @Security ApiKeyAuth
// @Param id path string true "Review ID"
// @Success 200 {object} booking.ReviewResponse
// @Failure 400 {object} string "Invalid review ID"
// @Failure 404 {object} string "Review not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /admin/reviews/{id} [get]
func (h *Handler) GetReview1(c *gin.Context) {
	h.Log.Info("GetReview handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("review ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Review.GetReviewById(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting review").Error()
		if err.Error() == "review not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetReview handler is completed")
	c.JSON(http.StatusOK, gin.H{"review": resp})
}

// UpdateReview godoc
// @Summary Updates a review by ID
// @Description Updates the details of a review by its ID
// @Tags reviews
// @Security ApiKeyAuth
// @Param id path string true "Review ID"
// @Param review body models.Review true "Review details"
// @Success 200 {object} booking.ReviewResponse
// @Failure 400 {object} string "Invalid review data"
// @Failure 404 {object} string "Review not found"
// @Failure 500 {object} string "Server error while processing review"
// @Router /customer/reviews/{id} [put]
func (h *Handler) UpdateReview(c *gin.Context) {
	h.Log.Info("UpdateReview handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("review ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	var req *models.Review
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid review data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Review.UpdateReview(ctx, &pb.UpdateReviewRequest{
		Id:      id,
		Rating:  req.Rating,
		Comment: req.Comment,
	})
	if err != nil {
		er := errors.Wrap(err, "error updating review").Error()
		if err.Error() == "review not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("UpdateReview handler is completed")
	c.JSON(http.StatusOK, gin.H{"review": resp})
}

// DeleteReview godoc
// @Summary Deletes a review by ID
// @Description Deletes a review by its ID
// @Tags reviews
// @Security ApiKeyAuth
// @Param id path string true "Review ID"
// @Success 200 {object} string "Review successfully deleted"
// @Failure 400 {object} string "Invalid review ID"
// @Failure 404 {object} string "Review not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/reviews/{id} [delete]
func (h *Handler) DeleteReview(c *gin.Context) {
	h.Log.Info("DeleteReview handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("review ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	_, err := h.Review.DeleteReview(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error deleting review").Error()
		if err.Error() == "review not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("DeleteReview handler is completed")
	c.JSON(http.StatusOK, gin.H{"message": "Review successfully deleted"})
}
