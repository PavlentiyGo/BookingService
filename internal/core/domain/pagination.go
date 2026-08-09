package domain

import (
	core_errors "avitoBooking/internal/core/errors"
	"fmt"
)

const DefaultPage = 1
const DefaultPageSize = 20

type Pagination struct {
	Page     int
	PageSize int
	Total    int
}

func NewPagination(page int, pageSize int) (Pagination, error) {
	var pagination Pagination
	if page != 0 {
		if page < 0 {
			return Pagination{}, fmt.Errorf("page must be positive number: %w", core_errors.ErrInvalidRequest)
		}
		pagination.Page = page
	} else {
		pagination.Page = DefaultPage
	}
	if pageSize != 0 {
		if pageSize > 100 || pageSize < 0 {
			return Pagination{}, fmt.Errorf("pageSize must be less than 100 and higher than 0: %w", core_errors.ErrInvalidRequest)
		}
		pagination.PageSize = pageSize
	} else {
		pagination.PageSize = DefaultPageSize
	}

	return pagination, nil
}
