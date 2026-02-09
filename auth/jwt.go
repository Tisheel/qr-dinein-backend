package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	StaffID      int    `json:"staffId"`
	RestaurantID int    `json:"restaurantId"`
	Role         string `json:"role"`
	Username     string `json:"username"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret      []byte
	expiryHours int
}

func NewJWTManager() (*JWTManager, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	expiryHours := 24
	if e := os.Getenv("JWT_EXPIRY_HOURS"); e != "" {
		if parsed, err := strconv.Atoi(e); err == nil && parsed > 0 {
			expiryHours = parsed
		}
	}

	return &JWTManager{
		secret:      []byte(secret),
		expiryHours: expiryHours,
	}, nil
}

func (m *JWTManager) GenerateToken(staffID, restaurantID int, role, username string) (string, int64, error) {
	expiresAt := time.Now().Add(time.Duration(m.expiryHours) * time.Hour)

	claims := &Claims{
		StaffID:      staffID,
		RestaurantID: restaurantID,
		Role:         role,
		Username:     username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt.Unix(), nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
