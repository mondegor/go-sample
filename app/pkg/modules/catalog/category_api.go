package catalog

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrerr"
)

type (
	CategoryAPI interface {
		// CheckingAvailability - error:
		//   - FactoryErrCategoryRequired
		//   - FactoryErrCategoryNotAvailable
		//   - FactoryErrCategoryNotFound
		//   - Failed
		CheckingAvailability(ctx context.Context, itemID uuid.UUID) error
	}
)

var (
	FactoryErrCategoryRequired = mrerr.NewFactory(
		"errCatalogCategoryRequired", mrerr.ErrorKindUser, "category ID is required")

	FactoryErrCategoryNotAvailable = mrerr.NewFactory(
		"errCatalogCategoryNotAvailable", mrerr.ErrorKindUser, "category with ID={{ .id }} is not available")

	FactoryErrCategoryNotFound = mrerr.NewFactory(
		"errCatalogCategoryNotFound", mrerr.ErrorKindUser, "category with ID={{ .id }} not found")
)
