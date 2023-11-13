package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	TrademarkService interface {
		GetList(ctx context.Context, params entity.TrademarkParams) ([]entity.Trademark, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Trademark, error)
		Create(ctx context.Context, item *entity.Trademark) error
		Store(ctx context.Context, item *entity.Trademark) error
		ChangeStatus(ctx context.Context, item *entity.Trademark) error
		Remove(ctx context.Context, id mrtype.KeyInt32) error
	}

	TrademarkStorage interface {
		NewFetchParams(params entity.TrademarkParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Trademark, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.Trademark) error
		FetchStatus(ctx context.Context, row *entity.Trademark) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
		Insert(ctx context.Context, row *entity.Trademark) error
		Update(ctx context.Context, row *entity.Trademark) error
		UpdateStatus(ctx context.Context, row *entity.Trademark) error
		Delete(ctx context.Context, id mrtype.KeyInt32) error
	}
)
