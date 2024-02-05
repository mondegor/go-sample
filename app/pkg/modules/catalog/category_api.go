package catalog

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryAPI interface {
		// CheckingAvailability - error: FactoryErrCategoryNotFound or Failed
		CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error
	}
)

var (
	FactoryErrCategoryNotFound = mrerr.NewFactory(
		"errCatalogCategoryNotFound", mrerr.ErrorKindUser, "category with ID={{ .id }} not found")
)
