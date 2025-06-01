package handler

import (
	"context"
	"mood-bridge-v2/server/internal/model/request"
	"mood-bridge-v2/server/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FriendHandler interface {
	AddFriend(c *gin.Context)
	AcceptRequest(c *gin.Context)
	GetFriends(c *gin.Context)
	Delete(c *gin.Context)
	GetFriendRequests(c *gin.Context)
	GetFriendRecommendation(c *gin.Context)
}

type FriendHandlerImpl struct {
	FriendService service.FriendService
	validate validator.Validate
}

func NewFriendHandler(friendService service.FriendService, validate validator.Validate) *FriendHandlerImpl {
	return &FriendHandlerImpl{
		FriendService: friendService,
		validate: validate,
	}
}

func (h *FriendHandlerImpl) AddFriend(c *gin.Context) {
	var req request.FriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response, err := h.FriendService.AddFriend(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"message": "Friend request sent successfully",
			"data": response,
		})
	}
}

func (h *FriendHandlerImpl) AcceptRequest(c *gin.Context) {
	var req request.FriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response, err := h.FriendService.AcceptRequest(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"message": "Friend request accepted successfully",
			"data": response,
		})
	}
}

func (h *FriendHandlerImpl) GetFriends(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "User ID is required",
		})
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "Invalid User ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	friends, err := h.FriendService.GetFriends(ctx, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"message": "Friends retrieved successfully",
			"data": friends,
		})
	}
}

func (h *FriendHandlerImpl) Delete(c *gin.Context) {
	friendID := c.Param("id")
	if friendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "Friend ID is required",
		})
		return
	}

	friendIDInt, err := strconv.Atoi(friendID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "Invalid Friend ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	message, err := h.FriendService.Delete(ctx, friendIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	if message == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":   http.StatusNotFound,
			"message": "Friend not found",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"message": message,
		})
	}
}

func (h *FriendHandlerImpl) GetFriendRequests(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "User ID is required",
		})
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "Invalid User ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	friendRequests, err := h.FriendService.GetFriendRequests(ctx, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	if len(friendRequests) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code":   http.StatusNotFound,
			"message": "No friend requests found",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"message": "Friend requests retrieved successfully",
			"data": friendRequests,
		})
	}
}

func (h *FriendHandlerImpl) GetFriendRecommendation(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "User ID is required",
		})
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"message": "Invalid User ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	recommendations, err := h.FriendService.GetFriendRecommendation(ctx, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	if len(recommendations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code":   http.StatusNotFound,
			"message": "No friend recommendations found",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"message": "Friend recommendations retrieved successfully",
			"data": recommendations,
		})
	}
}