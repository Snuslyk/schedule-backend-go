package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/config"
	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
	"github.com/Snuslyk/schedule-backend-go/internal/storage"
	"github.com/Snuslyk/schedule-backend-go/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidToken      = errors.New("invalid token")
)

type AuthService struct {
	userRepo      *postgres.UserRepository
	avatarService *storage.AvatarService
	cfg           *config.Config
}

func NewAuthService(
	userRepo *postgres.UserRepository,
	avatarService *storage.AvatarService,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		avatarService: avatarService,
		cfg:           cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, w http.ResponseWriter, dto domain.RegisterDTO) (*domain.AuthResponse, error) {
	existing, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	hashed, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	roles := dto.Roles
	if len(roles) == 0 {
		roles = []string{string(domain.RoleStudent)}
	}

	u := &domain.User{
		ID:       uuid.New().String(),
		Name:     dto.Name,
		Email:    dto.Email,
		Password: hashed,
		Roles:    roles,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	return s.auth(w, u.ID, u.Roles)
}

func (s *AuthService) Login(ctx context.Context, w http.ResponseWriter, dto domain.LoginDTO) (*domain.AuthResponse, error) {
	u, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrInvalidCredentials
	}

	if !utils.CheckPassword(dto.Password, u.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.auth(w, u.ID, u.Roles)
}

func (s *AuthService) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) (*domain.AuthResponse, error) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil || cookie.Value == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}

	return s.auth(w, u.ID, u.Roles)
}

func (s *AuthService) Logout(w http.ResponseWriter) {
	s.setCookie(w, "", time.Unix(0, 0))
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*domain.UserProfile, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}

	var avatarURL *string
	if s.avatarService != nil {
		avatarURL, _ = s.avatarService.GetAvatarURL(ctx, u.ID)
	}

	return &domain.UserProfile{
		Name:   u.Name,
		Email:  u.Email,
		Roles:  u.Roles,
		Avatar: avatarURL,
	}, nil
}

func (s *AuthService) UpdateRoles(ctx context.Context, id, role string) error {
	return s.userRepo.UpdateRole(ctx, id, role)
}

func (s *AuthService) auth(w http.ResponseWriter, id string, roles []string) (*domain.AuthResponse, error) {
	accessToken, err := s.generateToken(id, roles, s.cfg.JWTAccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(id, roles, s.cfg.JWTRefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	expires := time.Now().Add(s.cfg.JWTRefreshTokenTTL)
	s.setCookie(w, refreshToken, expires)

	return &domain.AuthResponse{AccessToken: accessToken}, nil
}

func (s *AuthService) generateToken(id string, roles []string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":    id,
		"roles": roles,
		"exp":   time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) setCookie(w http.ResponseWriter, value string, expires time.Time) {
	sameSite := http.SameSiteLaxMode
	secure := !s.cfg.IsDev()
	if !s.cfg.IsDev() {
		sameSite = http.SameSiteNoneMode
	}

	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    value,
		Domain:   s.cfg.CookieDomain,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	}
	http.SetCookie(w, cookie)
}
