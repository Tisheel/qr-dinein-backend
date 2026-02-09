package service

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"

	"qr-dinein-backend/model"
)

const (
	otpLength            = 6
	otpExpiryMinutes     = 5
	sessionExpiryMinutes = 30
	maxOTPAttempts       = 5
	resendCooldownSecs   = 60
	maxOTPPerHour        = 5
)

// Redis key prefixes
const (
	keyPrefixOTP       = "otp:"
	keyPrefixCooldown  = "otp_cooldown:"
	keyPrefixRateLimit = "otp_rate:"
	keyPrefixSession   = "customer_session:"
)

// Customer handles customer OTP operations
type Customer struct {
	smsService *SMSService
}

// NewCustomer creates a new customer service
func NewCustomer(smsService *SMSService) *Customer {
	return &Customer{
		smsService: smsService,
	}
}

// SendOTP generates and sends OTP to the customer
func (svc *Customer) SendOTP(ctx *gofr.Context, req *model.SendOTPRequest) (*model.OTPResponse, error) {
	if req.PhoneNumber == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	if req.RestaurantID <= 0 {
		return nil, fmt.Errorf("restaurant ID is required")
	}

	// Check rate limit (max 5 OTP requests per phone per hour)
	rateLimitKey := keyPrefixRateLimit + req.PhoneNumber
	rateCount, _ := ctx.Redis.Get(ctx, rateLimitKey).Int()
	if rateCount >= maxOTPPerHour {
		return nil, fmt.Errorf("too many OTP requests. Please try again later")
	}

	// Check resend cooldown (60 seconds between requests)
	cooldownKey := keyPrefixCooldown + req.PhoneNumber
	exists, _ := ctx.Redis.Exists(ctx, cooldownKey).Result()
	if exists > 0 {
		ttl, _ := ctx.Redis.TTL(ctx, cooldownKey).Result()
		return nil, fmt.Errorf("please wait %d seconds before requesting another OTP", int(ttl.Seconds()))
	}

	// Generate 6-digit OTP
	otp, err := generateOTP(otpLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	ctx.Logger.Infof("OTP generated: %s", otp)

	// Store OTP data in Redis
	otpKey := keyPrefixOTP + req.PhoneNumber
	otpData := model.OTPData{
		OTP:      otp,
		Attempts: 0,
	}
	otpDataJSON, _ := json.Marshal(otpData)
	ctx.Redis.Set(ctx, otpKey, string(otpDataJSON), time.Duration(otpExpiryMinutes)*time.Minute)

	// Set resend cooldown
	ctx.Redis.Set(ctx, cooldownKey, "1", time.Duration(resendCooldownSecs)*time.Second)

	// Increment rate limit counter
	if rateCount == 0 {
		ctx.Redis.Set(ctx, rateLimitKey, "1", time.Hour)
	} else {
		ctx.Redis.Incr(ctx, rateLimitKey)
	}

	// Send OTP asynchronously
	go func() {
		if err := svc.smsService.SendOTP(req.PhoneNumber, otp); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to send OTP to %s: %v\n", req.PhoneNumber, err)
		}
	}()

	return &model.OTPResponse{
		Message:   "OTP sent successfully",
		ExpiresIn: otpExpiryMinutes * 60,
	}, nil
}

// VerifyOTP verifies the OTP and creates a session
func (svc *Customer) VerifyOTP(ctx *gofr.Context, req *model.VerifyOTPRequest) (*model.VerifyOTPResponse, error) {
	if req.PhoneNumber == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	if req.OTP == "" {
		return nil, fmt.Errorf("OTP is required")
	}

	if req.RestaurantID <= 0 {
		return nil, fmt.Errorf("restaurant ID is required")
	}

	// Get OTP data from Redis
	otpKey := keyPrefixOTP + req.PhoneNumber
	otpDataJSON, err := ctx.Redis.Get(ctx, otpKey).Result()
	if err != nil {
		return nil, fmt.Errorf("OTP expired or not found. Please request a new OTP")
	}

	var otpData model.OTPData
	if err := json.Unmarshal([]byte(otpDataJSON), &otpData); err != nil {
		return nil, fmt.Errorf("invalid OTP data")
	}

	// Check attempts
	if otpData.Attempts >= maxOTPAttempts {
		ctx.Redis.Del(ctx, otpKey)
		return nil, fmt.Errorf("too many failed attempts. Please request a new OTP")
	}

	// Verify OTP
	if otpData.OTP != req.OTP {
		// Increment attempts
		otpData.Attempts++
		updatedJSON, _ := json.Marshal(otpData)

		// Get remaining TTL and update with same expiry
		ttl, _ := ctx.Redis.TTL(ctx, otpKey).Result()
		ctx.Redis.Set(ctx, otpKey, string(updatedJSON), ttl)

		remainingAttempts := maxOTPAttempts - otpData.Attempts
		return nil, fmt.Errorf("invalid OTP. %d attempts remaining", remainingAttempts)
	}

	// OTP verified - delete it
	ctx.Redis.Del(ctx, otpKey)

	// Generate session token
	sessionToken := uuid.New().String()
	expiresAt := time.Now().Add(time.Duration(sessionExpiryMinutes) * time.Minute).Unix()

	// Store session in Redis
	sessionKey := keyPrefixSession + sessionToken
	session := model.CustomerSession{
		PhoneNumber:  req.PhoneNumber,
		RestaurantID: req.RestaurantID,
		Verified:     true,
	}
	sessionJSON, _ := json.Marshal(session)
	ctx.Redis.Set(ctx, sessionKey, string(sessionJSON), time.Duration(sessionExpiryMinutes)*time.Minute)

	return &model.VerifyOTPResponse{
		SessionToken: sessionToken,
		ExpiresAt:    expiresAt,
		PhoneNumber:  req.PhoneNumber,
	}, nil
}

// GetSession retrieves a customer session by token
func (svc *Customer) GetSession(ctx *gofr.Context, sessionToken string) (*model.CustomerSession, error) {
	if sessionToken == "" {
		return nil, fmt.Errorf("session token is required")
	}

	sessionKey := keyPrefixSession + sessionToken
	sessionJSON, err := ctx.Redis.Get(ctx, sessionKey).Result()
	if err != nil {
		return nil, fmt.Errorf("session expired or invalid")
	}

	var session model.CustomerSession
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, fmt.Errorf("invalid session data")
	}

	return &session, nil
}

// InvalidateSession removes a customer session (call after order is placed)
func (svc *Customer) InvalidateSession(ctx *gofr.Context, sessionToken string) error {
	if sessionToken == "" {
		return fmt.Errorf("session token is required")
	}

	sessionKey := keyPrefixSession + sessionToken
	ctx.Redis.Del(ctx, sessionKey)
	return nil
}

// generateOTP generates a cryptographically secure random OTP
func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}
