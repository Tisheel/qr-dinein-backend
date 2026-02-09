package handler

import (
	"fmt"

	"gofr.dev/pkg/gofr"

	"qr-dinein-backend/model"
	"qr-dinein-backend/service"
)

// Customer handles customer OTP endpoints
type Customer struct {
	service *service.Customer
}

// NewCustomer creates a new customer handler
func NewCustomer(svc *service.Customer) *Customer {
	return &Customer{
		service: svc,
	}
}

// SendOTP handles POST /api/v1/customer/send-otp
func (h *Customer) SendOTP(ctx *gofr.Context) (interface{}, error) {
	var req model.SendOTPRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	return h.service.SendOTP(ctx, &req)
}

// VerifyOTP handles POST /api/v1/customer/verify-otp
func (h *Customer) VerifyOTP(ctx *gofr.Context) (interface{}, error) {
	var req model.VerifyOTPRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	return h.service.VerifyOTP(ctx, &req)
}

// GetSession handles GET /api/v1/customer/session?token={sessionToken}
func (h *Customer) GetSession(ctx *gofr.Context) (interface{}, error) {
	sessionToken := ctx.Param("token")
	if sessionToken == "" {
		return nil, fmt.Errorf("session token is required")
	}

	return h.service.GetSession(ctx, sessionToken)
}
