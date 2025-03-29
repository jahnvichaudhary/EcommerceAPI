package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rasadov/EcommerceMicroservices/account"
)

// AuthorizeJWT is a Gin middleware that checks for a valid JWT in the "Authorization" header.
// Usage (for example):
//
//	router.GET("/protected", AuthorizeJWT(jwtService), protectedHandler)
func AuthorizeJWT(jwtService account.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header not provided",
			})
			return
		}

		// Typically: "Authorization: Bearer <token>"
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		tokenString := splitToken[1]
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Token is valid, extract our custom claims
		if claims, ok := token.Claims.(*account.JWTCustomClaims); ok && token.Valid {
			// Put user ID into context for subsequent handlers
			c.Set("userID", claims.UserID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		c.Next()
	}
}
