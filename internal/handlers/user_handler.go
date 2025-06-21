package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/dhifanrazaqa/kumparan-article/internal/services"
	"github.com/dhifanrazaqa/kumparan-article/pkg/utils"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Request body is not valid")
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Username and password cannot be empty")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetUsers(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := h.userService.GetUserByID(r.Context(), vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), id, req, id)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.userService.DeleteUser(r.Context(), id, id)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
