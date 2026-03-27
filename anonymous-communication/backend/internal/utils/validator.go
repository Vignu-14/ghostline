package utils

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"anonymous-communication/backend/internal/models"

	"github.com/google/uuid"
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9._]{3,50}$`)

func NormalizeUsername(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func NormalizeEmail(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func ValidateRegisterRequest(input models.RegisterRequest) error {
	fields := make(map[string]string)

	username := NormalizeUsername(input.Username)
	if username == "" {
		fields["username"] = "username is required"
	} else if !usernamePattern.MatchString(username) {
		fields["username"] = "username must be 3-50 characters and use letters, numbers, dots, or underscores"
	}

	email := NormalizeEmail(input.Email)
	if email == "" {
		fields["email"] = "email is required"
	} else if _, err := mail.ParseAddress(email); err != nil {
		fields["email"] = "email must be valid"
	}

	if passwordMessage := validatePassword(input.Password); passwordMessage != "" {
		fields["password"] = passwordMessage
	}

	return models.NewValidationError(fields)
}

func ValidateLoginRequest(input models.LoginRequest) error {
	fields := make(map[string]string)

	if NormalizeUsername(input.Username) == "" {
		fields["username"] = "username is required"
	}

	if strings.TrimSpace(input.Password) == "" {
		fields["password"] = "password is required"
	}

	return models.NewValidationError(fields)
}

func ValidateImpersonateRequest(input models.ImpersonateRequest) error {
	fields := make(map[string]string)

	if _, err := uuid.Parse(strings.TrimSpace(input.TargetUserID)); err != nil {
		fields["target_user_id"] = "target_user_id must be a valid UUID"
	}

	if strings.TrimSpace(input.ImpersonationPassword) == "" {
		fields["impersonation_password"] = "impersonation_password is required"
	}

	return models.NewValidationError(fields)
}

func validatePassword(password string) string {
	if len(password) < 8 {
		return "password must be at least 8 characters long"
	}

	var hasUpper bool
	var hasLower bool
	var hasDigit bool
	var hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return "password must include uppercase, lowercase, number, and special character"
	}

	return ""
}
