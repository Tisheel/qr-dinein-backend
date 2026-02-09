package service

import (
	"fmt"
	"qr-dinein-backend/auth"
	"qr-dinein-backend/model"
	"qr-dinein-backend/store"

	"gofr.dev/pkg/gofr"
)

type Auth struct {
	staffStore        *store.Staff
	jwtManager        *auth.JWTManager
	superuserUsername  string
	superuserPassword string
}

func NewAuth(staffStore *store.Staff, jwtManager *auth.JWTManager, superuserUsername, superuserPassword string) *Auth {
	return &Auth{
		staffStore:        staffStore,
		jwtManager:        jwtManager,
		superuserUsername:  superuserUsername,
		superuserPassword: superuserPassword,
	}
}

func (s *Auth) Login(ctx *gofr.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if req.Pin == "" {
		return nil, fmt.Errorf("pin is required")
	}

	staff, err := s.staffStore.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if staff.Pin != req.Pin {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !staff.Active {
		return nil, fmt.Errorf("staff account is inactive")
	}

	token, expiresAt, err := s.jwtManager.GenerateToken(staff.ID, staff.RestaurantID, staff.Role, staff.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Staff:     staff,
	}, nil
}

func (s *Auth) SuperuserLogin(ctx *gofr.Context, req *model.SuperuserLoginRequest) (*model.SuperuserLoginResponse, error) {
	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	if s.superuserUsername == "" || s.superuserPassword == "" {
		return nil, fmt.Errorf("superuser login is not configured")
	}

	if req.Username != s.superuserUsername || req.Password != s.superuserPassword {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, expiresAt, err := s.jwtManager.GenerateToken(0, 0, "superuser", req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.SuperuserLoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Role:      "superuser",
		Username:  req.Username,
	}, nil
}
