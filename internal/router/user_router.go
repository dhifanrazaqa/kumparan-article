package router

import (
	"net/http"

	"github.com/dhifanrazaqa/kumparan-article/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, h *handlers.UserHandler) {
	userRouter := router.PathPrefix("/users").Subrouter()

	userRouter.HandleFunc("", h.CreateUser).Methods(http.MethodPost)
	userRouter.HandleFunc("", h.GetUsers).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", h.GetUserByID).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", h.UpdateUser).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", h.DeleteUser).Methods(http.MethodDelete)
}