package middleware

import (
	"omiai-server/internal/data"
	"omiai-server/pkg/auth"
	"omiai-server/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/redis"
)

func Authorization(db *data.DB, redis *redis.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.MiddlewareErrorResponse(c, response.ParamsCommonError, "未登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.MiddlewareErrorResponse(c, response.ParamsCommonError, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			response.MiddlewareErrorResponse(c, response.ParamsCommonError, "登录已失效")
			c.Abort()
			return
		}

		// 存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
