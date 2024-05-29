package logic

import "errors"

var (
	ErrInvalidEmail            = errors.New("invalid email")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrUserNotFound            = errors.New("user not found")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
)
