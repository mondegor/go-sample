package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/category/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Category struct {
		storage       CategoryStorage
		usecaseHelper *mrcore.UsecaseHelper
		imgBaseURL    mrlib.BuilderPath
	}
)

func NewCategory(
	storage CategoryStorage,
	usecaseHelper *mrcore.UsecaseHelper,
	imgBaseURL mrlib.BuilderPath,
) *Category {
	return &Category{
		storage:       storage,
		usecaseHelper: usecaseHelper,
		imgBaseURL:    imgBaseURL,
	}
}

func (uc *Category) GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error) {
	fetchParams := uc.storage.NewFetchParams(params)
	total, err := uc.storage.FetchTotal(ctx, fetchParams.Where)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	if total < 1 {
		return []entity.Category{}, 0, nil
	}

	items, err := uc.storage.Fetch(ctx, fetchParams)

	if err != nil {
		return nil, 0, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	for i := range items {
		items[i].ImageURL = uc.imgBaseURL.FullPath(items[i].ImageURL)
	}

	return items, total, nil
}

func (uc *Category) GetItem(ctx context.Context, itemID mrtype.KeyInt32, languageID uint16) (entity.Category, error) {
	if itemID < 1 {
		return entity.Category{}, mrcore.FactoryErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)

	if err != nil {
		return entity.Category{}, uc.usecaseHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, itemID)
	}

	item.ImageURL = uc.imgBaseURL.FullPath(item.ImageURL)

	return item, nil
}
