package middleware

import (
	"fmt"
	"omiai-server/internal/data"
	"omiai-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-version"
	"github.com/iWuxc/go-wit/log"
	"github.com/iWuxc/go-wit/redis"
)

type CodeError struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	TokenKey string `json:"token_key"`
}

func Authorization(db *data.DB, redis *redis.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {

		v := c.Request.Header.Get("version")
		// 创建版本对象
		v1, err := version.NewVersion("1.1.0")
		if err != nil {
			fmt.Printf("Error parsing version: %v\n", err)
			return
		}

		v2, err := version.NewVersion(v)
		if err != nil {
			fmt.Printf("Error parsing version: %v\n", err)
			return
		}
		var loginError *CodeError
		// 版本比较
		if v1.GreaterThan(v2) {
			log.WithContext(c).Infof("checkLogin")
		}
		if loginError != nil {
			response.MiddlewareErrorResponse(c, response.Code(loginError.Code), loginError.Msg)
			return
		}
		c.Next()
	}

}
