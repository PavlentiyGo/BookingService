package core_errors

import "errors"

var ErrInternalError = errors.New("INTERNAL_ERROR")
var ErrInvalidRequest = errors.New("INVALID_REQUEST")

var ErrExpiredToken = errors.New("token is expired")
var ErrForbidden = errors.New("you are not allowed to do this action")
var ErrNotAuthorized = errors.New("you are not authorized")
var ErrWrongAuthType = errors.New("token has wrong authorization type")

var ErrInvalidRoomCapacity = errors.New("capacity mustn't be negative")

var ErrWrongRoomId = errors.New("invalid room id")
var ErrInvalidDays = errors.New("wrong days of week argument for room schedule")
var ErrInvalidTime = errors.New("invalid room time for start or end: start must be less than end and in format 23:59")
var ErrMissingDate = errors.New("missing date argument in path")
var ErrInvalidDateTime = errors.New("invalid date time in query")

var ErrConferenceServiceUnavailable = errors.New("conference service is unavailable now")

var ErrInvalidSlotId = errors.New("invalid slot id in request")
var ErrSlotInThePast = errors.New("you can't book slot in the past")

var ErrInvalidBookingId = errors.New("invalid booking id")
var ErrEmailExists = errors.New("such email is already in use")
