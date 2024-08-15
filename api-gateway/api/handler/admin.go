package handler

import (
	pb "Api_Gateway/genproto/users"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// GetUser godoc
// @Summary Gets user
// @Description Retrieves user info by ID
// @Tags admin
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} users.GetProfileResponse
// @Failure 400 {object} string "Invalid user id"
// @Failure 500 {object} string "Server error while processing request"
// @Router /admin/user/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	h.Log.Info("GetUser handler is invoked")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid user id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Admin.GetProfile(ctx, &pb.Id{UserId: id})
	if err != nil {
		er := errors.Wrap(err, "error getting user").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("GetUser handler is completed")
	c.JSON(http.StatusOK, resp)
}

// UpdateUser godoc
// @Summary Updates user
// @Description Updates user info by ID
// @Tags admin
// @Security ApiKeyAuth
// @Param user body users.UpdateProfileRequest true "New user data"
// @Success 200 {object} users.UserResponse
// @Failure 400 {object} string "Invalid user id"
// @Failure 500 {object} string "Server error while processing request"
// @Router /admin/user [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	h.Log.Info("UpdateUser handler is invoked")

	var req pb.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		er := errors.Wrap(err, "invalid data provided").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Admin.UpdateProfileA(ctx, &req)
	if err != nil {
		er := errors.Wrap(err, "error updating user").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("UpdateUser handler is completed")
	c.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
// @Summary Deletes user
// @Description Deletes user info by ID
// @Tags admin
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 400 {object} string "Invalid user id"
// @Failure 500 {object} string "Server error while processing request"
// @Router /admin/user/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	h.Log.Info("DeleteUser handler is invoked")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid user id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	_, err = h.Admin.DeleteUser(ctx, &pb.Id{UserId: id})
	if err != nil {
		er := errors.Wrap(err, "error deleting user").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("DeleteUser handler is completed")
	c.JSON(http.StatusOK, "User deleted successfully")
}

// FetchUsers godoc
// @Summary Fetches users
// @Description Retrieves users from the database by filtering
// @Tags admin
// @Security ApiKeyAuth
// @Param full_name query string false "Full name"
// @Param role query string false "Role"
// @Param page query int false "Page number"
// @Param limit query int false "Number of users per page"
// @Success 200 {object} users.UserResponses
// @Failure 400 {object} string "Invalid pagination parameters"
// @Failure 500 {object} string "Server error while processing request"
// @Router /admin/users [get]
func (h *Handler) FetchUsers(c *gin.Context) {
	h.Log.Info("FetchUsers handler is invoked")

	name := c.Query("full_name")
	role := c.Query("role")
	page := c.Query("page")
	limit := c.Query("limit")
	var p, l int

	if page != "" {
		pag, err := strconv.Atoi(page)
		if err != nil {
			er := errors.Wrap(err, "invalid pagination parameters").Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
			h.Log.Error(er)
			return
		}
		p = pag
	}

	if limit != "" {
		lim, err := strconv.Atoi(limit)
		if err != nil {
			er := errors.Wrap(err, "invalid pagination parameters").Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
			h.Log.Error(er)
			return
		}
		l = lim
	}

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	resp, err := h.Admin.FetchUsers(ctx, &pb.Filter{
		FirstName: name,
		Role:      role,
		Page:      int32(p),
		Limit:     int32(l),
	})
	if err != nil {
		er := errors.Wrap(err, "error fetching users").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("FetchUsers handler is completed")
	c.JSON(http.StatusOK, resp)
}
