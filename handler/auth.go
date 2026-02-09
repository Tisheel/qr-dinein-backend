package handler

import (
	"qr-dinein-backend/auth"
	"qr-dinein-backend/model"
	"qr-dinein-backend/service"

	"gofr.dev/pkg/gofr"
)

type Auth struct {
	service *service.Auth
}

func NewAuth(svc *service.Auth) *Auth {
	return &Auth{service: svc}
}

func (h *Auth) Login(ctx *gofr.Context) (interface{}, error) {
	var req model.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	return h.service.Login(ctx, &req)
}

func (h *Auth) SuperuserLogin(ctx *gofr.Context) (interface{}, error) {
	var req model.SuperuserLoginRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	return h.service.SuperuserLogin(ctx, &req)
}

func (h *Auth) Me(ctx *gofr.Context) (interface{}, error) {
	claims := auth.GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, nil
	}

	return map[string]interface{}{
		"staffId":      claims.StaffID,
		"restaurantId": claims.RestaurantID,
		"role":         claims.Role,
		"username":     claims.Username,
	}, nil
}
