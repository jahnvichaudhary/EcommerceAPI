package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetSecretKey() string
}

type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJwtService(secretKey, issuer string) AuthService {
	return &jwtService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (service *jwtService) GenerateToken(userID string) (string, error) {
	claims := &JWTCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    service.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(service.secretKey))
}

func (service *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		encodedToken,
		&JWTCustomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(service.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		if claims.Issuer != service.issuer {
			return nil, errors.New("invalid issuer in token")
		}
		return token, nil
	}

	return nil, errors.New("invalid token claims")
}

func (service *jwtService) GetSecretKey() string {
	return service.secretKey
}

func GetUserId(ctx context.Context, abort bool) string {
	accountId, ok := ctx.Value("accountId").(string)
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
	if idString != "" {
		idInt, err := strconv.ParseInt(idString, 10, 64)
		if err == nil {
			return int(idInt), nil
		}
	}
	return 0, errors.New("Some Error happened")
}
