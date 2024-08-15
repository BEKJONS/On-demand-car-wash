package handler

import (
	pb "Api_Gateway/genproto/booking"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// CreateService godoc
// @Summary Creates a new service
// @Description Creates a new service and returns the service details
// @Tags services
// @Security ApiKeyAuth
// @Param service body booking.CreateServiceRequest true "Service details"
// @Success 200 {object} booking.ServiceResponse
// @Failure 400 {object} string "Invalid service data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider/services [post]
func (h *Handler) CreateService(c *gin.Context) {
	h.Log.Info("CreateService handler is starting")

	var req *pb.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid service data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Service.CreateService(ctx, req)
	if err != nil {
		er := errors.Wrap(err, "error creating service").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("CreateService handler is completed")
	c.JSON(http.StatusOK, gin.H{"service": resp})
}

// ListServices godoc
// @Summary Lists all services with pagination
// @Description Retrieves a list of services with pagination
// @Tags services
// @Security ApiKeyAuth
// @Param page query int true "Page number"
// @Param limit query int true "Page limit"
// @Success 200 {object} booking.ListServicesResponse
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/services [get]
func (h *Handler) ListServices(c *gin.Context) {
	h.Log.Info("ListServices handler is starting")

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

	resp, err := h.Service.ListServices(ctx, &pb.ListServicesRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		er := errors.Wrap(err, "error listing services").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("ListServices handler is completed")
	c.JSON(http.StatusOK, gin.H{"services": resp})
}

// ListServices1 godoc
// @Summary Lists all services with pagination
// @Description Retrieves a list of services with pagination
// @Tags services
// @Security ApiKeyAuth
// @Param page query int true "Page number"
// @Param limit query int true "Page limit"
// @Success 200 {object} booking.ListServicesResponse
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider/services [get]
func (h *Handler) ListServices1(c *gin.Context) {
	h.Log.Info("ListServices handler is starting")

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

	resp, err := h.Service.ListServices(ctx, &pb.ListServicesRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		er := errors.Wrap(err, "error listing services").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("ListServices handler is completed")
	c.JSON(http.StatusOK, gin.H{"services": resp})
}

// GetService godoc
// @Summary Retrieves a service by ID
// @Description Retrieves service details by service ID
// @Tags services
// @Security ApiKeyAuth
// @Param id path string true "Service ID"
// @Success 200 {object} booking.ServiceResponse
// @Failure 400 {object} string "Invalid service ID"
// @Failure 404 {object} string "Service not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/services/{id} [get]
func (h *Handler) GetService(c *gin.Context) {
	h.Log.Info("GetService handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("service ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Service.GetServiceByID(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting service").Error()
		if err.Error() == "service not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetService handler is completed")
	c.JSON(http.StatusOK, gin.H{"service": resp})
}

// GetService1 godoc
// @Summary Retrieves a service by ID
// @Description Retrieves service details by service ID
// @Tags services
// @Security ApiKeyAuth
// @Param id path string true "Service ID"
// @Success 200 {object} booking.ServiceResponse
// @Failure 400 {object} string "Invalid service ID"
// @Failure 404 {object} string "Service not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /admin/services/{id} [get]
func (h *Handler) GetService1(c *gin.Context) {
	h.Log.Info("GetService handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("service ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Service.GetServiceByID(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting service").Error()
		if err.Error() == "service not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetService handler is completed")
	c.JSON(http.StatusOK, gin.H{"service": resp})
}

// UpdateService godoc
// @Summary Updates a service by ID
// @Description Updates service details by service ID
// @Tags services
// @Security ApiKeyAuth
// @Param id path string true "Service ID"
// @Param service body booking.UpdateServiceRequest true "Updated service details"
// @Success 200 {object} booking.ServiceResponse
// @Failure 400 {object} string "Invalid service data or ID"
// @Failure 404 {object} string "Service not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider/services/{id} [put]
func (h *Handler) UpdateService(c *gin.Context) {
	h.Log.Info("UpdateService handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("service ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	var req pb.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid service data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	req.Id = id

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Service.UpdateService(ctx, &req)
	if err != nil {
		er := errors.Wrap(err, "error updating service").Error()
		if err.Error() == "service not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("UpdateService handler is completed")
	c.JSON(http.StatusOK, gin.H{"service": resp})
}

// GetPopular godoc
// @Summary Retrieves a popular service
// @Description Retrieves service details
// @Tags services
// @Security ApiKeyAuth
// @Success 200 {object} booking.ListServicesResponse
// @Failure 400 {object} string "Invalid service ID"
// @Failure 404 {object} string "Service not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/services/popular [get]
func (h *Handler) GetPopular(c *gin.Context) {
	h.Log.Info("GetPopular handler is starting")

	res, err := h.Service.PopularServices(c, &pb.Void{})
	if err != nil {
		er := errors.Wrap(err, "error getting popular services").Error()
		if err.Error() == "service not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetService handler is completed")
	c.JSON(http.StatusOK, gin.H{"service": res})
}
