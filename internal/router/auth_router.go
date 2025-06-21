package router

import (
	"net/http"

	"github.com/dhifanrazaqa/kumparan-article/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router, h *handlers.AuthHandler) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	authRouter.HandleFunc("/refresh", h.RefreshToken).Methods(http.MethodPost)
}
