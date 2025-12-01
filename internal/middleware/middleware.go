package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/pkg/jwt"
	"github.com/codepnw/simple-bank/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	token *jwt.JWTToken
}

func NewMiddleware(token *jwt.JWTToken) *AuthMiddleware {
	return &AuthMiddleware{token: token}
}

func (m *AuthMiddleware) Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "header is missing")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid token format")
			c.Abort()
			return
		}

		claims, err := m.token.VerifyAccessToken(parts[1])
		if err != nil {
			log.Printf("verify token failed: %v", err)
			response.Unauthorized(c, "verify token failed")
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, consts.ContextUserClaimsKey, claims)
		ctx = context.WithValue(ctx, consts.ContextUserIDKey, claims.UserID)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
