package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func registerUser(t *testing.T, router http.Handler, user models.CreateUserRequest) models.UserResponse {
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var res struct {
		Data    models.UserResponse `json:"data"`
		Message string              `json:"message"`
	}

	err := json.Unmarshal(rr.Body.Bytes(), &res)
	require.NoError(t, err)
	return res.Data
}

func TestUserHandler_CRUD(t *testing.T) {
	clearDatabase(testDbPool)

	var createdUserID string
	var userToken string
	userCredentials := models.CreateUserRequest{
		Username: "user_crud_test",
		Name:     "User CRUD Test",
		Password: "password123",
	}

	t.Run("sukses membuat user baru", func(t *testing.T) {
		resUser := registerUser(t, testRouter, userCredentials)
		assert.Equal(t, userCredentials.Username, resUser.Username)
		createdUserID = resUser.ID
		require.NotEmpty(t, createdUserID)
	})

	t.Run("login untuk mendapatkan token", func(t *testing.T) {
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: userCredentials.Username,
			Password: userCredentials.Password,
		})
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		testRouter.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		var res struct {
			Data    models.AuthResponse `json:"data"`
			Message string              `json:"message"`
		}
		err := json.Unmarshal(rr.Body.Bytes(), &res)

		require.NoError(t, err)
		userToken = res.Data.AccessToken
		require.NotEmpty(t, userToken)
	})

	t.Run("sukses mendapatkan user berdasarkan ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", createdUserID), nil)
		rr := httptest.NewRecorder()
		testRouter.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var res struct {
			Data models.UserResponse `json:"data"`
		}
		err := json.Unmarshal(rr.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.Equal(t, userCredentials.Username, res.Data.Username)
	})

	t.Run("sukses mendapatkan semua user", func(t *testing.T) {
		registerUser(t, testRouter, models.CreateUserRequest{Username: "user_lain", Name: "User Lain", Password: "passwordlain"})

		req, _ := http.NewRequest("GET", "/users", nil)
		rr := httptest.NewRecorder()
		testRouter.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		var res struct {
			Data []models.UserResponse `json:"data"`
		}
		err := json.Unmarshal(rr.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(res.Data), 2)
	})

	t.Run("sukses memperbarui user sendiri", func(t *testing.T) {
		updateBody, _ := json.Marshal(models.UpdateUserRequest{Username: "user_crud_updated"})
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%s", createdUserID), bytes.NewBuffer(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)

		rr := httptest.NewRecorder()
		testRouter.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		var res struct {
			Data models.UserResponse `json:"data"`
		}
		err := json.Unmarshal(rr.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.Equal(t, "user_crud_updated", res.Data.Username)
	})

	t.Run("gagal memperbarui user lain", func(t *testing.T) {
		registerUser(t, testRouter, models.CreateUserRequest{Username: "user_lain_2", Name: "User Lain 2", Password: "passwordlain2"})

		updateBody, _ := json.Marshal(models.UpdateUserRequest{Username: "tidak_akan_berhasil"})
		updateReq, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%s", createdUserID), bytes.NewBuffer(updateBody))
		updateReq.Header.Set("Content-Type", "application/json")
	})

	t.Run("sukses menghapus user sendiri", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%s", createdUserID), nil)
		req.Header.Set("Authorization", "Bearer "+userToken)

		rr := httptest.NewRecorder()
		testRouter.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("memverifikasi user telah dihapus", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", createdUserID), nil)
		rr := httptest.NewRecorder()
		testRouter.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
