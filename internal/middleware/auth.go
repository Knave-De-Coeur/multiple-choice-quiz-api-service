package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-api-service/internal/api"
)

type IAuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
}

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		authSplit := strings.Split(auth, "Bearer ")
		if len(authSplit) < 2 {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				api.GenerateMessageResponse("no token", nil, errors.New("missing token in request")),
			)
			return
		}
		tokenString := authSplit[1]
		if tokenString == "" {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				api.GenerateMessageResponse("no token", nil, errors.New("missing token in request")),
			)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
			}
			return []byte(a.jwtSecret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				api.GenerateMessageResponse("something went wrong with the token", nil, err),
			)
			return
		}

		var tokenUID int
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			tokenUID, _ = strconv.Atoi(fmt.Sprint(claims["sub"]))
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatusJSON(
					http.StatusForbidden,
					api.GenerateMessageResponse("expired token", nil, errors.New("token is no longer valid")),
				)
				return
			}
		} else {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				api.GenerateMessageResponse("bad token", nil, errors.New("token is not valid")),
			)
			return
		}

		paramUserID := c.Param("uID")
		userIDint, _ := strconv.Atoi(paramUserID)

		if tokenUID != userIDint {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				api.GenerateMessageResponse("incorrect token for user", nil, errors.New("token is not valid")),
			)
			return
		}

		c.Set("user_id", tokenUID)

		c.Next()
	}
}
