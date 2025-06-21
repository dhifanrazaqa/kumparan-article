package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/dhifanrazaqa/kumparan-article/internal/services"
	"github.com/dhifanrazaqa/kumparan-article/pkg/middleware"
	"github.com/dhifanrazaqa/kumparan-article/pkg/utils"
	"github.com/gorilla/mux"
)

type ArticleHandler struct {
	articleService services.ArticleService
}

func NewArticleHandler(s services.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: s}
}

func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*models.Claims)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get user data from token")
		return
	}

	var req models.CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	article, err := h.articleService.CreateArticle(r.Context(), req, claims.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, "Article created successfully", article)
}

func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	page, _ := strconv.Atoi(queryParams.Get("page"))
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	params := models.ListArticlesParams{
		Query:  queryParams.Get("query"),
		Author: queryParams.Get("author"),
		Limit:  limit,
		Offset: offset,
	}

	paginatedResult, err := h.articleService.GetArticles(r.Context(), params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.WriteJSON(w, http.StatusOK, "Articles retrieved successfully", paginatedResult)
}

func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	article, err := h.articleService.GetArticleByID(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Article retrived successfully", article)
}

func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*models.Claims)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get user data from token")
		return
	}

	var req models.UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	article, err := h.articleService.UpdateArticle(r.Context(), id, req, claims.UserID)
	if err != nil {
		if err.Error() == "article not found" {
			utils.WriteError(w, http.StatusNotFound, err.Error())
		} else {
			utils.WriteError(w, http.StatusForbidden, err.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Article updated successfully", article)
}

func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*models.Claims)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get user data from token")
		return
	}

	err := h.articleService.DeleteArticle(r.Context(), id, claims.UserID)
	if err != nil {
		if err.Error() == "article not found" {
			utils.WriteError(w, http.StatusNotFound, err.Error())
		} else {
			utils.WriteError(w, http.StatusForbidden, err.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Article deleted successfully", nil)
}
