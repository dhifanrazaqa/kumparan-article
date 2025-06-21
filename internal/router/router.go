package router

import (
	"github.com/dhifanrazaqa/kumparan-article/internal/handlers"
	"github.com/gorilla/mux"
)

type Deps struct {
	UserHandler *handlers.UserHandler
}

func SetupRouter(d Deps) *mux.Router {
	router := mux.NewRouter()

	RegisterUserRoutes(router, d.UserHandler)

	return router
}
