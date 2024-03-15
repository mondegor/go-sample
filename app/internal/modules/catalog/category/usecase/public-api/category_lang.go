package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/category/entity/public-api"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryLangDecorator struct {
		useCase CategoryUseCase
		dict    *mrlang.MultiLangDictionary
	}
)

func NewCategoryLangDecorator(
	useCase CategoryUseCase,
	dict *mrlang.MultiLangDictionary,
) *CategoryLangDecorator {
	return &CategoryLangDecorator{
		useCase: useCase,
		dict:    dict,
	}
}

func (uc *CategoryLangDecorator) GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error) {
	items, total, err := uc.useCase.GetList(ctx, params)

	if err != nil {
		return nil, 0, err
	}

	if dict := uc.getDict(ctx, params.LanguageID); dict != nil {
		for i := range items {
			items[i].Caption = dict.ItemByID(int(items[i].ID)).Attr("caption", items[i].Caption)
		}
	}

	return items, total, nil
}

func (uc *CategoryLangDecorator) GetItem(ctx context.Context, itemID mrtype.KeyInt32, languageID uint16) (entity.Category, error) {
	item, err := uc.useCase.GetItem(ctx, itemID, languageID)

	if err != nil {
		return entity.Category{}, err
	}

	item.Caption = uc.getDict(ctx, languageID).ItemByID(int(item.ID)).Attr("caption", item.Caption)

	return item, nil
}

func (uc *CategoryLangDecorator) getDict(ctx context.Context, languageID uint16) *mrlang.Dictionary {
	dict, err := uc.dict.ByLangID(languageID)

	if err != nil {
		mrlog.Ctx(ctx).Warn().Err(err)
		return nil
	}

	return dict
}
