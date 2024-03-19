package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/category/entity/public-api"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	CategoryUseCase interface {
		GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error)
		GetItem(ctx context.Context, itemID uuid.UUID, languageID uint16) (entity.Category, error)
	}

	CategoryStorage interface {
		NewSelectParams(params entity.CategoryParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Category, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.Category, error)
	}
)
