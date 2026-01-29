package server

import (
	"fmt"
	"omiai-server/internal/api"
	"omiai-server/internal/conf"

	"github.com/iWuxc/go-wit/pprof"
	"github.com/iWuxc/go-wit/transport"
	"github.com/iWuxc/go-wit/transport/http"
)

// NewHTTPServer . new HTTP server.
func NewHTTPServer(router api.RouterInterface) []transport.ServerInterface {
	// http server
	httpServer := http.NewServer(
		http.Alias("http"),
		http.Address(fmt.Sprintf(":%d", conf.GetConfig().Server.HTTPPort)),
		//http.Timeout(120*time.Second),
		http.Handler(router.Register()),
	)

	// monitor server
	monitorServer := http.NewServer(
		http.Alias("monitor"),
		http.Address(fmt.Sprintf(":%d", conf.GetConfig().Server.MonitorPort)),
		http.Handler(pprof.NewHandler()),
	)

	return []transport.ServerInterface{httpServer, monitorServer}
}
