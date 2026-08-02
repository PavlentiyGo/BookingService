package core_routes

import (
	"avitoBooking/internal/core/middleware"
	"net/http"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []core_middleware.Middleware
}
