package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
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
	ErrCategoryRequired = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogCategoryRequired", mrerr.ErrorKindUser, "category ID is required")

	// ErrCategoryNotAvailable - category with ID is not available.
	ErrCategoryNotAvailable = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogCategoryNotAvailable", mrerr.ErrorKindUser, "category with ID={{ .id }} is not available")

	// ErrCategoryNotFound - category with ID not found.
	ErrCategoryNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogCategoryNotFound", mrerr.ErrorKindUser, "category with ID={{ .id }} not found")
)

// CategoryErrors - comment func.
func CategoryErrors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrCategoryRequired,
		ErrCategoryNotAvailable,
		ErrCategoryNotFound,
	}
}
