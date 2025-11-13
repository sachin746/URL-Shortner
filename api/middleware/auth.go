package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"URL-Shortner/business"

	"URL-Shortner/constants"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	log.Default().Println("AuthMiddleware initialized")
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, constants.ErrUnauthorized.SetErr(errors.New("authorization header is missing")))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, constants.ErrUnauthorized.SetErr(errors.New("invalid authorization header format")))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := business.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, constants.ErrUnauthorized.SetErr(err))
			c.Abort()
			return
		}

		c.Set("username", (*claims)["username"])
		c.Set("email", (*claims)["email"])
		c.Next()
	}
}
