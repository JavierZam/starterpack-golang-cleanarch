package auth

import (
	"encoding/json"
	"net/http"

	"starterpack-golang-cleanarch/internal/utils"
	globalErrors "starterpack-golang-cleanarch/internal/utils/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	service   *AuthService
	validator *validator.Validate
}

func NewAuthHandler(s *AuthService, v *validator.Validate) *AuthHandler {
	return &AuthHandler{service: s, validator: v}
}

func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", h.Register).Methods("POST")
	router.HandleFunc("/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/auth/refresh", h.Refresh).Methods("POST")
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleHTTPError(w, globalErrors.NewBadRequest("Invalid request payload", nil), r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleHTTPError(w, globalErrors.NewBadRequest(err.Error(), nil), r)
		return
	}

	userResp, err := h.service.RegisterUser(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, userResp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleHTTPError(w, globalErrors.NewBadRequest("Invalid request payload", nil), r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleHTTPError(w, globalErrors.NewBadRequest(err.Error(), nil), r)
		return
	}

	authResp, err := h.service.LoginUser(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, authResp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleHTTPError(w, globalErrors.NewBadRequest("Invalid request payload", nil), r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleHTTPError(w, globalErrors.NewBadRequest(err.Error(), nil), r)
		return
	}

	authResp, err := h.service.RefreshTokens(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, authResp)
}
