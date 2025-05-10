package handler

import (
	"context"
	"mood-bridge-v2/server/internal/model/request"
	"mood-bridge-v2/server/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Create(c *gin.Context)
	Find(c *gin.Context)
	FindByEmail(c *gin.Context)
}

type UserHandlerImpl struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &UserHandlerImpl{
		UserService: userService,
	}
}

func (h *UserHandlerImpl) Create(c *gin.Context) {
	// step 1: ambil request dari body
	var request request.CreateUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	// step 2: buat context buat ngatur time-out (handle connection time-out)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// step 3: call service-nya buat create task-nya
	response, err := h.UserService.Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "User created successfully",
			"data":    response,
		})
		return
	}
}

func (h *UserHandlerImpl) Find(c *gin.Context) {
	// step 1: ambil username dari path
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Username is required",
		})
		return
	}

	// step 2: buat context buat ngatur time-out (handle connection time-out)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// step 3: call service-nya buat find user-nya
	response, err := h.UserService.Find(ctx, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to find user",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "User found successfully",
			"data":    response,
		})
		return
	}
}

func (h *UserHandlerImpl) FindByEmail(c *gin.Context) {
	// step 1: ambil email dari params
	email := c.Query("email")
	
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Email is required",
		})
		return
	}

	// step 2: buat context buat ngatur time-out (handle connection time-out)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// step 3: call service-nya buat find user-nya
	response, err := h.UserService.FindByEmail(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to find user with this email",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "User found successfully",
			"data":    response,
		})
		return
	}
}
