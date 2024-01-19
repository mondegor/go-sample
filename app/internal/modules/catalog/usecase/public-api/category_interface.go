package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/public-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryService interface {
		GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32, languageID uint16) (*entity.Category, error)
	}

	CategoryStorage interface {
		NewFetchParams(params entity.CategoryParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Category, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.Category) error
	}
)
