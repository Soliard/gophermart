package errs

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrEmptyLoginOrPassword = errors.New("login or password is empty")
	ErrUserAlreadyExists    = errors.New("login already exists")
	ErrWrongLoginOrPassword = errors.New("wrong login or password")

	ErrTokenExpired  = errors.New("jwt expired")
	ErrTokenNotFound = errors.New("token not found in headers")
	ErrTokenInvalid  = errors.New("invalid token")

	ErrEmptyContextUser = errors.New("user info not found in context")

	ErrOrderNotFound                   = errors.New("order not uploaded yet")
	ErrOrderAlreadyUploadedByOtherUser = errors.New("order already uploaded by other user")
	ErrOrderAlreadyUploadedByThisUser  = errors.New("order already uploaded by this user")
	ErrOrderIsNotValid                 = errors.New("order is not valid")

	ErrBalanceInsufficient = errors.New("not enough points on balance")

	ErrWithdrawalAlreadyProcessed = errors.New("this withdraw already was processed")
	ErrWithdrawalsNotFound        = errors.New("withdrawals not found")

	ErrUnexpectedStatusAccrualService = errors.New("unexpected status code from accrual service")
)
