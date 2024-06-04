package errors

import (
	"errors"
	"fmt"
	"net/http"
	"webdp/internal/api/http/response"
)

var (
	ErrBadFormatting  = errors.New("bad format")
	ErrBadInput       = errors.New("bad input")
	ErrBadRequest     = errors.New("bad request")
	ErrBadType        = errors.New("type error")
	ErrDatabase       = errors.New("database error")
	ErrForbidden      = errors.New("you are not allowed to perform this action")
	ErrInvalidToken   = errors.New("invalid authorization token")
	ErrMissingEnv     = errors.New("missing environment variable")
	ErrNotFound       = errors.New("could not find resource")
	ErrNotImplemented = errors.New("not implemented")
	ErrTimeout        = errors.New("timeout error")
	ErrUnauthorized   = errors.New("you are not authorized to perform this action")
	ErrUnexpected     = errors.New("something bad happened")
)

func ExpandError(err error) response.HttpResponse[any] {
	var status int
	var desc string
	unwrapped := err
	if errors.Unwrap(err) != nil {
		unwrapped = errors.Unwrap(err)
	}
	switch unwrapped {
	case ErrBadInput:
		fallthrough
	case ErrBadFormatting:
		fallthrough
	case ErrBadType:
		fallthrough
	case ErrBadRequest: // 400
		status, desc = http.StatusBadRequest, "Bad Request"
	case ErrInvalidToken:
		fallthrough
	case ErrUnauthorized: // 401
		status, desc = http.StatusUnauthorized, "Unauthorized"
	case ErrForbidden: // 403
		status, desc = http.StatusForbidden, "Forbidden"
	case ErrNotFound: // 404
		status, desc = http.StatusNotFound, "Not Found"
	case ErrDatabase:
		fallthrough
	case ErrMissingEnv:
		fallthrough
	case ErrTimeout:
		fallthrough
	case ErrUnexpected: // 500
		status, desc = http.StatusInternalServerError, "Unexpected error"
	case ErrNotImplemented: // 501
		status, desc = http.StatusNotImplemented, "Not Implemented"
	default:
		fmt.Println("--- ERROR EXPAND: could not match error", unwrapped.Error(), "!!! ---")
		fmt.Println("--- please handle this error.")
		status, desc = http.StatusInternalServerError, "Internal Server Error"
	}
	// TODO: Print propagated error message for testing only. Printing error
	// messages will leak information about internal mechanisms. Should be
	// masked in PROD to print a generic error message.
	return response.NewFail[any](status).WithTitle(desc).WithType(desc).WithDetail(err.Error())
}

/*
Converts "sql: no rows in result set" errors to NotFound errors. Converts the rest to Database errors.
Note: Can be extended as necessary.
*/
func WrapDBError(err error, method string, info string) error {
	if err == ErrNotFound || err.Error() == "sql: no rows in result set" {
		return fmt.Errorf("%w: %s not found", ErrNotFound, info)
	} else {
		return fmt.Errorf("%w: failed to %s %s: %s", ErrDatabase, method, info, err.Error())
	}
}
