package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/service"
	"github.com/Snuslyk/schedule-backend-go/internal/storage"
)

type AuthHandler struct {
	authService   *service.AuthService
	avatarService *storage.AvatarService
}

func NewAuthHandler(authService *service.AuthService, avatarService *storage.AvatarService) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		avatarService: avatarService,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register user and return JWT access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      domain.RegisterDTO    true  "Registration payload"
// @Success      201    {object}  domain.AuthResponse
// @Failure      400    {object}  domain.ErrorResponse
// @Failure      409    {object}  domain.ErrorResponse
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto domain.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if dto.Email == "" || dto.Password == "" || dto.Name == "" {
		BadRequest(w, "name, email, and password are required")
		return
	}

	resp, err := h.authService.Register(r.Context(), w, dto)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			Conflict(w, err.Error())
			return
		}
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, resp)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      domain.LoginDTO       true  "Login payload"
// @Success      200    {object}  domain.AuthResponse
// @Failure      400    {object}  domain.ErrorResponse
// @Failure      404    {object}  domain.ErrorResponse
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto domain.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if dto.Email == "" || dto.Password == "" {
		BadRequest(w, "email and password are required")
		return
	}

	resp, err := h.authService.Login(r.Context(), w, dto)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			NotFound(w, "User not found")
			return
		}
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusOK, resp)
}

// Refresh godoc
// @Summary      Refresh token
// @Description  Refresh JWT access token using cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  domain.AuthResponse
// @Failure      401  {object}  domain.ErrorResponse
// @Router       /api/auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	resp, err := h.authService.Refresh(r.Context(), w, r)
	if err != nil {
		Unauthorized(w, "Refresh token should be updated")
		return
	}

	JSON(w, http.StatusOK, resp)
}

// Logout godoc
// @Summary      User logout
// @Description  Clear refresh token cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  domain.MessageResponse
// @Router       /api/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.authService.Logout(w)
	JSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Return profile of the authenticated user
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  domain.UserProfile
// @Failure      401  {object}  domain.ErrorResponse
// @Router       /api/user/profile [get]
// @Router       /api/auth/profile [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := GetAuthUser(r)
	if user == nil {
		Unauthorized(w, "Unauthorized")
		return
	}

	profile, err := h.authService.GetProfile(r.Context(), user.ID)
	if err != nil {
		NotFound(w, "User not found")
		return
	}

	JSON(w, http.StatusOK, profile)
}

// AddRole godoc
// @Summary      Update user role (Admin only)
// @Description  Add or update role for specified user ID
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Param        id    query     string  true  "User UUID"
// @Param        role  query     string  true  "Role name"
// @Success      200   {object}  domain.MessageResponse
// @Failure      400   {object}  domain.ErrorResponse
// @Router       /api/user/role [post]
// @Router       /api/auth/role [post]
func (h *AuthHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	role := r.URL.Query().Get("role")

	if id == "" || role == "" {
		BadRequest(w, "id and role are required query parameters")
		return
	}

	if err := h.authService.UpdateRoles(r.Context(), id, role); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"message": "Role updated"})
}

// UploadAvatar godoc
// @Summary      Upload user avatar
// @Description  Upload avatar image with optional cropData to S3
// @Tags         user
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        avatar    formData  file    true   "Avatar image file"
// @Param        cropData  formData  string  false  "JSON cropData {x, y, width, height}"
// @Success      200       {object}  domain.MessageResponse
// @Failure      400       {object}  domain.ErrorResponse
// @Router       /api/user/avatar [post]
// @Router       /api/auth/avatar [post]
func (h *AuthHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	user := GetAuthUser(r)
	if user == nil {
		Unauthorized(w, "Unauthorized")
		return
	}

	if err := r.ParseMultipartForm(2 << 20); err != nil {
		BadRequest(w, "File exceeds max allowed size of 2MB")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		BadRequest(w, "avatar file is required")
		return
	}
	defer file.Close()

	if header.Size > 2097152 {
		BadRequest(w, "File size must not exceed 2MB")
		return
	}

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		BadRequest(w, "File must be an image")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		InternalError(w, "Failed to read file")
		return
	}

	cropDataStr := r.FormValue("cropData")
	var crop storage.CropData
	if cropDataStr != "" {
		if err := json.Unmarshal([]byte(cropDataStr), &crop); err != nil {
			BadRequest(w, "Invalid cropData json")
			return
		}
	}

	if err := h.avatarService.Upload(r.Context(), user.ID, fileBytes, crop); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"message": "Avatar uploaded successfully"})
}
