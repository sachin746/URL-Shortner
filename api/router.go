package api

import (
	"time"

	"URL-Shortner/log"

	"URL-Shortner/api/middleware"
	apiv1 "URL-Shortner/api/v1"

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

	}

	// Serve static files
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")
	router.StaticFile("/", "./frontend/index.html")
	router.StaticFile("/index.html", "./frontend/index.html")
	router.StaticFile("/login.html", "./frontend/login.html")
	router.StaticFile("/register.html", "./frontend/register.html")
	router.StaticFile("/dashboard.html", "./frontend/dashboard.html")
	router.StaticFile("/tech.html", "./frontend/tech.html")
	router.StaticFile("/favicon.ico", "./frontend/favicon.png")

	// add rate limiting middleware
	router.Use(middleware.RateLimitingMiddleware())
	router.POST("/custom/shortcode", apiv1.HandleCustomShortenURL)
	router.POST("/shortcode", apiv1.HandleShortenURL)
	router.GET("/:shortcode", apiv1.HandleGetURL)

	log.Sugar.Infof("Router initialized with version 1 endpoints")
	return router, nil
}
