package router

import (
	"net/http"

	"github.com/pi-prakhar/go-gcp-auth/internal/handlers"
	"github.com/pi-prakhar/go-gcp-auth/internal/middleware"
)

type Router struct {
	Mux        *http.ServeMux
	handler    *handlers.AuthHandler
	middleware *middleware.AuthMiddleware
}

func NewRouter(handler *handlers.AuthHandler) *Router {
	mux := http.NewServeMux()
	mw := middleware.AuthMiddleware{}

	router := &Router{
		Mux:        mux,
		handler:    handler,
		middleware: &mw,
	}

	router.initRoutes()
	return router
}

func (r *Router) initRoutes() {
	r.Mux.HandleFunc("/api/v1/auth/home", r.handler.HandleHome)
	r.Mux.HandleFunc("/api/v1/auth/google/login", r.handler.HandleGoogleLogin)
	r.Mux.HandleFunc("/api/v1/auth/google/callback", r.handler.HandleGoogleCallback)
	r.Mux.HandleFunc("/api/v1/auth/protected", r.middleware.IsAuthenticated(r.handler.HandleProtected))
}
