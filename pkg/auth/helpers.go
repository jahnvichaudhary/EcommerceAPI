package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetUserId(ctx context.Context, abort bool) string {
	userId, err := GetUserIdInt(ctx, abort)
	if err != nil {
		return ""
	}
	return strconv.Itoa(userId)
}

func GetUserIdInt(ctx context.Context, abort bool) (int, error) {
	accountId, ok := ctx.Value("userID").(uint64)
	log.Println("userID", accountId)
	if !ok {
		if abort {
			ginContext, _ := ctx.Value("GinContextKey").(*gin.Context)
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
		return 0, errors.New("UserId not found in context")
	}
	return int(accountId), nil
}
