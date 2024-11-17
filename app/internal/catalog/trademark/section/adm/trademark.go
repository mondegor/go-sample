package adm

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/entity"
)

type (
	// TrademarkUseCase - comment interface.
	TrademarkUseCase interface {
		GetList(ctx context.Context, params entity.TrademarkParams) (items []entity.Trademark, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.Trademark, error)
		Create(ctx context.Context, item entity.Trademark) (itemID uint64, err error)
		Store(ctx context.Context, item entity.Trademark) error
		ChangeStatus(ctx context.Context, item entity.Trademark) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// TrademarkStorage - comment interface.
	TrademarkStorage interface {
		FetchWithTotal(ctx context.Context, params entity.TrademarkParams) (rows []entity.Trademark, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.Trademark, error)
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Trademark) (rowID uint64, err error)
		Update(ctx context.Context, row entity.Trademark) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.Trademark) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
