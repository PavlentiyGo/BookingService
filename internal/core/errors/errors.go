package core_errors

import "errors"

var ErrInternalError = errors.New("INTERNAL_ERROR")
var ErrInvalidRequest = errors.New("INVALID_REQUEST")

var ErrExpiredToken = errors.New("token is expired")
var ErrForbidden = errors.New("you are not allowed to do this action")
var ErrNotAuthorized = errors.New("you are not authorized")
