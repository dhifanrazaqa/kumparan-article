package router

import (
	"net/http"

	"github.com/dhifanrazaqa/kumparan-article/internal/handlers"
	"github.com/dhifanrazaqa/kumparan-article/pkg/middleware"
	"github.com/gorilla/mux"
)

func RegisterArticleRoutes(r *mux.Router, h *handlers.ArticleHandler, jwtSecret string) {
	articleRouter := r.PathPrefix("/articles").Subrouter()

	articleRouter.HandleFunc("", h.GetArticles).Methods(http.MethodGet)
	articleRouter.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", h.GetArticleByID).Methods(http.MethodGet)

	authed := articleRouter.PathPrefix("").Subrouter()
	authed.Use(func(next http.Handler) http.Handler {
		return middleware.JWT(next, jwtSecret)
	})
	authed.HandleFunc("", h.CreateArticle).Methods(http.MethodPost)
	authed.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", h.UpdateArticle).Methods(http.MethodPut)
	authed.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", h.DeleteArticle).Methods(http.MethodDelete)
}
