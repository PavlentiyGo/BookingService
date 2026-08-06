package core_http_request

import (
	core_errors "avitoBooking/internal/core/errors"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var Validator = validator.New()

func DecodeAndValidate(r *http.Request, dest any) error {

	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		return fmt.Errorf("failed to decode request: %w: %w", err, core_errors.ErrInvalidRequest)
	}

	if err := Validator.Struct(dest); err != nil {
		return fmt.Errorf("failed to validate request: %w: %w", err, core_errors.ErrInvalidRequest)
	}

	return nil
}
