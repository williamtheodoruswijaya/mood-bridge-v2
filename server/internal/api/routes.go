package api

import (
	"database/sql"
	"mood-bridge-v2/server/internal/handler"
	"mood-bridge-v2/server/internal/middleware"
	"mood-bridge-v2/server/internal/repository"
	"mood-bridge-v2/server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// step 1: define sebuah Handler sebagai struct yang akan digunakan untuk mengumpulkan semua handler yang ada dalam rest api kita.
type Handlers struct {
	UserHandler handler.UserHandler
}

// step 2: buat method untuk setiap route yang ada dalam api kita. misal kita mau bikin route untuk create user, kita bisa bikin method CreateUser
func SetupRoutes(db *sql.DB) *gin.Engine {
	return initRoutes(initHandler(db))
}

func initHandler(db *sql.DB) Handlers {
	// Inisialisasi validator juga
	validator := validator.New()

	// Inisialisasi repository, handler, dan services disini
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(db, userRepository)
	userHandler := handler.NewUserHandler(userService, *validator)

	return Handlers{
		UserHandler: userHandler,
	}
}

func initRoutes(h Handlers) *gin.Engine {
	// Inisialisasi router
	router := gin.Default()

	// Terapkan middleware untuk menangani panic
	router.Use(middleware.HandlePanic())

	// Lakukan grouping
	api := router.Group("/api")

	// Buat routes untuk user
	user := api.Group("/user")
	{
		user.POST("/register", h.UserHandler.Create)
		user.POST("/login", h.UserHandler.Login)
		
		user.Use(middleware.Authenticate()) // Terapkan middleware untuk semua route di bawah ini
		user.GET("/by-username/:username", h.UserHandler.Find)
		user.GET("/by-email", h.UserHandler.FindByEmail)
		user.GET("/by-id/:id", h.UserHandler.FindByID)
		user.GET("/all", h.UserHandler.FindAll)
	}


	return router
}
