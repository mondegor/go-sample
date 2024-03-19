package catalog

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	TrademarkAPI interface {
		// CheckingAvailability - error: FactoryErrTrademarkRequired | FactoryErrTrademarkNotFound |
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	FactoryErrTrademarkRequired = mrerr.NewFactory(
		"errCatalogTrademarkRequired", mrerr.ErrorKindUser, "trademark ID is required")

	FactoryErrTrademarkNotFound = mrerr.NewFactory(
		"errCatalogTrademarkNotFound", mrerr.ErrorKindUser, "trademark with ID={{ .id }} not found")
)
