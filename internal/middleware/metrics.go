package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/metrics/stat"
)

func metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		w := c.Writer

		c.Next()

		uri := c.FullPath()
		if len(uri) == 0 {
			if w.Status() == http.StatusNotFound {
				uri = "404"
			} else {
				uri = c.Request.RequestURI
			}
		}
		stat.APPRequestTotalCount.With("uri", uri, "code", strconv.Itoa(w.Status())).Inc()

		latency := float64(time.Since(start)) / float64(time.Second)
		stat.APPRequestHistogram.With("uri", uri).Observe(latency)
	}
}
