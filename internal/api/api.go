package api

import (
	"net/http"
)

// RouterInterface routers .
type RouterInterface interface {
	Register() http.Handler
}
