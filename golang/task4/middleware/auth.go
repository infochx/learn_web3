package middleware

import (
	"blog/api"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		// 获取请求头中的JWT Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "bearer token is required"})
			ctx.Abort()
			return
		}

		// 解析 JWT Token
		claims := &api.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return api.JwtKey, nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid toke"})
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文中
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}
