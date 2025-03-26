package error

import "errors"

var (
	ErrInternalServer          = errors.New("internal server error")
	ErrRegisterOrLoginRequired = errors.New("email and password are required")
)

var (
	ErrRegisterInvalidEmail     = errors.New("invalid email")
	ErrRegisterInvalidPassword  = errors.New("password must be at least 8 characters")
	ErrRegisterDuplicatedEmail  = errors.New("email already exists")
	ErrRegisterEmailRequired    = errors.New("email is required")
	ErrRegisterPasswordRequired = errors.New("password is required")
	ErrRegisterNameRequired     = errors.New("name is required")
)

var (
	ErrLoginEmailNotFound   = errors.New("email doesn't exist")
	ErrLoginInvalidPassword = errors.New("wrong password")
)

var (
	ErrVerificationTokenInvalid    = errors.New("token invalid")
	ErrResetTokenStillValid        = errors.New("link reset password still active")
	ErrVerificationTokenStillValid = errors.New("link verification still active")
)
