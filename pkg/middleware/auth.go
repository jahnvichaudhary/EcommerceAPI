package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"

	"github.com/rasadov/EcommerceAPI/pkg/auth"
)

func AuthorizeJWT(jwtService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Cookie("token")
		if err != nil || authCookie == "" {
			c.Set("userID", "")
			c.Next()
			return
		}

		token, err := jwtService.ValidateToken(authCookie)
		if err != nil {
			c.Set("userID", "")
			c.Next()
			return
		}

		if claims, ok := token.Claims.(*auth.JWTCustomClaims); ok && token.Valid {
			log.Println("Successfully validated token")
			log.Println("User ID from token:", claims.UserID)

			// Here we are setting the userID in the both go default context and gin context
			c.Set("userID", claims.UserID)
			ctxWithVal := context.WithValue(c.Request.Context(), "userID", claims.UserID)
			c.Request = c.Request.WithContext(ctxWithVal)
		} else {
			c.Set("userID", "")
		}

		c.Next()
	}
}
