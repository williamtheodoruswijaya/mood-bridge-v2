package handler

import (
	"mood-bridge-v2/server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CommentHandler interface {
	Create(c *gin.Context)
	GetAllByPostID(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
}

type CommentHandlerImpl struct {
	CommentService service.CommentService
	validate validator.Validate
}

func NewCommentHandler(commentService service.CommentService, validate validator.Validate) CommentHandler {
	return &CommentHandlerImpl{
		CommentService: commentService,
		validate:       validate,
	}
}

func (h *CommentHandlerImpl) Create(c *gin.Context) {
}

func (h *CommentHandlerImpl) GetAllByPostID(c *gin.Context) {
}

func (h *CommentHandlerImpl) Delete(c *gin.Context) {
}

func (h *CommentHandlerImpl) GetByID(c *gin.Context) {
}