package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	TrademarkAvailabilityName = "Catalog.API.TrademarkAvailability" // TrademarkAvailabilityName - название API
)

type (
	// TrademarkAvailability - comment interface.
	TrademarkAvailability interface {
		// CheckingAvailability - error:
		//    - ErrTrademarkRequired
		//	  - ErrTrademarkNotAvailable
		//	  - ErrTrademarkNotFound
		//	  - Failed
		CheckingAvailability(ctx context.Context, itemID mrtype.KeyInt32) error
	}
)

var (
	// ErrTrademarkRequired - trademark ID is required.
	ErrTrademarkRequired = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogTrademarkRequired", mrerr.ErrorKindUser, "trademark ID is required")

	// ErrTrademarkNotAvailable - trademark with ID is not available.
	ErrTrademarkNotAvailable = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogTrademarkNotAvailable", mrerr.ErrorKindUser, "trademark with ID={{ .id }} is not available")

	// ErrTrademarkNotFound - trademark with ID not found.
	ErrTrademarkNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogTrademarkNotFound", mrerr.ErrorKindUser, "trademark with ID={{ .id }} not found")
)

// TrademarkErrors - comment func.
func TrademarkErrors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrTrademarkRequired,
		ErrTrademarkNotAvailable,
		ErrTrademarkNotFound,
	}
}
