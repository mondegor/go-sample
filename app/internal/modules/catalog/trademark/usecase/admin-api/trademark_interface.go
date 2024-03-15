package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/trademark/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	TrademarkUseCase interface {
		GetList(ctx context.Context, params entity.TrademarkParams) ([]entity.Trademark, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Trademark, error)
		Create(ctx context.Context, item entity.Trademark) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.Trademark) error
		ChangeStatus(ctx context.Context, item entity.Trademark) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	TrademarkStorage interface {
		NewFetchParams(params entity.TrademarkParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Trademark, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Trademark, error)
		FetchStatus(ctx context.Context, row entity.Trademark) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, rowID mrtype.KeyInt32) error
		Insert(ctx context.Context, row entity.Trademark) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.Trademark) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Trademark) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
