package handler

import (
	pb "Api_Gateway/genproto/booking"
	"Api_Gateway/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// RegisterProvider godoc
// @Summary Registers a new provider
// @Description Registers a new provider and returns the provider details
// @Tags providers
// @Security ApiKeyAuth
// @Param provider body models.Provider true "Provider details"
// @Success 200 {object} booking.ProviderResponse
// @Failure 400 {object} string "Invalid provider data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider [post]
func (h *Handler) RegisterProvider(c *gin.Context) {
	h.Log.Info("RegisterProvider handler is starting")

	Id, ok := c.Get("user_id")
	if !ok {
		er := errors.New("user id not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	var req *models.Provider
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid provider data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Provider.RegisterProvider(ctx, &pb.RegisterProviderRequest{
		UserId:        Id.(string),
		CompanyName:   req.CompanyName,
		Description:   req.Description,
		Services:      req.Services,
		Availability:  req.Availability,
		AverageRating: req.AverageRating,
		Location:      &pb.GeoPoint{Latitude: req.Location.Latitude, Longitude: req.Location.Longitude},
	})
	if err != nil {
		er := errors.Wrap(err, "error registering provider").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("RegisterProvider handler is completed")
	c.JSON(http.StatusOK, gin.H{"provider": resp})
}

// ListProviders godoc
// @Summary Lists all providers with pagination
// @Description Retrieves a list of providers with pagination
// @Tags providers
// @Security ApiKeyAuth
// @Param page query int true "Page number"
// @Param limit query int true "Page limit"
// @Success 200 {object} booking.ListProvidersResponse
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/providers [get]
func (h *Handler) ListProviders(c *gin.Context) {
	h.Log.Info("ListProviders handler is starting")

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

	resp, err := h.Provider.ListProviders(ctx, &pb.ListProvidersRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		er := errors.Wrap(err, "error listing providers").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("ListProviders handler is completed")
	c.JSON(http.StatusOK, gin.H{"providers": resp})
}

// ListProviders1 godoc
// @Summary Lists all providers with pagination
// @Description Retrieves a list of providers with pagination
// @Tags providers
// @Security ApiKeyAuth
// @Param page query int true "Page number"
// @Param limit query int true "Page limit"
// @Success 200 {object} booking.ListProvidersResponse
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider [get]
func (h *Handler) ListProviders1(c *gin.Context) {
	h.Log.Info("ListProviders handler is starting")

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

	resp, err := h.Provider.ListProviders(ctx, &pb.ListProvidersRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		er := errors.Wrap(err, "error listing providers").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("ListProviders handler is completed")
	c.JSON(http.StatusOK, gin.H{"providers": resp})
}

// GetProvider godoc
// @Summary Retrieves a provider by ID
// @Description Retrieves provider details by provider ID
// @Tags providers
// @Security ApiKeyAuth
// @Param id path string true "Provider ID"
// @Success 200 {object} booking.ProviderResponse
// @Failure 400 {object} string "Invalid provider ID"
// @Failure 404 {object} string "Provider not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /customer/providers/{id} [get]
func (h *Handler) GetProvider(c *gin.Context) {
	h.Log.Info("GetProvider handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("provider ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Provider.GetProvider(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting provider").Error()
		if err.Error() == "provider not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetProvider handler is completed")
	c.JSON(http.StatusOK, gin.H{"provider": resp})
}

// GetProvider1 godoc
// @Summary Retrieves a provider by ID
// @Description Retrieves provider details by provider ID
// @Tags providers
// @Security ApiKeyAuth
// @Param id path string true "Provider ID"
// @Success 200 {object} booking.ProviderResponse
// @Failure 400 {object} string "Invalid provider ID"
// @Failure 404 {object} string "Provider not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /admin/providers/{id} [get]
func (h *Handler) GetProvider1(c *gin.Context) {
	h.Log.Info("GetProvider handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("provider ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Provider.GetProvider(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting provider").Error()
		if err.Error() == "provider not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetProvider handler is completed")
	c.JSON(http.StatusOK, gin.H{"provider": resp})
}

// GetProvider2 godoc
// @Summary Retrieves a provider by ID
// @Description Retrieves provider details by provider ID
// @Tags providers
// @Security ApiKeyAuth
// @Param id path string true "Provider ID"
// @Success 200 {object} booking.ProviderResponse
// @Failure 400 {object} string "Invalid provider ID"
// @Failure 404 {object} string "Provider not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider/{id} [get]
func (h *Handler) GetProvider2(c *gin.Context) {
	h.Log.Info("GetProvider handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("provider ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Provider.GetProvider(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting provider").Error()
		if err.Error() == "provider not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetProvider handler is completed")
	c.JSON(http.StatusOK, gin.H{"provider": resp})
}

// UpdateProvider godoc
// @Summary Updates a provider by ID
// @Description Updates provider details by provider ID
// @Tags providers
// @Security ApiKeyAuth
// @Param id path string true "Provider ID"
// @Param provider body booking.UpdateProviderRequest true "Updated provider details"
// @Success 200 {object} booking.ProviderResponse
// @Failure 400 {object} string "Invalid provider data or ID"
// @Failure 404 {object} string "Provider not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider/{id} [put]
func (h *Handler) UpdateProvider(c *gin.Context) {
	h.Log.Info("UpdateProvider handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("provider ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	var req pb.UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid provider data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	req.Id = id

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Provider.UpdateProvider(ctx, &req)
	if err != nil {
		er := errors.Wrap(err, "error updating provider").Error()
		if err.Error() == "provider not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("UpdateProvider handler is completed")
	c.JSON(http.StatusOK, gin.H{"provider": resp})
}

// DeleteProvider godoc
// @Summary Deletes a provider by ID
// @Description Deletes a provider by provider ID
// @Tags providers
// @Security ApiKeyAuth
// @Param id path string true "Provider ID"
// @Success 204 {object} string "No Content"
// @Failure 400 {object} string "Invalid provider ID"
// @Failure 404 {object} string "Provider not found"
// @Failure 500 {object} string "Server error while processing request"
// @Router /provider/{id} [delete]
func (h *Handler) DeleteProvider(c *gin.Context) {
	h.Log.Info("DeleteProvider handler is starting")

	id := c.Param("id")
	if id == "" {
		er := errors.New("provider ID not provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	_, err := h.Provider.DeleteProvider(ctx, &pb.IdRequest{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error deleting provider").Error()
		if err.Error() == "provider not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": er})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		}
		h.Log.Error(er)
		return
	}

	h.Log.Info("DeleteProvider handler is completed")
	c.Status(http.StatusNoContent)
}
