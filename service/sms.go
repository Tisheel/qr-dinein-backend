package service

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSService handles sending SMS via Twilio
type SMSService struct {
	client      *twilio.RestClient
	fromNumber  string
	enabled     bool
}

// NewSMSService creates a new SMS service
func NewSMSService() *SMSService {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")

	// If Twilio credentials are not set, disable SMS
	if accountSID == "" || authToken == "" || fromNumber == "" {
		return &SMSService{
			enabled: false,
		}
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &SMSService{
		client:     client,
		fromNumber: fromNumber,
		enabled:    true,
	}
}

// SendOTP sends an OTP to the given phone number
func (s *SMSService) SendOTP(phoneNumber, otp string) error {
	if !s.enabled {
		// Log OTP for development when Twilio is not configured
		fmt.Printf("[DEV] OTP for %s: %s\n", phoneNumber, otp)
		return nil
	}

	message := fmt.Sprintf("Your verification code is: %s. Valid for 5 minutes.", otp)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(s.fromNumber)
	params.SetBody(message)

	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}

// IsEnabled returns whether SMS service is enabled
func (s *SMSService) IsEnabled() bool {
	return s.enabled
}
