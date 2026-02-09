package model

// SendOTPRequest represents the request to send OTP
type SendOTPRequest struct {
	PhoneNumber  string `json:"phoneNumber"`
	RestaurantID int    `json:"restaurantId"`
}

// VerifyOTPRequest represents the request to verify OTP
type VerifyOTPRequest struct {
	PhoneNumber  string `json:"phoneNumber"`
	OTP          string `json:"otp"`
	RestaurantID int    `json:"restaurantId"`
}

// OTPResponse represents the response after sending OTP
type OTPResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expiresIn"` // seconds
}

// VerifyOTPResponse represents the response after successful verification
type VerifyOTPResponse struct {
	SessionToken string `json:"sessionToken"`
	ExpiresAt    int64  `json:"expiresAt"`
	PhoneNumber  string `json:"phoneNumber"`
}

// CustomerSession represents the verified customer session stored in Redis
type CustomerSession struct {
	PhoneNumber  string `json:"phoneNumber"`
	RestaurantID int    `json:"restaurantId"`
	Verified     bool   `json:"verified"`
}

// OTPData represents the OTP data stored in Redis
type OTPData struct {
	OTP      string `json:"otp"`
	Attempts int    `json:"attempts"`
}
