package core_http_server

import (
	core_config "avitoBooking/internal/core/config"
	"avitoBooking/internal/core/middleware"
	core_routes "avitoBooking/internal/core/routes"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ChainRoutes(routes ...core_routes.Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {

		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		handler := core_middleware.ChainMiddlewares(route.Handler, route.Middlewares...)

		mux.Handle(pattern, handler)

	}
	return mux
}

func NewServer(
	config core_config.Config,
	routes []core_routes.Route,
	middlewares ...core_middleware.Middleware,
) *http.Server {

	mux := ChainRoutes(routes...)

	mux.Handle("/metrics", promhttp.Handler())

	handler := core_middleware.ChainMiddlewares(mux, middlewares...)

	server := &http.Server{
		Addr:    ":" + config.Addr,
		Handler: handler,
	}

	return server
}
