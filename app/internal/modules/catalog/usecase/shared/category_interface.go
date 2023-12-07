package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryServiceAPI interface {
		CheckingAvailability(ctx context.Context, id mrtype.KeyInt32) error
	}

	CategoryStorage interface {
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
	}
)
