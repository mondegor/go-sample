package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrerr"
)

const (
	CategoryAvailabilityName = "Catalog.API.CategoryAvailability" // CategoryAvailabilityName - название API
)

type (
	// CategoryAvailability - comment interface.
	CategoryAvailability interface {
		// CheckingAvailability - error:
		//   - ErrCategoryRequired
		//   - ErrCategoryNotAvailable
		//   - ErrCategoryNotFound
		//   - Failed
		CheckingAvailability(ctx context.Context, itemID uuid.UUID) error
	}
)

var (
	// ErrCategoryRequired - category ID is required.
	ErrCategoryRequired = mrerr.NewProto(
		"catalog.errCategoryRequired", mrerr.ErrorKindUser, "category ID is required")

	// ErrCategoryNotAvailable - category with ID is not available.
	ErrCategoryNotAvailable = mrerr.NewProto(
		"catalog.errCategoryNotAvailable", mrerr.ErrorKindUser, "category with ID={{ .id }} is not available")

	// ErrCategoryNotFound - category with ID not found.
	ErrCategoryNotFound = mrerr.NewProto(
		"catalog.errCategoryNotFound", mrerr.ErrorKindUser, "category with ID={{ .id }} not found")
)
