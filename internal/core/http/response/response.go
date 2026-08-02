package response

import (
	core_errors "avitoBooking/internal/core/errors"
	core_logger "avitoBooking/internal/core/logger"
	"context"
	"encoding/json"
	"errors"
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

func (r *Responser) ErrorResponse(
	err error,
) {

	if errors.Is(err, core_errors.ErrExpiredToken) || errors.Is(err, core_errors.ErrInvalidRequest) {
		r.logger.Debug("got error", zap.Error(err))
		r.writeErrorJson(http.StatusBadRequest, core_errors.ErrInvalidRequest.Error(), err.Error())
	} else if errors.Is(err, core_errors.ErrNotAuthorized) {
		r.logger.Debug("got error", zap.Error(err))
		r.writeErrorJson(http.StatusUnauthorized, core_errors.ErrInvalidRequest.Error(), err.Error())
	} else if errors.Is(err, core_errors.ErrForbidden) {
		r.logger.Debug("got error", zap.Error(err))
		r.writeErrorJson(http.StatusForbidden, core_errors.ErrInvalidRequest.Error(), err.Error())
	} else {
		r.logger.Error("got INTERNAL error", zap.Error(err))
		r.writeErrorJson(http.StatusInternalServerError, core_errors.ErrInternalError.Error(), err.Error())
	}

}
func (r *Responser) writeErrorJson(
	statusCode int,
	error string,
	message string,
) {
	response := map[string]string{
		"error":   error,
		"message": message,
	}
	bytes, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		r.logger.Error("failed to marshal data into json", zap.Error(err))
		return
	}
	http.Error(r.rw, string(bytes), statusCode)
}
