package errs

import "errors"

var (
	UserNotFound                    = errors.New("User not found")
	EmptyLoginOrPassword            = errors.New("Login or password is empty")
	LoginAlreadyExists              = errors.New("Login already exists")
	WrongLoginOrPassword            = errors.New("Wrong login or password")
	TokenExpired                    = errors.New("JWT expired")
	TokenNotFound                   = errors.New("Token not found in headers")
	TokenInvalid                    = errors.New("Invalid token")
	EmptyContextUser                = errors.New("User info not found in context")
	OrderNotFound                   = errors.New("Order not uploaded yet")
	OrderAlreadyUploadedByOtherUser = errors.New("Order already uploaded by other user")
	OrderAlreadyUploadedByThisUser  = errors.New("Order already uploaded by this user")
)
