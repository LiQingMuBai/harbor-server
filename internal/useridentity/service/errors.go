package service

import "errors"

var (
	ErrInvalidUID         = errors.New("invalid uid")
	ErrInvalidID          = errors.New("invalid id")
	ErrInvalidType        = errors.New("invalid type")
	ErrUserNotFound       = errors.New("user not found")
	ErrAuthNotFound       = errors.New("auth not found")
	ErrAuthAlreadyExists  = errors.New("auth already exists")
	ErrRealNameTooLong    = errors.New("real name too long")
	ErrCardPhotosRequired = errors.New("card photos required")
)
