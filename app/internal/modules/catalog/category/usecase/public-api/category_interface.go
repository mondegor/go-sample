package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/category/entity/public-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryUseCase interface {
		GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32, languageID uint16) (entity.Category, error)
	}

	CategoryStorage interface {
		NewFetchParams(params entity.CategoryParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Category, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Category, error)
	}
)
