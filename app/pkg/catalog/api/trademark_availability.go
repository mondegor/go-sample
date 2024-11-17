package api

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
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
		CheckingAvailability(ctx context.Context, itemID uint64) error
	}
)

var (
	// ErrTrademarkRequired - trademark ID is required.
	ErrTrademarkRequired = mrerr.NewProto(
		"catalog.errTrademarkRequired", mrerr.ErrorKindUser, "trademark ID is required")

	// ErrTrademarkNotAvailable - trademark with ID is not available.
	ErrTrademarkNotAvailable = mrerr.NewProto(
		"catalog.errTrademarkNotAvailable", mrerr.ErrorKindUser, "trademark with ID={{ .id }} is not available")

	// ErrTrademarkNotFound - trademark with ID not found.
	ErrTrademarkNotFound = mrerr.NewProto(
		"catalog.errTrademarkNotFound", mrerr.ErrorKindUser, "trademark with ID={{ .id }} not found")
)
