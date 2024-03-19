package catalog

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrerr"
)

type (
	CategoryAPI interface {
		// CheckingAvailability - error: FactoryErrCategoryRequired | FactoryErrCategoryNotFound | Failed
		CheckingAvailability(ctx context.Context, itemID uuid.UUID) error
	}
)

var (
	FactoryErrCategoryRequired = mrerr.NewFactory(
		"errCatalogCategoryRequired", mrerr.ErrorKindUser, "category ID is required")

	FactoryErrCategoryNotFound = mrerr.NewFactory(
		"errCatalogCategoryNotFound", mrerr.ErrorKindUser, "category with ID={{ .id }} not found")
)
