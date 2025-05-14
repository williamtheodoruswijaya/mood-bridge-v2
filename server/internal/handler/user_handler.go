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

type UserHandler interface {
	Create(c *gin.Context)
	Find(c *gin.Context)
	FindByEmail(c *gin.Context)
	FindByID(c *gin.Context)
	FindAll(c *gin.Context)
	Login(c *gin.Context)
	Update(c *gin.Context)
}

type UserHandlerImpl struct {
	UserService service.UserService
	validate validator.Validate // nyobain pake validator buat validasi request
}

func NewUserHandler(userService service.UserService, validate validator.Validate) UserHandler {
	return &UserHandlerImpl{
		UserService: userService,
		validate:    validate,
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

func (h *UserHandlerImpl) FindByID(c *gin.Context) {
	// step 1: ambil username dari path
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "ID is required",
		})
		return
	}
	// step 1.1: convert id ke int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid ID format",
		})
		return
	}

	// step 2: buat context buat ngatur time-out (handle connection time-out)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// step 3: call service-nya buat find user-nya
	response, err := h.UserService.FindByID(ctx, idInt)
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

func (h *UserHandlerImpl) FindAll(c *gin.Context) {
	// step 1: buat context buat ngatur time-out (handle connection time-out)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// step 2: call service-nya buat find all user-nya
	response, err := h.UserService.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to find all users",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Users found successfully",
			"data":    response,
		})
		return
	}
}

func (h *UserHandlerImpl) Login(c *gin.Context) {
	// step 1: ambil request dari body sekaligus validasi
	var request request.ValidateUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}
	// Validasi request (apakah sama dengan struct atau engga)
	err = h.validate.Struct(request)
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

	// step 3: call service-nya buat login user-nya
	token, err := h.UserService.Login(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to login user",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "User logged in successfully",
			"data":    token,
		})
		return
	}
}

func (h *UserHandlerImpl) Update(c *gin.Context) {
	// step 1: ambil id dari path
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "ID is required",
		})
		return
	}
	// step 1.1: convert id ke int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid ID format",
		})
		return
	}

	// step 2: ambil request dari body
	var request request.UpdateUserRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request",
		})
		return
	}

	// step 3: buat context buat ngatur time-out (handle connection time-out)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// step 4: call service-nya buat update user-nya
	response, err := h.UserService.Update(ctx, idInt, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update user",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "User updated successfully",
			"data":    response,
		})
		return
	}
}