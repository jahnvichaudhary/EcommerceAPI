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
	accountId, ok := ctx.Value("userID").(string)
	log.Println("userID", accountId)
	if !ok {
		if abort {
			ginContext, _ := ctx.Value("GinContextKey").(*gin.Context)
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
		return ""
	}
	return accountId
}

func GetUserIdInt(ctx context.Context, abort bool) (int, error) {
	idString := GetUserId(ctx, abort)
	log.Println("userID", idString)
	if idString != "" {
		idInt, err := strconv.ParseInt(idString, 10, 64)
		if err == nil {
			return int(idInt), nil
		}
	}
	return 0, errors.New("some error happened")
}
