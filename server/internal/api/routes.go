package api

import (
	"database/sql"
	"mood-bridge-v2/server/internal/handler"
	"mood-bridge-v2/server/internal/middleware"
	"mood-bridge-v2/server/internal/repository"
	"mood-bridge-v2/server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

// step 1: define sebuah Handler sebagai struct yang akan digunakan untuk mengumpulkan semua handler yang ada dalam rest api kita.
type Handlers struct {
	UserHandler    handler.UserHandler
	PostHandler    handler.PostHandler
	CommentHandler handler.CommentHandler
	FriendHandler handler.FriendHandler
	ChatHandler    handler.ChatHandler
}

// step 2: buat method untuk setiap route yang ada dalam api kita. misal kita mau bikin route untuk create user, kita bisa bikin method CreateUser
func SetupRoutes(db *sql.DB, redisClient *redis.Client) *gin.Engine {
	return initRoutes(initHandler(db, redisClient))
}

func initHandler(db *sql.DB, redisClient *redis.Client) Handlers {
	// Inisialisasi validator juga
	validator := validator.New()

	// Inisialisasi repository, handler, dan services disini
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(db, userRepository)
	userHandler := handler.NewUserHandler(userService, *validator)

	postRepository := repository.NewPostRepository()
	postService := service.NewPostService(db, postRepository, userRepository, service.NewMoodPredictionService(), redisClient)
	postHandler := handler.NewPostHandler(postService, *validator)

	commentRepository := repository.NewCommentRepository()
	commentService := service.NewCommentService(commentRepository, userRepository, postRepository, db, redisClient)
	commentHandler := handler.NewCommentHandler(commentService, *validator)

	friendRepository := repository.NewFriendRepository()
	friendService := service.NewFriendService(friendRepository, userRepository, db, redisClient)
	friendHandler := handler.NewFriendHandler(friendService, *validator)

	chatRepository := repository.NewChatRepository(db)
	websocketHub := service.NewConcreteHub(chatRepository)
	chatService := service.NewChatService(chatRepository, websocketHub)
	chatHandler := handler.NewChatHandler(chatService)

	return Handlers{
		UserHandler:    userHandler,
		PostHandler:    postHandler,
		CommentHandler: commentHandler,
		FriendHandler: friendHandler,
		ChatHandler:    chatHandler,
	}
}

func initRoutes(h Handlers) *gin.Engine {
	// Inisialisasi router
	router := gin.Default()

	// Terapkan middleware untuk menangani panic
	router.Use(middleware.HandlePanic())
	// Terapkan middleware untuk CORS
	router.Use(middleware.CORSMiddleware())

	// Lakukan grouping
	api := router.Group("/api")

	// Buat routes untuk user
	user := api.Group("/user")
	{
		user.POST("/register", h.UserHandler.Create)
		user.POST("/login", h.UserHandler.Login)
		user.GET("/by-username/:username", h.UserHandler.Find)
		user.GET("/by-id/:id", h.UserHandler.FindByID)

		user.Use(middleware.Authenticate())
		user.GET("/by-email", h.UserHandler.FindByEmail)
		user.GET("/all", h.UserHandler.FindAll)
		user.PUT("/update/:id", h.UserHandler.Update)
	}

	post := api.Group("/post")
	{
		post.GET("/all", h.PostHandler.FindAll)
		post.GET("/by-id/:id", h.PostHandler.Find)
		post.GET("/by-userid/:id", h.PostHandler.FindByUserID)

		post.Use(middleware.Authenticate())
		post.POST("/create", h.PostHandler.Create)
		post.PUT("/update/:id", h.PostHandler.Update)
		post.DELETE("/delete/:id", h.PostHandler.Delete)
		post.GET("/friend-posts/:id", h.PostHandler.GetFriendPosts)
	}

	comment := api.Group("/comment")
	{
		comment.GET("/by-postid/:id", h.CommentHandler.GetAllByPostID)
		comment.GET("/by-id/:id", h.CommentHandler.GetByID)

		comment.Use(middleware.Authenticate())
		comment.POST("/create", h.CommentHandler.Create)
		comment.DELETE("/delete/:id", h.CommentHandler.Delete)
	}

	friend := api.Group("/friend")
	{
		friend.GET("/all/:id", h.FriendHandler.GetFriends)
		friend.GET("/requests/:id", h.FriendHandler.GetFriendRequests)
		friend.Use(middleware.Authenticate())
		friend.POST("/add", h.FriendHandler.AddFriend)
		friend.POST("/accept", h.FriendHandler.AcceptRequest)
		friend.DELETE("/delete/:id", h.FriendHandler.Delete)
		friend.GET("/recommendation/:id", h.FriendHandler.GetFriendRecommendation)
	}

	chat := api.Group("/chat")
	{
		chat.Use(middleware.Authenticate()) // gaboleh dituker
		chat.GET("/ws", h.ChatHandler.HandleWebSocketConnection)
		chat.GET("/history", h.ChatHandler.HandleFetchChatHistory)
		chat.POST("/messages/:message_id/read", h.ChatHandler.HandleMarkMessageAsRead)
	}

	return router
}
