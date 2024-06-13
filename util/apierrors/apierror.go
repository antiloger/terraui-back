package apierror

import "errors"

var (
	ErrAuthFail     = errors.New("authantication fail ex: password or username")
	ErrDynamoClient = errors.New("dynamodb client failed")
	ErrZeroData     = errors.New("data is empty")
)
