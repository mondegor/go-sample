package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Category struct {
		storage       CategoryStorage
		serviceHelper *mrtool.ServiceHelper
	}
)

func NewCategory(
	storage CategoryStorage,
	serviceHelper *mrtool.ServiceHelper,
) *Category {
	return &Category{
		storage:       storage,
		serviceHelper: serviceHelper,
	}
}

func (uc *Category) GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	if total < 1 {
		return []entity.Category{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	return items, total, nil
}

func (uc *Category) GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Category, error) {
	if id < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	item := &entity.Category{ID: id}

	if err := uc.storage.LoadOne(ctx, item); err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, id)
	}

	return item, nil
}
