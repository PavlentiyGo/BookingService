package response

import (
	core_errors "avitoBooking/internal/core/errors"
	core_logger "avitoBooking/internal/core/logger"
	repository_errors "avitoBooking/internal/repository/erorrs"
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Responser struct {
	rw     http.ResponseWriter
	logger *core_logger.Logger
}

func NewResponser(
	rw http.ResponseWriter,
	ctx context.Context,
) *Responser {

	logger := core_logger.FromContext(ctx)

	return &Responser{
		rw:     rw,
		logger: logger,
	}
}

type errorValue struct {
	statusCode int
	error      string
	logLevel   string
}

var errorMapper = map[error]errorValue{
	core_errors.ErrExpiredToken:             {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	core_errors.ErrWrongAuthType:            {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	core_errors.ErrInvalidRequest:           {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	repository_errors.ErrRoomAlreadyExists:  {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	core_errors.ErrInvalidRoomCapacity:      {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	core_errors.ErrWrongRoomId:              {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	core_errors.ErrInvalidDays:              {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	core_errors.ErrInvalidTime:              {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	core_errors.ErrMissingDate:              {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	repository_errors.ErrRoomScheduleExists: {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
	core_errors.ErrInvalidDateTime:          {statusCode: http.StatusBadRequest, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},

	repository_errors.ErrRoomNotFound: {statusCode: http.StatusNotFound, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},

	core_errors.ErrNotAuthorized: {statusCode: http.StatusUnauthorized, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},

	core_errors.ErrForbidden: {statusCode: http.StatusForbidden, error: core_errors.ErrInvalidRequest.Error(), logLevel: "DEBUG"},
}

func (r *Responser) ErrorResponse(
	err error,
) {

	val, ok := errorMapper[err]
	if !ok {
		r.logger.Error("got INTERNAL error", zap.Error(err))
		r.writeErrorJson(http.StatusInternalServerError, core_errors.ErrInternalError.Error(), err.Error())
		return
	}
	switch val.logLevel {
	case "DEBUG":
		r.logger.Debug("got error", zap.Error(err))
	}
	r.writeErrorJson(val.statusCode, val.error, err.Error())
}
func (r *Responser) writeErrorJson(
	statusCode int,
	error string,
	message string,
) {

	response := map[string]struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		"error": {
			Code:    error,
			Message: message,
		},
	}
	r.WriteJson(statusCode, response)
}
func (r *Responser) WriteJson(
	statusCode int,
	body any,
) {
	bytes, err := json.MarshalIndent(body, "", "    ")
	if err != nil {
		r.logger.Error("failed to marshal data into json", zap.Error(err))
		return
	}
	r.rw.WriteHeader(statusCode)
	if _, err = r.rw.Write(bytes); err != nil {
		r.logger.Error("failed to write data", zap.Error(err))
		r.rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
