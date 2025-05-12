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

type PostHandler interface {
	Create (c *gin.Context)
	Find (c *gin.Context)
}

type PostHandlerImpl struct {
	PostService service.PostService
	validate validator.Validate
}

func NewPostHandler(postService service.PostService, validate validator.Validate) PostHandler {
	return &PostHandlerImpl{
		PostService: postService,
		validate:    validate,
	}
}

func (h *PostHandlerImpl) Create(c *gin.Context) {
	var request request.CreatePostRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request format, please check the data you sent",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response, err := h.PostService.Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Post created successfully",
			"data":    response,
		})
		return
	}
}

func (h *PostHandlerImpl) Find(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Post ID is required",
		})
		return
	}
	postID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid Post ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response, err := h.PostService.Find(ctx, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Success",
			"data":    response,
		})
		return
	}
}

