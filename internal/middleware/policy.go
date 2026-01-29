package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iWuxc/go-wit/log"
	"github.com/iWuxc/go-wit/utils"
)

// SetPolicy sets the policy for the current request.
func SetPolicy() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = utils.GetUUID()
		}

		blw := bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		if origin := c.Request.Header.Get("Origin"); origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization, x-api-key, x-code, x-channel,x-signature,x-timestamp,x-platform")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Content-Type", "application/json;charset=UTF-8")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Set("request_id", requestID)
		c.Set("api-key", c.Request.Header.Get("x-api-key"))
		c.Set("request_url", c.Request.URL.String())

		log.WithContext(c).Printf(c.Request.URL.String())

		if c.Request.Method == http.MethodPost && c.Request.ContentLength < 8192 {
			rawData, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				log.WithContext(c).Errorf(err.Error())
			}
			log.WithContext(c).Printf("requestBody: %s", rawData)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))
		}

		c.Header("x-requestID", requestID)
		c.Next()
		if blw.body.Len() <= 8192 {
			log.WithContext(c).Infof("responseBody: %s", blw.body.String())
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
