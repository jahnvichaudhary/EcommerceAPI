package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rasadov/EcommerceAPI/account/internal/user"
)

func AuthorizeJWT(jwtService user.AuthService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authCookie, err := context.Cookie("token")
		if err != nil || authCookie == "" {
			context.Set("userID", "")
			context.Next()
			return
		}

		token, err := jwtService.ValidateToken(authCookie)
		if err != nil {
			// Token is invalid => treat as anonymous or invalid user
			context.Set("userID", "")
			context.Next()
			return
		}

		// Token is valid => set user info
		if claims, ok := token.Claims.(*user.JWTCustomClaims); ok && token.Valid {
			context.Set("userID", claims.UserID)
		} else {
			context.Set("userID", "")
		}

		// Continue the request
		context.Next()
	}
}
