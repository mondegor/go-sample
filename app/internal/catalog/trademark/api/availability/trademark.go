package availability

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// TrademarkStorage - comment interface.
	TrademarkStorage interface {
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
	}
)
