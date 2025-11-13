package api

import (
	"habit-tracker/log"
	"time"

	"habit-tracker/api/middleware"
	apiv1 "habit-tracker/api/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetRouter() (*gin.Engine, error) {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Or use your frontend domain e.g. http://localhost:5173
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/v1")
	{
		v1.POST("/register", apiv1.HandleRegisterUser)
		v1.POST("/login", apiv1.HandleLoginUser)
		v1.GET("/auth/google", apiv1.HandleGoogleAuth)
		v1.GET("/auth/google/callback", apiv1.HandleGoogleAuthCallback)
		v1.GET("/auth/github", apiv1.HandleGithubAuth)
		v1.GET("/auth/github/callback", apiv1.HandleGithubAuthCallback)

		// Protected route
		profile := v1.Group("/profile")
		profile.Use(middleware.AuthMiddleware())
		profile.GET("", apiv1.HandleGetUser)
		// profile.PUT("/profile", apiv1.HandleUpdateUser)
		// profile.DELETE("/profile", apiv1.HandleDeleteUser)
		// profile.POST("/logout", apiv1.HandleLogoutUser)
		// Posts
		posts := v1.Group("/posts", middleware.AuthMiddleware())
		{
			posts.POST("", apiv1.CreatePost)
			posts.GET("", apiv1.ListPosts)
			posts.GET("/:id", apiv1.GetPost)
			posts.POST("/comments", apiv1.CreateComment)
			posts.POST("/vote", apiv1.HandleVote)
		}

		// Habits
		habits := v1.Group("/habits", middleware.AuthMiddleware())
		{
			habits.POST("", apiv1.CreateHabit)
			habits.POST("/:habitId/members", apiv1.AddHabitMember)
			habits.GET("/:habitId/members", apiv1.GetHabitMembers)
		}
	}

	log.Sugar.Infof("Router initialized with version 1 endpoints")
	return router, nil
}
