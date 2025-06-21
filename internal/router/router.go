package router

import (
	"github.com/dhifanrazaqa/kumparan-article/internal/handlers"
	"github.com/gorilla/mux"
)

type Deps struct {
	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
}

func SetupRouter(d Deps) *mux.Router {
	router := mux.NewRouter()

	RegisterAuthRoutes(router, d.AuthHandler)
	RegisterUserRoutes(router, d.UserHandler)

	return router
}
