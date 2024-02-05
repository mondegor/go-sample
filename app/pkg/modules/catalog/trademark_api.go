package catalog

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	TrademarkAPI interface {
		// CheckingAvailability - error: FactoryErrTrademarkNotFound or Failed
		CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error
	}
)

var (
	FactoryErrTrademarkNotFound = mrerr.NewFactory(
		"errCatalogTrademarkNotFound", mrerr.ErrorKindUser, "trademark with ID={{ .id }} not found")
)
