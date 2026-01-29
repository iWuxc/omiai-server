package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	defaultLimiter = 100
	defaultRate    = 10 * time.Millisecond
)

// Throttle is a middleware that limits the number of requests per second.
func Throttle() gin.HandlerFunc {
	limit := rate.Every(defaultRate)
	limiter := rate.NewLimiter(limit, defaultLimiter)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
