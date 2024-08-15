package handler

import (
	booking "Api_Gateway/genproto/booking"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SearchServices godoc
// @Summary Searchs Services
// @Description Searchs service and returns details
// @Tags search
// @Security ApiKeyAuth
// @Param by_price query bool false "Sort by price"
// @Param location query string false "Location"
// @Param price query string false "Price"
// @Param duration query string false "Duration"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} booking.ListServicesResponses
// @Failure 500 {object} string "Server error while processing SearchServices"
// @Router /customer/search/service [get]
func (h *Handler) SearchServices(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "10")
	byPrice := c.DefaultQuery("by_price", "false")
	loc := c.Query("location")
	price := c.Query("price")
	duration := c.Query("duration")

	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		priceFloat = 0
	}
	durInt, err := strconv.Atoi(duration)
	if err != nil {
		durInt = 0
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	req := &booking.Filter{
		ByPrice:  byPrice == "true",
		Location: loc,
		Price:    float32(priceFloat),
		Duration: int32(durInt),
		Page:     int32(pageInt),
		Limit:    int32(limitInt),
	}
	fmt.Println("req: ", req)

	res, err := h.Search.SearchServices(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// SearchProviders godoc
// @Summary Searchs Providers
// @Description Searchs Providers and returns details
// @Tags search
// @Security ApiKeyAuth
// @Param by_rating query bool false "Sort by rating"
// @Param by_noc query bool false "Sort by number of comments"
// @Param company query string false "Company Name"
// @Param location query string false "Location"
// @Param availability query string false "Ishlash Vaqt oralig`i (8:00-20:00)"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} booking.ListProvidersResponses
// @Failure 500 {object} string "Server error while processing SearchProviders"
// @Router /customer/search/provider [get]
func (h *Handler) SearchProviders(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "10")
	loc := c.Query("location")
	byRating := c.DefaultQuery("by_rating", "false")
	byNoc := c.DefaultQuery("by_noc", "false")
	availability := c.Query("availability")
	company := c.Query("company")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	req := &booking.Filter{
		Location:         loc,
		ByRating:         byRating == "true",
		NumberOfComments: byNoc == "true",
		CompanyName:      company,
		ScheduledTime:    availability,
		Page:             int32(pageInt),
		Limit:            int32(limitInt),
	}

	res, err := h.Search.SearchProviders(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
