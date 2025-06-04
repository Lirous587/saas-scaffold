package auth

import (
	"github.com/pkg/errors"
	"sass-scaffold/internal/common/reskit/codes"
	"sass-scaffold/internal/common/reskit/response"
	"sass-scaffold/internal/user/adapters"
	"sass-scaffold/internal/user/domain"
	"sass-scaffold/internal/user/service"
	"strings"

	"github.com/gin-gonic/gin"
)

var tokenServer domain.TokenService

func init() {
	tokenCache := adapters.NewRedisTokenCache()
	userRepo := adapters.NewPSQLUserRepository()
	tokenServer = service.NewTokenService(tokenCache, userRepo)
}

const (
	authHeaderKey	= "Authorization"
	bearerPrefix	= "Bearer "
)

// 解析 Authorization 头部的 Token
func parseTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(authHeaderKey)
	if authHeader == "" {
		return "", errors.New("token为空")
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("token格式错误")
	}

	return strings.TrimPrefix(authHeader, bearerPrefix), nil
}

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头解析 Token
		tokenStr, err := parseTokenFromHeader(c)
		if err != nil {
			response.Error(c, err)
			return
		}

		// 2. 解析 Token
		payload, isExpire, err := tokenServer.ValidateAccessToken(tokenStr)
		if err != nil {
			if isExpire {
				response.Error(c, codes.ErrTokenExpired)
			} else {
				response.Error(c, codes.ErrTokenInvalid)
			}
			return
		}

		// 3. 将用户 相关信息存入上下文
		c.Set("user_id", payload.UserID)
		//c.Set("random_code", payload.RandomCode)

		c.Next()
	}
}
