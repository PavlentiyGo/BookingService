package handlers

import (
	"avitoBooking/internal/core/auth"
	core_domain "avitoBooking/internal/core/domain"
	core_middleware "avitoBooking/internal/core/middleware"
	"avitoBooking/internal/core/routes"
	handlers_service "avitoBooking/internal/handlers/service"
)

type Handlers struct {
	jwtProvider auth.JwtProvider
	service     handlers_service.BookingService
}

func NewHandlers(
	provider auth.JwtProvider,
	service handlers_service.BookingService,
) *Handlers {
	return &Handlers{
		jwtProvider: provider,
		service:     service,
	}
}
func (h *Handlers) RoomsRoutes() []core_routes.Route {
	return []core_routes.Route{
		{
			Method:  "GET",
			Path:    "/rooms/list",
			Handler: h.GetRooms,
			Middlewares: []core_middleware.Middleware{
				core_middleware.Auth(h.jwtProvider, core_domain.UserRole, core_domain.AdminRole),
			},
		}, {
			Method:  "POST",
			Path:    "/rooms/create",
			Handler: h.CreateRoom,
			Middlewares: []core_middleware.Middleware{
				core_middleware.Auth(h.jwtProvider, core_domain.AdminRole),
			},
		},
	}
}
func (h *Handlers) AuthRoutes() []core_routes.Route {
	return []core_routes.Route{
		{
			Method:  "POST",
			Path:    "/register",
			Handler: h.Register,
		}, {
			Method:  "POST",
			Path:    "/login",
			Handler: h.Login,
		},
		{
			Method:  "POST",
			Path:    "/dummyLogin",
			Handler: h.DummyLogin,
		},
	}
}
func (h *Handlers) ScheduleRoutes() []core_routes.Route {
	return []core_routes.Route{
		{
			Method:  "POST",
			Path:    "/rooms/{roomId}/schedule/create",
			Handler: h.CreateSchedule,
			Middlewares: []core_middleware.Middleware{
				core_middleware.Auth(h.jwtProvider, core_domain.AdminRole),
			},
		},
	}
}

func (h *Handlers) SlotsRoutes() []core_routes.Route {
	return []core_routes.Route{
		{
			Method:  "GET",
			Path:    "/rooms/{roomId}/slots/list",
			Handler: h.GetSlots,
			Middlewares: []core_middleware.Middleware{
				core_middleware.Auth(h.jwtProvider, core_domain.AdminRole, core_domain.UserRole),
			},
		},
	}
}
func (h *Handlers) BookingRoutes() []core_routes.Route {
	return []core_routes.Route{
		{
			Method:  "POST",
			Path:    "/bookings/create",
			Handler: h.CreateBooking,
			Middlewares: []core_middleware.Middleware{
				core_middleware.Auth(h.jwtProvider, core_domain.UserRole),
			},
		}, {
			Method:  "GET",
			Path:    "/bookings/list",
			Handler: h.GetBookings,
			Middlewares: []core_middleware.Middleware{
				core_middleware.Auth(h.jwtProvider, core_domain.AdminRole),
			},
		},
		{
			Method:  "GET",
			Path:    "/bookings/my",
			Handler: h.GetMyBooking,
			Middlewares: []core_middleware.Middleware{
				core_middleware.Auth(h.jwtProvider, core_domain.UserRole),
			},
		}, {
			Method:  "POST",
			Path:    "/bookings/{bookingId}/cancel",
			Handler: h.CancelMyBooking,
			Middlewares: []core_middleware.Middleware{
				core_middleware.Auth(h.jwtProvider, core_domain.UserRole),
			},
		},
	}
}
