package utils

import "errors"

// Controller errors
var (
	// ErrBadRequestBody is used when the request body is not valid format
	ErrBadRequestBody = errors.New("bad request body")

	// ErrInvalidCredentials is used when the user's credentials are invalid
	ErrInvalidCredentials = errors.New("invalid username or password")

	// ErrDidntHavePermission is used when the user doesn't have permission to access or modify the resource
	ErrDidntHavePermission = errors.New("you didn't have permission to do this action")

	// ErrInvalidNumber is used when the number covertion is invalid
	ErrInvalidNumber = errors.New("invalid number")
)

// Service errors
var (
	// ErrFieldNotMatch is used when the field in the request body is not match with the field in the template that saved in the database
	ErrFieldNotMatch = errors.New("document fields doesn't match with template fields")
)

// Repository errors
var (
	// ErrUsernameAlreadyExist is used when the username is already exist in the database
	ErrUsernameAlreadyExist = errors.New("user with provided username already exist")

	// ErrUserNotFound is used when the user is not found in the database
	ErrUserNotFound = errors.New("user not found")
)

//OTP
var (
	ErrOtpInvalid = errors.New("OTP invalid")
	ErrOtpExpired = errors.New("OTP Expired")
)

//Movie
var (
	ErrMovieNotFound     = errors.New("Movie Not Found")
	ErrTitleAlreadyExist = errors.New("Movie with porvided title already exist")
)
