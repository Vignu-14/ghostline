package models

import "errors"

var (
	ErrUnauthorized                 = errors.New("unauthorized")
	ErrForbidden                    = errors.New("forbidden")
	ErrAdminOnly                    = errors.New("admin access required")
	ErrUserNotFound                 = errors.New("user not found")
	ErrInvalidCredentials           = errors.New("invalid credentials")
	ErrUsernameTaken                = errors.New("username is already taken")
	ErrEmailTaken                   = errors.New("email is already taken")
	ErrPostNotFound                 = errors.New("post not found")
	ErrMessageNotFound              = errors.New("message not found")
	ErrCannotLikeOwnPost            = errors.New("users cannot like their own posts")
	ErrCannotMessageSelf            = errors.New("users cannot message themselves")
	ErrDeleteForEveryoneNotAllowed  = errors.New("delete for everyone is only allowed on your own messages")
	ErrCannotImpersonateSelf        = errors.New("admins cannot impersonate themselves")
	ErrInvalidImpersonationPassword = errors.New("invalid impersonation password")
	ErrImpersonationNotConfigured   = errors.New("impersonation password is not configured")
	ErrImpersonationNotActive       = errors.New("impersonation is not active")
	ErrStorageNotConfigured         = errors.New("storage is not configured")
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func NewValidationError(fields map[string]string) error {
	if len(fields) == 0 {
		return nil
	}

	return &ValidationError{Fields: fields}
}
