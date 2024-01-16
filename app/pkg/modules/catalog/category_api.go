package catalog

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryAPI interface {
		// CheckingAvailability - error: FactoryErrCategoryNotFound or Failed
		CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error
	}
)
