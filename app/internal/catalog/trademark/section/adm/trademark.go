package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/entity"
)

type (
	// TrademarkUseCase - comment interface.
	TrademarkUseCase interface {
		GetList(ctx context.Context, params entity.TrademarkParams) ([]entity.Trademark, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Trademark, error)
		Create(ctx context.Context, item entity.Trademark) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.Trademark) error
		ChangeStatus(ctx context.Context, item entity.Trademark) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	// TrademarkStorage - comment interface.
	TrademarkStorage interface {
		NewSelectParams(params entity.TrademarkParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Trademark, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Trademark, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Trademark) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.Trademark) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Trademark) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
